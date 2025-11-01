package api

import store "github.com/nhx-finance/kesy/internal/stores"


type UserHandler struct {
	UserStore store.PostgresUserStore
}

func NewUserHandler(userStore store.PostgresUserStore) *UserHandler {
	return &UserHandler{
		UserStore: userStore,
	}
}

type CreateUserRequest struct {
	Email string `json:"email"`
	Password string `json:"password"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	DOB string `json:"dob"`
	ResidenceCountry string `json:"residence_country"`
	Province string `json:"province"`
	Timezone string `json:"timezone"`
	AcceptedTerms bool `json:"accepted_terms"`
}

