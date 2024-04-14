package middleware

import (
	"context"
	"github.com/SlavaShagalov/avito-intern-task/internal/pkg/config"
	pErrors "github.com/SlavaShagalov/avito-intern-task/internal/pkg/errors"
	pHTTP "github.com/SlavaShagalov/avito-intern-task/internal/pkg/http"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"net/http"
)

func NewCheckAuth(log *zap.Logger) Middleware {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			tokenString := r.Header.Get("token")

			parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					log.Warn("Unexpected signing method", zap.Any("alg", token.Header["alg"]))
					return nil, pErrors.ErrInvalidAuthToken
				}

				return []byte(viper.GetString(config.AuthKey)), nil
			})
			if err != nil {
				pHTTP.HandleError(w, r, pErrors.ErrInvalidAuthToken)
				return
			}

			claims, ok := parsedToken.Claims.(jwt.MapClaims)
			if !ok {
				pHTTP.HandleError(w, r, pErrors.ErrInvalidAuthToken)
				return
			}

			ctx := context.WithValue(r.Context(), ContextUserID, claims["user_id"])
			ctx = context.WithValue(ctx, ContextIsAdmin, claims["is_admin"])

			h(w, r.WithContext(ctx))
		}
	}
}
