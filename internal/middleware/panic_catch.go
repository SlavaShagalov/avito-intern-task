package middleware

import (
	pErrors "github.com/SlavaShagalov/avito-intern-task/internal/pkg/errors"
	pHTTP "github.com/SlavaShagalov/avito-intern-task/internal/pkg/http"
	"go.uber.org/zap"
	"net/http"
)

func NewPanicCatch(log *zap.Logger) func(handler http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					log.Error("Panic occurred", zap.Any("error", err))
					pHTTP.HandleError(w, r, pErrors.ErrUnknown)
				}
			}()
			handler.ServeHTTP(w, r)
		})
	}
}
