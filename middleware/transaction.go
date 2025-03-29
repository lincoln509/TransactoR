package middleware

import (
	"context"
	"net/http"

	"TransactoR/context"
	"TransactoR/database"
	"TransactoR/logging"

	"gorm.io/gorm"
)

type TransactionHandler struct {
	DB     *gorm.DB
	Logger logging.Logger
}

func NewTransactionHandler(logger logging.Logger) *TransactionHandler {
	return &TransactionHandler{
		DB:     database.GetDB(),
		Logger: logger,
	}
}

func (th *TransactionHandler) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tx := th.DB.Begin()
		th.Logger.Info(r, "Transaction démarrée")

		// Création d'un wrapper pour capturer le status
		rw := &responseWriter{w, http.StatusOK}

		defer func() {
			if rec := recover(); rec != nil {
				tx.Rollback()
				th.Logger.Error(r, "Transaction annulée (panic): %v", rec)
				http.Error(w, "Erreur interne", http.StatusInternalServerError)
			}

			if rw.status >= 400 {
				tx.Rollback()
				th.Logger.Warn(r, "Transaction annulée (status %d)", rw.status)
			} else {
				tx.Commit()
				th.Logger.Info(r, "Transaction validée")
			}
		}()

		ctx := context.WithValue(r.Context(), dbctx.TxKey, tx)
		next.ServeHTTP(rw, r.WithContext(ctx))
	})
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}
