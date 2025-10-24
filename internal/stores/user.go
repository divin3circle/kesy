package store

import (
	"errors"

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
	PasswordHash			string `json:"-"`
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
