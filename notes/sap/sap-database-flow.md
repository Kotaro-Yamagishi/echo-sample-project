# SAP データベース接続フローについて

## 概要

SAPパッケージにおける、DBに接続してクエリを叩いて値を取ってくるまでの詳細なフローを説明します。このフローにより、開発者は通常通りにコードを書き、適切なDBが自動的に選択されるようになっています。

## 基本的なフロー

### 1. リクエスト開始

```go
// ユーザーがAPIリクエストを送信
GET /api/v1/users/123

// コントローラーでリクエストを受け取り
func (ctr *UserController) Get(c echo.Context) error {
    ctx := newContext(c)  // ← sap.Contextを作成
    
    user, err := ctr.userRepository.FindByID(ctx, userID)  // ← リポジトリを呼び出し
    return c.JSON(http.StatusOK, user)
}
```

### 2. リポジトリ層での処理

```go
// リポジトリ層
func (r *userRepository) FindByID(ctx context.Context, id uint) (*entity.User, error) {
    // 1. Executorを取得
    ex := getDBExecutor(ctx, r.db)  // ← ここでDB選択が行われる
    
    // 2. SQLBoilerでクエリ実行
    return entity.Users(
        entity.UserWhere.ID.EQ(id),
        entity.UserWhere.DeletedAt.IsNull(),
    ).One(ctx, ex)  // ← 実際のDBクエリが実行される
}
```

### 3. Executor選択の詳細

```go
func getDBExecutor(ctx context.Context, db *sap.DB) boil.ContextExecutor {
    // 1. トランザクション内かどうかを判定
    tctx := sap.FromTxContext(ctx)
    if tctx != nil {
        return tctx.GetTx()  // ← トランザクション内: *sql.Tx
    }
    
    // 2. トランザクション外の場合
    return db  // ← 通常のDB: *sap.DB
}
```

### 4. SAP層でのDB選択

```go
// トランザクション外の場合（*sap.DB）
func (db *DB) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
    // 1. recordsModifiedフラグをチェック
    if db.recordsModified(ctx) {
        return db.Primary().QueryContext(ctx, query, args...)  // ← Primary DB
    }
    
    // 2. recordsModified = false の場合
    return db.ReadReplica().QueryContext(ctx, query, args...)  // ← Read Replica
}

// Read Replica の選択
func (db *DB) ReadReplica() *sql.DB {
    return db.connections[db.rotate(len(db.connections))]  // ← ラウンドロビン
}
```

## 具体的なフロー例

### シナリオ1: 通常の読み取り（トランザクション外）

```go
// 1. リクエスト開始
GET /api/v1/users

// 2. コントローラー
func (ctr *UserController) GetAll(c echo.Context) error {
    ctx := newContext(c)  // ← sap.Context作成（recordsModified = false）
    
    users, err := ctr.userRepository.FindAll(ctx)  // ← リポジトリ呼び出し
    return c.JSON(http.StatusOK, users)
}

// 3. リポジトリ
func (r *userRepository) FindAll(ctx context.Context) ([]*entity.User, error) {
    ex := getDBExecutor(ctx, r.db)  // ← tctx = nil なので r.db を返す
    
    return entity.Users().All(ctx, ex)  // ← SQLBoilerでクエリ実行
}

// 4. SAP層
func (db *DB) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
    if db.recordsModified(ctx) {  // ← recordsModified = false
        return db.Primary().QueryContext(ctx, query, args...)
    }
    return db.ReadReplica().QueryContext(ctx, query, args...)  // ← Read Replica を使用
}

// 5. Read Replica 選択
func (db *DB) ReadReplica() *sql.DB {
    return db.connections[db.rotate(len(db.connections))]  // ← ラウンドロビンで選択
}

// 6. 実際のDBクエリ実行
// → Read Replica 1, 2, 3 のいずれかでクエリが実行される
```

### シナリオ2: 書き込み処理（トランザクション外）

```go
// 1. リクエスト開始
POST /api/v1/users

// 2. コントローラー
func (ctr *UserController) Create(c echo.Context) error {
    ctx := newContext(c)  // ← sap.Context作成（recordsModified = false）
    
    user, err := ctr.userRepository.Insert(ctx, input)  // ← リポジトリ呼び出し
    return c.JSON(http.StatusCreated, user)
}

// 3. リポジトリ
func (r *userRepository) Insert(ctx context.Context, usr *entity.User) error {
    ex := getDBExecutor(ctx, r.db)  // ← tctx = nil なので r.db を返す
    
    return usr.Insert(ctx, ex, boil.Infer())  // ← SQLBoilerでクエリ実行
}

// 4. SAP層（ExecContext）
func (db *DB) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
    db.setRecordsModified(ctx)  // ← recordsModified = true に設定
    return db.Primary().ExecContext(ctx, query, args...)  // ← Primary DB で実行
}

// 5. 実際のDBクエリ実行
// → Primary DB でINSERTクエリが実行される
```

### シナリオ3: トランザクション内の処理

