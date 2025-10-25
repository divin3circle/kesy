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

func (pts *PostgresTokenStore) Insert(token *tokens.Token) error {
	query := `
	INSERT INTO tokens (hash, user_id, expiry, scope)
	VALUES ($1, $2, $3, $4)
	`
	_, err := pts.db.Exec(query, string(token.Hash), token.UserID, token.Expiry, token.Scope)
	return err
}

func (pts *PostgresTokenStore) Create(userID string, ttl time.Duration, scope string) (*tokens.Token, error) {
	token, err := tokens.GenerateToken(userID, ttl, scope)
	if err != nil {
		return nil, err
	}
	err = pts.Insert(token)

	return token, err
}

func (pts *PostgresTokenStore) DeleteForUser(scope string, userID string) error {
	query := `
	DELETE FROM tokens
	WHERE user_id = $1 AND scope = $2
	`
	_, err := pts.db.Exec(query, userID, scope)
	return err
}