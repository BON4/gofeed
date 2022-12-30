package domain

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/multierr"
)

type sessionBody struct {
	Id           string             `json:"id"`
	Refreshtoken string             `json:"refreshtoken"`
	UserAgent    string             `json:"user_agent"`
	ClientIp     string             `json:"client_ip"`
	Blocked      bool               `json:"blocked"`
	IssuedAt     time.Time          `json:"issued_at"`
	ExpiresAt    time.Time          `json:"expires_at"`
	Instance     AccountCredentials `json:"instance"`
}

type Session struct {
	sessionBody
}

func (s *Session) IsBlocked() bool {
	return s.Blocked
}

func (s *Session) GetTTL() time.Duration {
	return s.ExpiresAt.Sub(time.Now())
}

func (s *Session) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.sessionBody)
}

const YEAR = time.Second * 31536000

type SessionFactoryConfig struct {
	SessionMinTTL time.Duration
	SessionMaxTTL time.Duration
}

func (f *SessionFactoryConfig) Validate() error {
	var err error

	if f.SessionMaxTTL < time.Minute && f.SessionMaxTTL > YEAR {
		err = multierr.Append(
			err,
			errors.Errorf(
				"SessionMaxTTL should be greater than 60s and less then one year (in seconds), but is %d",
				f.SessionMaxTTL,
			),
		)
	}

	if f.SessionMinTTL < time.Minute && f.SessionMinTTL > YEAR {
		err = multierr.Append(
			err,
			errors.Errorf(
				"SessionMaxTTL should be greater than 60s and less then one year (in seconds), but is %d",
				f.SessionMinTTL,
			),
		)
	}

	if f.SessionMinTTL > f.SessionMaxTTL {
		err = multierr.Append(
			err,
			errors.Errorf(
				"SessionMaxTTL: %d should be greater than SessionMinTTL: %d",
				f.SessionMaxTTL, f.SessionMinTTL,
			),
		)
	}

	return err
}

type SessionFactory struct {
	fc SessionFactoryConfig
}

func (f *SessionFactory) UnmarshalSessionJSON(b []byte) (*Session, error) {
	ss := sessionBody{}
	if err := json.Unmarshal(b, &ss); err != nil {
		return nil, err
	}

	if err := f.validateSession(&ss); err != nil {
		return nil, err
	}

	return &Session{ss}, nil
}

func NewSessionfactory(fc SessionFactoryConfig) (SessionFactory, error) {
	if err := fc.Validate(); err != nil {
		return SessionFactory{}, errors.Wrap(err, "invalid config passed to factory")
	}

	return SessionFactory{fc: fc}, nil
}

func (f *SessionFactory) NewSession(
	id,
	refreshtoken,
	userAgent,
	clientIp string,
	issuedAt time.Time,
	expiresAt time.Time,
	instance AccountCredentials) (*Session, error) {

	ss := sessionBody{
		Id:           id,
		Refreshtoken: refreshtoken,
		UserAgent:    userAgent,
		ClientIp:     clientIp,
		IssuedAt:     issuedAt,
		ExpiresAt:    expiresAt,
		Blocked:      false,
		Instance:     instance,
	}

	if err := f.validateSession(&ss); err != nil {
		return nil, err
	}

	return &Session{ss}, nil
}

type TooLongTTLError struct {
	MaxSessionTTL time.Duration
	ProvidedTime  time.Duration
}

func (e TooLongTTLError) Error() string {
	return fmt.Sprintf(
		"too long ttl, max ttl: %s, provided time: %s",
		e.MaxSessionTTL,
		e.ProvidedTime,
	)
}

type TooSmallTTLError struct {
	MinSessionTTL time.Duration
	ProvidedTime  time.Duration
}

func (e TooSmallTTLError) Error() string {
	return fmt.Sprintf(
		"too small ttl, min ttl: %s, provided time: %s",
		e.MinSessionTTL,
		e.ProvidedTime,
	)
}

func (f *SessionFactory) validateSession(ss *sessionBody) error {
	sessDur := ss.ExpiresAt.Sub(ss.IssuedAt)
	if !ss.Blocked {
		if sessDur > f.fc.SessionMaxTTL {
			return TooLongTTLError{
				MaxSessionTTL: f.fc.SessionMaxTTL,
				ProvidedTime:  ss.ExpiresAt.Sub(ss.IssuedAt),
			}
		}

		if sessDur < f.fc.SessionMinTTL {
			return TooSmallTTLError{
				MinSessionTTL: f.fc.SessionMinTTL,
				ProvidedTime:  ss.ExpiresAt.Sub(ss.IssuedAt),
			}
		}
	}

	if len(ss.Id) == 0 {
		return errors.New("Id not specified")
	}

	if len(ss.Refreshtoken) == 0 {
		return errors.New("refresh token not specified")
	}

	return nil
}
