package domain

import (
	"time"

	"github.com/BON4/gofeed/internal/common/tokens"
	"github.com/pkg/errors"
	"go.uber.org/multierr"
)

type Token struct {
	AccessToken    string
	AccessPayload  *tokens.Payload[tokens.InstanceCredentials]
	RefreshToken   string
	RefreshPayload *tokens.Payload[tokens.InstanceCredentials]
}

type TokenFactoryConfig struct {
	RefreshTokenDuration time.Duration
	AccessTokenDuration  time.Duration
	TokenSecret          string
}

func (f TokenFactoryConfig) Validate() error {
	var err error

	if len(f.TokenSecret) < 32 {
		err = multierr.Append(
			err,
			errors.Errorf("TokenSecret length should be greater then 32, but is: %d", len(f.TokenSecret)),
		)
	}

	if f.AccessTokenDuration < time.Minute {
		err = multierr.Append(
			err,
			errors.Errorf("AccessTokenDuration should be greater then %s, but is: %s", time.Minute, f.AccessTokenDuration),
		)
	}

	if f.RefreshTokenDuration < time.Minute*10 {
		err = multierr.Append(
			err,
			errors.Errorf("RefreshTokenDuration should be greater then %s, but is: %s", time.Hour, f.RefreshTokenDuration),
		)
	}

	return err
}

type TokenFactory struct {
	tokenGen tokens.Generator[tokens.InstanceCredentials]
	fc       TokenFactoryConfig
}

func NewAuthTokenFactory(fc TokenFactoryConfig) (*TokenFactory, error) {
	if err := fc.Validate(); err != nil {
		return nil, errors.Wrap(err, "invalid config passed to factory")
	}

	gen, err := tokens.NewJWTGenerator[tokens.InstanceCredentials](fc.TokenSecret)
	if err != nil {
		return nil, err
	}

	return &TokenFactory{fc: fc, tokenGen: gen}, nil
}

func (f *TokenFactory) NewTokenPair(instance tokens.InstanceCredentials) (*Token, error) {
	acessToken, acessPayload, err := f.tokenGen.CreateToken(instance, f.fc.AccessTokenDuration)
	if err != nil {
		return nil, err
	}

	refreshToken, refreshPayload, err := f.tokenGen.CreateToken(instance, f.fc.RefreshTokenDuration)
	if err != nil {
		return nil, err
	}

	return &Token{
		AccessToken:    acessToken,
		AccessPayload:  acessPayload,
		RefreshToken:   refreshToken,
		RefreshPayload: refreshPayload,
	}, nil
}

func (f *TokenFactory) NewAccesToken(refreshToken string) (*Token, error) {
	refreshPayload, err := f.tokenGen.VerifyToken(refreshToken)
	if err != nil {
		return nil, err
	}

	acessToken, acessPayload, err := f.tokenGen.CreateToken(refreshPayload.Instance, f.fc.AccessTokenDuration)
	if err != nil {
		return nil, err
	}

	return &Token{
		AccessToken:    acessToken,
		AccessPayload:  acessPayload,
		RefreshToken:   refreshToken,
		RefreshPayload: refreshPayload,
	}, nil
}

func (f *TokenFactory) VerifyToken(token string) (*tokens.Payload[tokens.InstanceCredentials], error) {
	return f.tokenGen.VerifyToken(token)
}

type TokenVerifier struct {
	tokenGen tokens.Generator[tokens.InstanceCredentials]
}

func NewTokenVerifier(secretToken string) (*TokenVerifier, error) {
	gen, err := tokens.NewJWTGenerator[tokens.InstanceCredentials](secretToken)
	if err != nil {
		return nil, err
	}

	return &TokenVerifier{
		tokenGen: gen,
	}, nil
}

func (f *TokenVerifier) VerifyToken(token string) (*tokens.Payload[tokens.InstanceCredentials], error) {
	return f.tokenGen.VerifyToken(token)
}
