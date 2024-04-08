package middleware

import (
	pErrors "github.com/SlavaShagalov/avito-intern-task/internal/pkg/errors"
	pHTTP "github.com/SlavaShagalov/avito-intern-task/internal/pkg/http"
	"go.uber.org/zap"
	"net/http"
)

func NewCheckAdminAccess(log *zap.Logger) Middleware {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			isAdmin, ok := r.Context().Value(ContextIsAdmin).(bool)
			if !ok {
				log.Error("Check admin access: is_admin field not found")
				pHTTP.HandleError(w, r, pErrors.ErrReadBody)
				return
			}
			if !isAdmin {
				pHTTP.HandleError(w, r, pErrors.ErrAdminRequired)
				return
			}
			h(w, r)
		}
	}
}
