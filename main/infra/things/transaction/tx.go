package transaction

import (
	"context"
	"echoProject/domain/repository"
	"echoProject/infra/things/sap"
)

type TxRepository struct {
	db *sap.DB
}

func NewTxRepository(db *sap.DB) repository.TxRepository {
	return &TxRepository{
		db: db,
	}
}

// Do は、トランザクションを実行します。
func (r *TxRepository) Do(ctx context.Context, f func(ctx context.Context) error) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// if !config.Env.DBEnableForeignKey {
	// 	logger.App().WithField("component", "transaction").Debug("SET FOREIGN_KEY_CHECKS = 0")
	// 	_, err = tx.ExecContext(ctx, "SET FOREIGN_KEY_CHECKS = 0")
	// 	if err != nil {
	// 		_ = tx.Rollback()
	// 		return err
	// 	}
	// }

	// tctx, err := sap.NewTxContext(ctx, tx)
	// if err != nil {
	// 	if innerErr := tx.Rollback(); innerErr != nil {
	// 		logger.App().Errorf("failed to rollback: %v", innerErr)
	// 	}
	// 	return err
	// }

	// defer func() {
	// 	if err := recover(); err != nil {
	// 		innerErr := tx.Rollback()
	// 		if innerErr != nil {
	// 			logger.App().Errorf("failed to rollback: %v", innerErr)
	// 		}
	// 		panic(err)
	// 	}
	// }()

	// // トランザクション用のコンテキストを渡してコールバックを実行
	// if err := f(tctx); err != nil {
	// 	if innerErr := tx.Rollback(); innerErr != nil {
	// 		logger.App().Errorf("failed to rollback: %v", innerErr)
	// 	}
	// 	return err
	// }

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
