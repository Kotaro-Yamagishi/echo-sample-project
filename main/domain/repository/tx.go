package repository

import "context"

type TxRepository interface {
	Do(ctx context.Context, f func(ctx context.Context) error) error
}
