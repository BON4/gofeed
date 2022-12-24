package adapters

type PostgresAccountsRepository struct{}

func NewPostgresAccountsRepository() *PostgresAccountsRepository {
	return &PostgresAccountsRepository{}
}
