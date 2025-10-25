package middleware

import (
	"context"
	"net/http"
	"strings"

	store "github.com/nhx-finance/kesy/internal/stores"
	"github.com/nhx-finance/kesy/internal/tokens"
	"github.com/nhx-finance/kesy/internal/utils"
)

type UserAuthMiddleware struct {
	UserStore store.UserStore
}

func NewUserAuthMiddleware(userStore store.UserStore) *UserAuthMiddleware {
	return &UserAuthMiddleware{
		UserStore: userStore,
	}
}

type contextKey string

const (
	UserContextKey contextKey = "user"
)

func SetUser(r *http.Request, user *store.User) *http.Request{
	ctx := context.WithValue(r.Context(), UserContextKey,user)
	return r.WithContext(ctx)
}

func GetUser(r *http.Request) *store.User {
	user, ok := r.Context().Value(UserContextKey).(*store.User)
	if !ok {
		panic("user not found in context")
	}
	return user
}

func (uam *UserAuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Vary", "Authorization")

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			r = SetUser(r, store.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "invalid authorization header"})
			return 
		}

		tokenPlainText := headerParts[1]
		user, err := uam.UserStore.GetToken(tokens.ScopeAuthentication, tokenPlainText)
		if err != nil {
			utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "invalid token"})
			return
		}
		if user == nil {
			utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "token expired or not found"})
			return
		}

		r = SetUser(r, user)
		next.ServeHTTP(w,r)
	})
}

func (uam *UserAuthMiddleware) RequireAuthenticatedUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := GetUser(r)
		if user.IsAnonymous() {
			utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "authentication required"})
			return
		}
		next.ServeHTTP(w, r)
	})
}