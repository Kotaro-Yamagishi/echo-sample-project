# SAP トランザクション機能について

## 概要

SAPパッケージのトランザクション機能は、データベースのトランザクションを透過的に管理し、ACID特性を保証するための仕組みです。開発者が意識する必要なく、適切なトランザクション管理が自動的に行われます。

## 目的

- **トランザクションの一貫性保証**: 全ての操作が同じトランザクション内で実行される
- **ロールバックの保証**: エラー時に全ての操作がロールバックされる
- **透過性**: 開発者がトランザクションを意識する必要がない
- **既存APIとの互換性**: リポジトリメソッドのAPIを変更せずに済む

## 核となるコンポーネント

### 1. TxContext（トランザクションコンテキスト）

```go
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
        return nil, errs.New("transaction is nil")
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
```

### 2. TxRepository（トランザクション管理）

```go
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
    // 1. トランザクション開始
    tx, err := r.db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }

    // 2. 外部キー制約の無効化（設定による）
    if !config.Env.DBEnableForeignKey {
        logger.App().WithField("component", "transaction").Debug("SET FOREIGN_KEY_CHECKS = 0")
        _, err = tx.ExecContext(ctx, "SET FOREIGN_KEY_CHECKS = 0")
        if err != nil {
            _ = tx.Rollback()
            return err
        }
    }

    // 3. TxContext作成
    tctx, err := sap.NewTxContext(ctx, tx)
    if err != nil {
        if innerErr := tx.Rollback(); innerErr != nil {
            logger.App().Errorf("failed to rollback: %v", innerErr)
        }
        return err
    }

    // 4. パニック処理
    defer func() {
        if err := recover(); err != nil {
            innerErr := tx.Rollback()
            if innerErr != nil {
                logger.App().Errorf("failed to rollback: %v", innerErr)
            }
            panic(err)
        }
    }()

    // 5. コールバック実行
    if err := f(tctx); err != nil {
        if innerErr := tx.Rollback(); innerErr != nil {
            logger.App().Errorf("failed to rollback: %v", innerErr)
        }
        return err
    }

    // 6. コミット
    if err := tx.Commit(); err != nil {
        return err
    }

    return nil
}
```

### 3. getDBExecutor（Executor選択）

```go
func getDBExecutor(ctx context.Context, db *sap.DB) boil.ContextExecutor {
    tctx := sap.FromTxContext(ctx)
    if tctx != nil {
        return tctx.GetTx()  // ← トランザクション内の場合は、トランザクションを使用
    }
    
    return db  // ← トランザクション外の場合は、通常のSAPの仕組みを使用
}
```

## 基本的な動作の流れ

### 1. トランザクション開始
```go
// ユースケース層でトランザクション開始
func (u *userUsecase) CreateUserWithProfile(ctx context.Context, input *CreateInput) error {
    return u.txRepository.Do(ctx, func(ctx context.Context) error {
        // このctxはTxContextになっている
        
        // ビジネスロジックを実行
        return nil
    })
}
```

### 2. トランザクション内でのDB操作
```go
// リポジトリ層での透過的なDB選択
func (r *userRepository) Insert(ctx context.Context, usr *entity.User) error {
    ex := getDBExecutor(ctx, r.db)  // ← 自動的に適切なExecutorが選択される
    return usr.Insert(ctx, ex, boil.Infer())
}
```

### 3. トランザクション終了
```go
// TxRepository.Do内で自動的にコミットまたはロールバック
// 成功時: tx.Commit()
// エラー時: tx.Rollback()
```

## 重要なポイント

### 1. 透過性（Transparency）

**開発者が意識する必要がない**

```go
// 透過的な実装（現在の設計）
func (r *userRepository) Insert(ctx context.Context, usr *entity.User) error {
    ex := getDBExecutor(ctx, r.db)  // ← 開発者は意識する必要なし
    return usr.Insert(ctx, ex, boil.Infer())
}

// 使用例
func (u *userUsecase) CreateUser(ctx context.Context, input *CreateInput) error {
    return u.txRepository.Do(ctx, func(ctx context.Context) error {
        // 開発者は「トランザクション内だ」と意識する必要なし
        user := &entity.User{Name: input.Name}
        err := u.userRepository.Insert(ctx, user)  // ← 普通に呼び出すだけ
        
        profile := &entity.Profile{UserID: user.ID}
        err = u.profileRepository.Insert(ctx, profile)  // ← 普通に呼び出すだけ
        
        return nil
    })
}
```

