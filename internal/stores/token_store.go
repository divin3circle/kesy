package store

import (
	"database/sql"
	"time"

	"github.com/nhx-finance/kesy/internal/tokens"
)

type PostgresTokenStore struct {
	db *sql.DB
}

func NewPostgresTokenStore(db *sql.DB) *PostgresTokenStore {
	return &PostgresTokenStore{
		db: db,
	}
}

type TokenStore interface {
	Insert(token *tokens.Token) error
	Create(userID string, ttl time.Duration, scope string) (*tokens.Token, error)
	DeleteForUser(scope string, userID string) error
}

func (pts *PostgresTokenStore) Insert(token *tokens.Token) error {}

func (pts *PostgresTokenStore) Create(userID string, ttl time.Duration, scope string) (*tokens.Token, error) {}

func (pts *PostgresTokenStore) DeleteForUser(scope string, userID string) error {}