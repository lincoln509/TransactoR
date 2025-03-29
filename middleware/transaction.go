package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type ctxKey string

const TxKey ctxKey = "db_tx"

func Transaction(db *gorm.DB) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tx := db.Begin()
			defer func() {
				if r := recover(); r != nil {
					tx.Rollback()
					panic(r)
				}
			}()

			ctx := context.WithValue(r.Context(), TxKey, tx)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)

			if responseWriter, ok := w.(*responseWriter); ok && responseWriter.status >= 400 {
				tx.Rollback()
			} else {
				tx.Commit()
			}
		})
	}
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func GetTx(ctx context.Context) (*gorm.DB, error) {
	tx, ok := ctx.Value(TxKey).(*gorm.DB)
	if !ok {
		return nil, fmt.Errorf("transaction non trouvée")
	}
	return tx, nil
}