### 2. 既存APIとの互換性

**リポジトリメソッドのAPIを変更せずに済む**

```go
// 既存のリポジトリメソッド
func (r *userRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
    return entity.Users(
        entity.UserWhere.Email.EQ(email),
        entity.UserWhere.DeletedAt.IsNull(),
    ).One(ctx, r.db)  // ← 直接r.dbを使用
}

// トランザクション内でも同じAPIを使える
func (u *userUsecase) CreateUserWithCheck(ctx context.Context, input *CreateInput) error {
    return u.txRepository.Do(ctx, func(ctx context.Context) error {
        // 既存のメソッドをそのまま使える
        existingUser, err := u.userRepository.FindByEmail(ctx, input.Email)  // ← 既存API
        if existingUser != nil {
            return errors.New("既に存在するユーザーです")
        }
        
        // 新しいユーザー作成
        user := &entity.User{Email: input.Email}
        err = u.userRepository.Insert(ctx, user)  // ← 既存API
        
        return nil
    })
}
```

### 3. 同じExecutorの重要性

**トランザクション内では必ず同じExecutorを使用する必要がある**

```go
// 正しい実装（同じExecutorを使用）
func (u *userUsecase) CreateUserWithProfile(ctx context.Context, input *CreateInput) error {
    return u.txRepository.Do(ctx, func(ctx context.Context) error {
        // 全て同じトランザクション内で実行される
        
        // 1. ユーザー作成
        user := &entity.User{Name: input.Name}
        err := user.Insert(ctx, getDBExecutor(ctx, u.db))  // ← 同じExecutor
        
        // 2. プロフィール作成
        profile := &entity.Profile{UserID: user.ID}
        err = profile.Insert(ctx, getDBExecutor(ctx, u.db))  // ← 同じExecutor
        
        return nil
    })
}

// 問題のある実装（異なるExecutorを使用）
func (u *userUsecase) CreateUserWithProfile(ctx context.Context, input *CreateInput) error {
    return u.txRepository.Do(ctx, func(ctx context.Context) error {
        // 問題: 異なるExecutorを使用
        
        // 1. ユーザー作成（トランザクション内）
        user := &entity.User{Name: input.Name}
        err := user.Insert(ctx, getDBExecutor(ctx, u.db))  // ← トランザクション内
        
        // 2. プロフィール作成（トランザクション外）
        profile := &entity.Profile{UserID: user.ID}
        err = profile.Insert(ctx, u.db)  // ← 直接u.dbを使用（トランザクション外）
        
        // 結果: ユーザー作成はロールバックされるが、プロフィール作成はコミットされる
        return nil
    })
}
```

## 実際の使用例

### シナリオ1: ユーザー登録とプロフィール作成

```go
func (u *userUsecase) CreateUserWithProfile(ctx context.Context, input *CreateInput) error {
    return u.txRepository.Do(ctx, func(ctx context.Context) error {
        // 1. 重複チェック
        existingUser, err := u.userRepository.FindByEmail(ctx, input.Email)
        if existingUser != nil {
            return errors.New("既に存在するユーザーです")
        }
        
        // 2. ユーザー作成
        user := &entity.User{
            Name:  input.Name,
            Email: input.Email,
        }
        err = u.userRepository.Insert(ctx, user)
        if err != nil {
            return err
        }
        
        // 3. プロフィール作成
        profile := &entity.Profile{
            UserID: user.ID,
            Bio:    input.Bio,
        }
        err = u.profileRepository.Insert(ctx, profile)
        if err != nil {
            return err
        }
        
        // 全て同じトランザクション内で実行される
        // どれか一つでも失敗すると、全てロールバックされる
        
        return nil
    })
}
```

### シナリオ2: 注文処理と在庫管理

