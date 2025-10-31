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

type CreateUserRequest struct {}

