// Package auth provides middlewares and functions for dealing with authentication.
package auth

import (
	"context"
	"firebase.google.com/go/auth"
	"github.com/steven-eckhoff/gopher-cache-open/internal/common/errors"
	"github.com/steven-eckhoff/gopher-cache-open/internal/common/server/httperr"
	"net/http"
	"strings"
)

// FirebaseHttpMiddleware performs authentications using Firebase.
type FirebaseHttpMiddleware struct {
	AuthClient *auth.Client
}

// Middleware will place a User type into the context after successfully authenticating
// with Firebase.
func (a FirebaseHttpMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		bearerToken := a.tokenFromHeader(r)
		if bearerToken == "" {
			httperr.Unauthorised("empty-bearer-token", nil, w, r)
			return
		}

		token, err := a.AuthClient.VerifyIDToken(ctx, bearerToken)
		if err != nil {
			httperr.Unauthorised("unable-to-verify-jwt", err, w, r)
			return
		}

		name, ok := token.Claims["name"].(string)
		if !ok {
			httperr.InternalError("invalid-token", nil, w, r)
			return
		}

		email, ok := token.Claims["email"].(string)
		if !ok {
			httperr.InternalError("invalid-token", nil, w, r)
			return
		}

		number, ok := token.Claims["number"].(string)
		if !ok {
			httperr.InternalError("invalid-token", nil, w, r)
			return
		}

		ctx = context.WithValue(ctx, userContextKey, User{
			UUID:        token.UID,
			DisplayName: name,
			Email:       email,
			Number:      number,
		})
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func (a FirebaseHttpMiddleware) tokenFromHeader(r *http.Request) string {
	headerValue := r.Header.Get("Authorization")

	if len(headerValue) > 7 && strings.ToLower(headerValue[0:6]) == "bearer" {
		return headerValue[7:]
	}

	return ""
}

// User holds the authenticated user's info
type User struct {
	UUID        string
	DisplayName string
	Email       string
	Number      string
}

type ctxKey int

const (
	userContextKey ctxKey = iota
)

var (
	ErrorNoUserInContext = errors.NewAuthorizationError("context has no user", "missing-context-user")
)

// UserFromContext will check the context for an authenticated user.
// If none is found then ErrorNoUserInContext will be returned.
func UserFromContext(ctx context.Context) (User, error) {
	u, ok := ctx.Value(userContextKey).(User)
	if ok {
		return u, nil
	}

	return User{}, ErrorNoUserInContext
}