```go
func (u *orderUsecase) ProcessOrder(ctx context.Context, input *OrderInput) error {
    return u.txRepository.Do(ctx, func(ctx context.Context) error {
        // 1. 在庫チェック
        stock, err := u.stockRepository.FindByProductID(ctx, input.ProductID)
        if err != nil {
            return err
        }
        if stock.Quantity < input.Quantity {
            return errors.New("在庫不足")
        }
        
        // 2. 注文作成
        order := &entity.Order{
            ProductID: input.ProductID,
            Quantity:  input.Quantity,
            UserID:    input.UserID,
        }
        err = u.orderRepository.Insert(ctx, order)
        if err != nil {
            return err
        }
        
        // 3. 在庫更新
        stock.Quantity -= input.Quantity
        err = u.stockRepository.Update(ctx, stock)
        if err != nil {
            return err
        }
        
        // 4. 支払い処理
        payment := &entity.Payment{
            OrderID: order.ID,
            Amount:  order.TotalAmount,
        }
        err = u.paymentRepository.Insert(ctx, payment)
        if err != nil {
            return err
        }
        
        // 全て同じトランザクション内で実行される
        // どれか一つでも失敗すると、全てロールバックされる
        
        return nil
    })
}
```

## 技術的なポイント

### 1. Contextの型変換

```go
// TxRepository.Do内で
tctx, err := sap.NewTxContext(ctx, tx)  // ← 通常のContextをTxContextに変換

// コールバック内で
tctx := sap.FromTxContext(ctx)  // ← ContextからTxContextを取得（型アサーション）
if tctx != nil {
    // トランザクション内
    executor := tctx.GetTx()
} else {
    // トランザクション外
    executor := db
}
```

### 2. Executorの型

```go
// boil.ContextExecutorは以下のいずれか
type ContextExecutor interface {
    // トランザクション内の場合: *sql.Tx
    // トランザクション外の場合: *sap.DB
}

// 実際の型
func getDBExecutor(ctx context.Context, db *sap.DB) boil.ContextExecutor {
    tctx := sap.FromTxContext(ctx)
    if tctx != nil {
        return tctx.GetTx()  // ← *sql.Tx
    }
    return db                // ← *sap.DB
}
```

### 3. エラーハンドリング

```go
// パニック処理
defer func() {
    if err := recover(); err != nil {
        innerErr := tx.Rollback()
        if innerErr != nil {
            logger.App().Errorf("failed to rollback: %v", innerErr)
        }
        panic(err)
    }
}()

// エラー時のロールバック
if err := f(tctx); err != nil {
    if innerErr := tx.Rollback(); innerErr != nil {
        logger.App().Errorf("failed to rollback: %v", innerErr)
    }
    return err
}
```

## 設計の優れた点

### 1. 透過性
- **開発者が意識する必要がない**: トランザクション内かどうかを意識する必要なし
- **既存APIがそのまま使える**: リポジトリメソッドのAPIを変更する必要なし
- **自動的に適切なDBが選択される**: `getDBExecutor`が自動判定

### 2. 安全性
- **トランザクションの一貫性保証**: 全ての操作が同じトランザクション内で実行される
- **ロールバックの保証**: エラー時に全ての操作がロールバックされる
- **パニック処理**: 予期しないエラーでもロールバックされる

### 3. 保守性
- **シンプルなAPI**: 複雑なトランザクション管理を隠蔽
- **エラーハンドリング**: 適切なロールバック処理
- **ログ出力**: デバッグしやすいログ

## まとめ

SAPパッケージのトランザクション機能により、以下のメリットを実現しています：

1. **トランザクションの一貫性保証**: 全ての操作が同じトランザクション内で実行される
2. **ロールバックの保証**: エラー時に全ての操作がロールバックされる
3. **透過性**: 開発者がトランザクションを意識する必要がない
4. **既存APIとの互換性**: リポジトリメソッドのAPIを変更せずに済む
5. **安全性**: 適切なエラーハンドリングとパニック処理

この設計により、**複雑なトランザクション管理を隠蔽し、開発者が安全にトランザクションを使用できる**ようになっています。
