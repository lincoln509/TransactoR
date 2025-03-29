package ctx

// package dbcontext

import (
	"TransactoR/middleware"
	"context"
	"fmt"

	// "TransactoR/middleware"

	"gorm.io/gorm"
)

// type contextKey string

// const TxKey contextKey = "db_tx"

func TxFromContext(ctx context.Context) (*gorm.DB, error) {
	tx, ok := ctx.Value(middleware.TxKey).(*gorm.DB)
	if !ok || tx == nil {
		return nil, fmt.Errorf("transaction non trouvée dans le contexte")
	}
	return tx, nil
}

func SetTx(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, middleware.TxKey, tx)
}

func GetTx(ctx context.Context) (*gorm.DB, error) {
	tx, ok := ctx.Value(middleware.TxKey).(*gorm.DB)
	if !ok {
		return nil, fmt.Errorf("transaction non trouvée")
	}
	return tx, nil
}
