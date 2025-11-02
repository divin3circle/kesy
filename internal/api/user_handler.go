package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	store "github.com/nhx-finance/kesy/internal/stores"
	"github.com/nhx-finance/kesy/internal/utils"
)


type UserHandler struct {
	UserStore store.PostgresUserStore
	Logger *log.Logger
}

func NewUserHandler(userStore store.PostgresUserStore, logger *log.Logger) *UserHandler {
	return &UserHandler{
		UserStore: userStore,
		Logger: logger,
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

func (uh *UserHandler) HandleCreateUser(w http.ResponseWriter, r *http.Request){
	var userReq CreateUserRequest

	err := json.NewDecoder(r.Body).Decode(&userReq)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		uh.Logger.Printf("failed to decode user request: %v", err)
		return
	}

	err = validateCreateUserRequest(userReq)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		uh.Logger.Printf("failed to validate user request: %v", err)
		return
	}

	userToCreate := &store.User{
		Email: userReq.Email,
		FirstName: userReq.FirstName,
		LastName: userReq.LastName,
		DOB: userReq.DOB,
		ResidenceCountry: userReq.ResidenceCountry,
		Province: userReq.Province,
		Timezone: userReq.Timezone,
		AcceptedTerms: userReq.AcceptedTerms,
		KYCStatus: "pending",
	}

	err = userToCreate.PasswordHash.Set(userReq.Password)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": err.Error()})
		uh.Logger.Printf("failed to set password: %v", err)
		return
	}

	user, err := uh.UserStore.Create(userToCreate)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": err.Error()})
		uh.Logger.Printf("failed to create user: %v", err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"user": user})
	uh.Logger.Printf("user created successfully: %v", user)
}

func validateCreateUserRequest(userReq CreateUserRequest) error {
	if userReq.Email == "" {
		return errors.New("email is required")
	}
	if userReq.Password == "" {
		return errors.New("password is required")
	}
	return nil
}



