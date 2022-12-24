package domain

import (
	"context"
	"fmt"
	"sync"

	"github.com/pkg/errors"
	"go.uber.org/multierr"
)

type getUsers func(ctx context.Context) []*User

type AccountRole string

// Enum values for AccountRole
const (
	AccountRoleAdmin AccountRole = "admin"
	AccountRoleBasic AccountRole = "basic"
)

var account_roles = [2]AccountRole{AccountRoleAdmin, AccountRoleBasic}

func (e AccountRole) IsValid() error {
	switch e {
	case AccountRoleAdmin, AccountRoleBasic:
		return nil
	default:
		return errors.New("enum is not valid")
	}
}

type Account struct {
	username string
	password string
	role     AccountRole

	users       []*User
	once        *sync.Once
	usersLoader getUsers
}

func (a *Account) GetUsers(ctx context.Context) []*User {
	a.once.Do(func() {
		a.users = a.usersLoader(ctx)
	})
	return a.users
}

func (a *Account) setUsers(f getUsers) {
	a.usersLoader = f
	a.once = &sync.Once{}
}

type FactoryConfig struct {
	MinUsernameLen int
	MinPasswordLen int
	DefaultRole    AccountRole
}

func (f FactoryConfig) Validate() error {
	var err error

	if f.MinUsernameLen < 1 {
		err = multierr.Append(
			err,
			errors.Errorf(
				"MinUsernameLen should be greater than 1, but is %d",
				f.MinUsernameLen,
			),
		)
	}

	if f.MinPasswordLen < 1 {
		err = multierr.Append(
			err,
			errors.Errorf(
				"MinPasswordLen should be greater than 1, but is %dn",
				f.MinPasswordLen,
			),
		)
	}

	if enumErr := f.DefaultRole.IsValid(); enumErr != nil {
		err = multierr.Append(
			err,
			enumErr,
		)
	}

	return err
}

type Factory struct {
	fc FactoryConfig
}

func NewFactory(fc FactoryConfig) (*Factory, error) {
	if err := fc.Validate(); err != nil {
		return nil, errors.Wrap(err, "invalid config passed to factory")
	}

	return &Factory{fc: fc}, nil
}

func (f *Factory) NewAccount(username, email, password string, role AccountRole) (*Account, error) {
	a := &Account{}
	if err := f.validateAccount(a); err != nil {
		return nil, err
	}

	// TODO: hash password before returning account

	return a, nil
}

// UnmarshalAccountFromDatabase - unmarshals account from the database.
//
// It should be used only for unmarshalling from the database!
func (f *Factory) UnmarshalAccountFromDatabase(username, email, password string, role AccountRole, lazyGetter getUsers) (*Account, error) {
	a := &Account{}
	if lazyGetter == nil {
		a.setUsers(func(ctx context.Context) []*User {
			return []*User{}
		})
	} else {
		a.setUsers(lazyGetter)
	}

	if err := f.validateAccount(a); err != nil {
		// Ignore password error, beacouse we unmarshaling from DB
		if _, ok := err.(TooShortPassword); !ok {
			return nil, err
		}
	}

	return a, nil
}

type TooShortUsername struct {
	MinUsernameLen   int
	ProvidedUsername string
}

func (e TooShortUsername) Error() string {
	return fmt.Sprintf(
		"Too short account username, min username length: %d, provided username: %s",
		e.MinUsernameLen,
		e.ProvidedUsername,
	)
}

type TooShortPassword struct {
	MinPasswordLen      int
	ProvidedPasswordLen int
}

func (e TooShortPassword) Error() string {
	return fmt.Sprintf(
		"Too short account password, min password length: %d, provided password len: %d",
		e.MinPasswordLen,
		e.ProvidedPasswordLen,
	)
}

type InvalidAccountRole struct {
	Roles               []AccountRole
	ProvidedAccountRole string
}

func (e InvalidAccountRole) Error() string {
	return fmt.Sprintf(
		"Invalid account role, valid roles: %v, provided role: %s",
		e.Roles,
		e.ProvidedAccountRole,
	)
}

func (f *Factory) validateAccount(acc *Account) error {
	if len(acc.username) < f.fc.MinUsernameLen {
		return TooShortUsername{
			MinUsernameLen:   f.fc.MinUsernameLen,
			ProvidedUsername: acc.username,
		}
	}

	if len(acc.password) < f.fc.MinPasswordLen {
		return TooShortPassword{
			MinPasswordLen:      f.fc.MinPasswordLen,
			ProvidedPasswordLen: len(acc.password),
		}
	}

	if err := acc.role.IsValid(); err != nil {
		return InvalidAccountRole{
			Roles:               account_roles[:],
			ProvidedAccountRole: string(acc.role),
		}
	}

	return nil
}
