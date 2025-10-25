package store

import (
	"crypto/sha256"
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type password struct {
	plainText string
	hash      string
}

func (p *password) Set(plainText string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainText), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	p.plainText = plainText
	p.hash = string(hash)
	return nil
}

func (p *password) Matches(plainText string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(p.hash), []byte(plainText))
	if err != nil {
		switch {
			case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
				return false, nil
			default:
				return false, err
		}
	}
	return true, nil
}

func (p *password) String() string {
	return p.plainText
}

type User struct {
	ID	   					string `json:"id"`
	Email  					string `json:"email"`
	PasswordHash			password `json:"-"`
	FirstName   			string `json:"first_name"`
	LastName    			string `json:"last_name"`
	DOB    					string `json:"dob"`
	ResidenceCountry		string `json:"residence_country"`
	Province 				string `json:"province"`
	Timezone				string `json:"timezone"`
	AcceptedTerms			bool   `json:"accepted_terms"`
	KYCStatus				string `json:"kyc_status"`
	CreatedAt 				string `json:"created_at"`
	UpdatedAt 				string `json:"updated_at"`
	DeletedAt 				*string `json:"deleted_at,omitempty"`
}

var AnonymousUser = &User{}

func (u *User) IsAnonymous() bool {
	return u.ID == AnonymousUser.ID
}

type PostgresUserStore struct {
	db *sql.DB
}

func NewPostgresUserStore(db *sql.DB) *PostgresUserStore {
	return &PostgresUserStore{
		db: db,
	}
}

type UserStore interface {
	Create(user *User) (*User, error)
	GetToken(scope string, tokenPlainText string) (*User, error)
}

func (pus *PostgresUserStore) Create(user *User) (*User, error) {
	query := `
	INSERT INTO users (email, password_hash, first_name, last_name, dob, residence_country, province, timezone, accepted_terms)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	RETURNING id, created_at, updated_at
	`
	err := pus.db.QueryRow(query, user.Email, user.PasswordHash.hash, user.FirstName, user.LastName, user.DOB, user.ResidenceCountry, user.Province, user.Timezone, user.AcceptedTerms).
	Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (pus *PostgresUserStore) GetToken(scope string, tokenPlainText string) (*User, error) {
	tokenHash := sha256.Sum256([]byte(tokenPlainText))

	query := `
	SELECT users.id, users.email, users.password_hash, users.first_name, users.last_name, users.dob, users.residence_country, users.province, users.timezone, users.accepted_terms, users.kyc_status, users.created_at, users.updated_at
	FROM users
	INNER JOIN tokens
	ON users.id = tokens.user_id
	WHERE tokens.hash = $1
	AND tokens.scope = $2
	AND tokens.expiry > $3
	`

	user := &User{
		PasswordHash: password{},
	}

	err := pus.db.QueryRow(query, string(tokenHash[:]), scope,  time.Now()).
	Scan(&user.ID, &user.Email, &user.PasswordHash.hash, &user.FirstName, &user.LastName, &user.DOB, &user.ResidenceCountry, &user.Province, &user.Timezone, &user.AcceptedTerms, &user.KYCStatus, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}
