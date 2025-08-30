package sap

import (
	"context"
	"database/sql"
	"errors"
)

type TxContext interface {
	context.Context
	GetTx() *sql.Tx
}

type txContext struct {
	context.Context
	tx *sql.Tx
}

func NewTxContext(ctx context.Context, tx *sql.Tx) (TxContext, error) {
	if tx == nil {
		return nil, errors.New("transaction is nil")
	}

	return &txContext{
		Context: ctx,
		tx:      tx,
	}, nil
}

func FromTxContext(ctx context.Context) TxContext {
	tctx, ok := ctx.(TxContext)
	if !ok {
		return nil
	}
	return tctx
}

func (c *txContext) GetTx() *sql.Tx {
	return c.tx
}