```go
// 1. リクエスト開始
POST /api/v1/users/with-profile

// 2. コントローラー
func (ctr *UserController) CreateWithProfile(c echo.Context) error {
    ctx := newContext(c)
    
    return ctr.userUsecase.CreateWithProfile(ctx, input)  // ← ユースケース呼び出し
}

// 3. ユースケース
func (u *userUsecase) CreateWithProfile(ctx context.Context, input *CreateInput) error {
    return u.txRepository.Do(ctx, func(ctx context.Context) error {  // ← トランザクション開始
        // このctxはTxContextになっている
        
        user, err := u.userRepository.Insert(ctx, input.User)  // ← リポジトリ呼び出し
        profile, err := u.profileRepository.Insert(ctx, input.Profile)  // ← リポジトリ呼び出し
        
        return nil
    })
}

// 4. トランザクション開始
func (r *TxRepository) Do(ctx context.Context, f func(ctx context.Context) error) error {
    tx, err := r.db.BeginTx(ctx, nil)  // ← トランザクション開始
    tctx, err := sap.NewTxContext(ctx, tx)  // ← TxContext作成
    
    if err := f(tctx); err != nil {  // ← TxContextを渡してコールバック実行
        tx.Rollback()
        return err
    }
    
    return tx.Commit()
}

// 5. リポジトリ（トランザクション内）
func (r *userRepository) Insert(ctx context.Context, usr *entity.User) error {
    ex := getDBExecutor(ctx, r.db)  // ← tctx != nil なので tctx.GetTx() を返す
    
    return usr.Insert(ctx, ex, boil.Infer())  // ← トランザクション内でクエリ実行
}

// 6. Executor選択
func getDBExecutor(ctx context.Context, db *sap.DB) boil.ContextExecutor {
    tctx := sap.FromTxContext(ctx)  // ← TxContextを取得
    if tctx != nil {
        return tctx.GetTx()  // ← トランザクションを返す
    }
    return db
}

// 7. 実際のDBクエリ実行
// → 同じトランザクション内でクエリが実行される
```

## フローのまとめ

### 1. リクエスト受信
```
HTTPリクエスト → コントローラー → sap.Context作成
```

### 2. リポジトリ呼び出し
```
リポジトリメソッド → getDBExecutor() → Executor選択
```

### 3. Executor選択
```
トランザクション内？ → Yes: *sql.Tx / No: *sap.DB
```

### 4. DB選択（*sap.DBの場合）
```
recordsModified = true？ → Yes: Primary DB / No: Read Replica
```

### 5. クエリ実行
```
SQLBoiler → 実際のDB接続 → クエリ実行 → 結果取得
```

## 詳細なフロー図

### 通常の読み取り処理
```
1. HTTPリクエスト (GET /api/v1/users)
   ↓
2. コントローラー (newContext() → sap.Context作成)
   ↓
3. リポジトリ (getDBExecutor() → *sap.DB取得)
   ↓
4. SAP層 (recordsModified = false → Read Replica選択)
   ↓
5. ラウンドロビン (Read Replica 1,2,3 から選択)
   ↓
6. 実際のDBクエリ実行
   ↓
7. 結果取得・返却
```

### 書き込み処理
```
1. HTTPリクエスト (POST /api/v1/users)
   ↓
2. コントローラー (newContext() → sap.Context作成)
   ↓
3. リポジトリ (getDBExecutor() → *sap.DB取得)
   ↓
4. SAP層 (ExecContext → recordsModified = true設定)
   ↓
5. Primary DB選択
   ↓
6. 実際のDBクエリ実行
   ↓
7. 結果取得・返却
```

### トランザクション内の処理
```
1. HTTPリクエスト (POST /api/v1/users/with-profile)
   ↓
2. コントローラー (newContext() → sap.Context作成)
   ↓
3. ユースケース (txRepository.Do() → トランザクション開始)
   ↓
4. TxContext作成 (NewTxContext() → トランザクション情報をContextに埋め込み)
   ↓
5. リポジトリ (getDBExecutor() → *sql.Tx取得)
   ↓
6. 実際のDBクエリ実行（トランザクション内）
   ↓
7. コミット・ロールバック
   ↓
8. 結果取得・返却
```

## 重要なポイント

### 1. 自動的なDB選択
- **トランザクション内**: 必ずPrimary DBのトランザクションを使用
- **書き込み処理**: 自動的にPrimary DBを使用
- **読み取り処理**: recordsModifiedフラグに応じてDB選択

### 2. 透過性
- **開発者は意識する必要なし**: 適切なDBが自動選択される
- **既存APIがそのまま使える**: リポジトリメソッドのAPIは変更不要

### 3. パフォーマンス最適化
- **読み取り専用**: Read Replicaで負荷分散
- **書き込み処理**: Primary DBで一貫性保証

### 4. データ整合性
- **トランザクション内**: 全ての操作が同じトランザクションで実行
- **書き込み後**: 以降の読み取りはPrimary DBから最新データを取得

## 技術的な詳細

### Contextの状態管理
```go
// リクエスト開始時
ctx := &Context{Context: req.Context(), recordsModified: false}

// 書き込み処理時
func (c *Context) SetRecordsModified() {
    c.recordsModified = true
}

// リクエスト終了時
// Contextが破棄されるので状態もリセット
```

### Executorの型
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

### ラウンドロビン負荷分散
```go
func (db *DB) rotate(n int) int {
    if n <= 1 {
        return 0
    }
    return int(1 + (atomic.AddUint64(&db.counter, 1) % uint64(n-1)))
}

// 例：3つのRead Replicaがある場合
// リクエスト1: Read Replica 1 を使用
// リクエスト2: Read Replica 2 を使用  
// リクエスト3: Read Replica 3 を使用
// リクエスト4: Read Replica 1 を使用（ローテーション）
```

## まとめ

SAPパッケージのデータベース接続フローにより、以下のメリットを実現しています：

1. **自動的なDB選択**: 適切なDBが自動的に選択される
2. **透過性**: 開発者が意識する必要がない
3. **パフォーマンス最適化**: 読み取り専用処理はRead Replicaで負荷分散
4. **データ整合性**: 書き込み処理はPrimary DBで一貫性保証
5. **トランザクション対応**: トランザクション内では適切なExecutorが選択される

このフローにより、**開発者は通常通りにコードを書き、複雑なDB選択ロジックを意識することなく、効率的で安全なデータベースアクセス**を実現できます。
