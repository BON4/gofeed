package adapters

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"

	"github.com/BON4/gofeed/internal/accounts/adapters/sqlc"
	"github.com/BON4/gofeed/internal/accounts/domain"
)

type PostgresAccountsRepository struct {
	querys sqlc.Store
	fc     *domain.Factory
}

func NewPostgresAccountsRepository(dbconn *sql.DB, accFactory *domain.Factory) *PostgresAccountsRepository {
	return &PostgresAccountsRepository{
		querys: sqlc.NewStore(dbconn),
		fc:     accFactory,
	}
}

func (p *PostgresAccountsRepository) CreateAccount(ctx context.Context, acc *domain.Account) (*domain.Account, error) {
	created, err := p.querys.CreateAccount(ctx, sqlc.CreateAccountParams{
		Username:  acc.GetUsername(),
		Password:  acc.GetPassword(),
		Email:     acc.GetEmail(),
		Role:      sqlc.AccountRole(acc.GetRole()),
		Activated: acc.Activeted(),
	})

	if err != nil {
		return nil, err
	}

	return p.unmarshalAccount(created, nil)
}

func (p *PostgresAccountsRepository) GetAccount(ctx context.Context, username string) (*domain.Account, error) {
	found, err := p.querys.GetAccount(ctx, username)
	if err != nil {
		return nil, err
	}

	return p.unmarshalAccount(found, func(ctx context.Context) ([]*domain.User, error) {
		users, err := p.querys.GetUsersByAccount(ctx, username)
		if err != nil {
			return []*domain.User{}, err
		}

		dUsers := make([]*domain.User, len(users))
		for _, usr := range users {
			dUsers = append(dUsers, p.unmarshalUser(usr))
		}
		return dUsers, nil
	})
}

func (p *PostgresAccountsRepository) unmarshalAccount(acc *sqlc.Account, f domain.GetUsers) (*domain.Account, error) {
	return p.fc.UnmarshalAccountFromDatabase(
		acc.Username,
		acc.Email,
		acc.Password,
		domain.AccountRole(acc.Role),
		acc.Activated,
		acc.PasswordChangedAt,
		f)
}

func (p *PostgresAccountsRepository) unmarshalUser(usr *sqlc.User) *domain.User {
	return &domain.User{
		Uuid:    usr.Uuid,
		Ip:      usr.Ip.String,
		Os:      usr.Os.String,
		Browser: usr.Browser.String,
	}
}

func NewPostgresConnection(connStr string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return db, nil
}
