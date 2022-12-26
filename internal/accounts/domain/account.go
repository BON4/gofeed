package domain

import (
	"context"
	"fmt"
	"sync"

	pswrd "github.com/BON4/gofeed/internal/common/password"
	"github.com/pkg/errors"
	"go.uber.org/multierr"
)

type GetUsers func(ctx context.Context) ([]*User, error)

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
	email    string
	password []byte
	role     AccountRole

	users       []*User
	once        *sync.Once
	usersLoader GetUsers
}

func (a *Account) GetUsername() string {
	return a.username
}

func (a *Account) GetEmail() string {
	return a.email
}

func (a *Account) GetPassword() []byte {
	return a.password
}

func (a *Account) GetRole() AccountRole {
	return a.role
}

func (a *Account) GetUsers(ctx context.Context) ([]*User, error) {
	var err error
	a.once.Do(func() {
		a.users, err = a.usersLoader(ctx)
	})

	if err != nil {
		return a.users, err
	}
	return a.users, nil
}

func (a *Account) setUsers(f GetUsers) {
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
	a := &Account{
		username: username,
		email:    email,
		role:     role,
	}

	err := f.validateAccount(a)
	if err != nil {
		return nil, err
	}

	a.password, err = pswrd.HashPassword(password)
	if err != nil {
		return nil, err
	}

	return a, nil
}

// UnmarshalAccountFromDatabase - unmarshals account from the database.
//
// It should be used only for unmarshalling from the database!
func (f *Factory) UnmarshalAccountFromDatabase(username, email string, password []byte, role AccountRole, lazyGetter getUsers) (*Account, error) {
	a := &Account{
		username: username,
		email:    email,
		password: password,
		role:     role,
	}

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
