package dbctx

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

type contextKey string

const (
	TxKey contextKey = "db_tx"
)

func TxFromContext(ctx context.Context) (*gorm.DB, error) {
	tx, ok := ctx.Value(TxKey).(*gorm.DB)
	if !ok || tx == nil {
		return nil, fmt.Errorf("transaction non trouvée dans le contexte")
	}
	return tx, nil
}
