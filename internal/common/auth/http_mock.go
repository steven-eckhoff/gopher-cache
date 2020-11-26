package auth

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/steven-eckhoff/gopher-cache-open/internal/common/server/httperr"
	"net/http"
)

// HttpMockMiddleware is used in the local environment which doesn't depend on Firebase.
// It expects a JWT in the request header signed with the word mock_secret. The JWT must
// contain a payload with the following strings:
// - id
// - name
// - email
// - number
func HttpMockMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var claims jwt.MapClaims
		token, err := request.ParseFromRequest(
			r,
			request.AuthorizationHeaderExtractor,
			func(token *jwt.Token) (i interface{}, e error) {
				return []byte("mock_secret"), nil
			},
			request.WithClaims(&claims),
		)
		if err != nil {
			httperr.BadRequest("unable-to-get-jwt", err, w, r)
			return
		}

		if !token.Valid {
			httperr.BadRequest("invalid-jwt", nil, w, r)
			return
		}

		id, ok := claims["id"].(string)
		if !ok {
			httperr.BadRequest("invalid-jwt", nil, w, r)
			return
		}

		name, ok := claims["name"].(string)
		if !ok {
			httperr.BadRequest("invalid-jwt", nil, w, r)
			return
		}

		email, ok := claims["email"].(string)
		if !ok {
			httperr.BadRequest("invalid-jwt", nil, w, r)
			return

		}

		number, ok := claims["number"].(string)
		if !ok {
			httperr.BadRequest("invalid-jwt", nil, w, r)
			return
		}

		ctx := context.WithValue(r.Context(), userContextKey, User{
			UUID:        id,
			DisplayName: name,
			Email:       email,
			Number:      number,
		})
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
