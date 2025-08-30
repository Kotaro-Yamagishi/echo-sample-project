# SAP (Sharded Access Pattern) データベース仕組みについて

## 概要

SAPパッケージは、データベースの読み取り専用レプリカ（Read Replica）とプライマリDB（Primary DB）を適切に使い分けるための仕組みです。シンプルで効率的な設計により、パフォーマンス向上とデータ整合性の両方を実現しています。

## 目的

- **パフォーマンス向上**: 読み取り専用の処理は Read Replica を使用
- **データ整合性**: 変更後の読み取りは Primary DB から最新データを取得
- **自動制御**: 開発者が意識せずとも適切なDBが選択される
- **柔軟性**: HTTPヘッダーで手動制御も可能
- **負荷分散**: 複数のRead Replica間でのロードバランシング

## DB構成

### 接続設定
```go
type DB struct {
    connections []*sql.DB  // 複数のDB接続を保持
    counter     uint64     // ラウンドロビン用カウンター
}
```

### データソース設定
```go
// セミコロンで区切って複数のDBを指定
dataSourceNames := "primary_db_url;read_replica1_url;read_replica2_url;read_replica3_url"

// 例：
// connections[0]: Primary DB (書き込み専用)
// connections[1]: Read Replica 1 (読み取り専用)
// connections[2]: Read Replica 2 (読み取り専用)
// connections[3]: Read Replica 3 (読み取り専用)
```

## 核となる仕組み

### 1. Context での状態管理
```go
type Context struct {
    context.Context
    recordsModified bool  // リクエスト単位で状態を管理
}

func (c *Context) SetRecordsModified() {
    c.recordsModified = true
}

func (c *Context) RecordsModified() bool {
    return c.recordsModified
}
```

### 2. HTTPヘッダーによる制御
```go
const HeaderKeyRecordModified = "X-Record-Modified"

func (c *Context) setRecordModifiedFromHeader(req *http.Request) {
    if req.Header.Get(HeaderKeyRecordModified) == "on" {
        c.SetRecordsModified()
    }
}
```

## 実際に使用されているメソッド

### 1. ExecContext（書き込み処理）
```go
func (db *DB) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
    db.setRecordsModified(ctx)  // ← ここで自動的に recordsModified = true に設定
    return db.Primary().ExecContext(ctx, query, args...)  // ← 必ずPrimary DBで実行
}
```

### 2. QueryContext（読み取り処理）
```go
func (db *DB) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
    if db.recordsModified(ctx) {
        return db.Primary().QueryContext(ctx, query, args...)  // ← Primary DB
    }
    return db.ReadReplica().QueryContext(ctx, query, args...)  // ← Read Replica
}
```

### 3. QueryRowContext（単一行読み取り）
```go
func (db *DB) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
    if db.recordsModified(ctx) {
        return db.Primary().QueryRowContext(ctx, query, args...)  // ← Primary DB
    }
    return db.ReadReplica().QueryRowContext(ctx, query, args...)  // ← Read Replica
}
```

### 4. BeginTx（トランザクション開始）
```go
func (db *DB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
    db.setRecordsModified(ctx)  // ← トランザクション開始時もフラグを設定
    return db.Primary().BeginTx(ctx, opts)  // ← 必ずPrimary DBでトランザクション開始
}
```

## 未使用のメソッド

### PrepareContext（プリペアドステートメント）
```go
func (db *DB) PrepareContext(ctx context.Context, query string) (Stmt, error) {
    if db.recordsModified(ctx) {
        return db.Primary().PrepareContext(ctx, query)
    }

    stmts := make([]*sql.Stmt, len(db.connections))
    err := scatter(len(db.connections), func(i int) (err error) {
        stmts[i], err = db.connections[i].PrepareContext(ctx, query)
        return err
    })

    return &stmt{ctx: ctx, db: db, stmts: stmts}, nil
}
```

**使用されていない理由:**
- SQLBoilerがプリペアドステートメントを使用しない
- 直接クエリ実行（QueryContext/ExecContext）を採用
- 動的クエリ生成を重視

## ラウンドロビン負荷分散

### Read Replica の選択
```go
func (db *DB) ReadReplica() *sql.DB {
    return db.connections[db.rotate(len(db.connections))]
}

func (db *DB) rotate(n int) int {
    if n <= 1 {
        return 0
    }
    return int(1 + (atomic.AddUint64(&db.counter, 1) % uint64(n-1)))  // ← 1から始まる（Primary DBを除く）
}
```

### 負荷分散の例
```go
// 複数のRead Replica間でローテーション
// リクエスト1: Read Replica 1 を使用
// リクエスト2: Read Replica 2 を使用
// リクエスト3: Read Replica 3 を使用
// リクエスト4: Read Replica 1 を使用（ローテーション）
```

## 基本的な動作の流れ

### 1. 初期状態（読み取り専用）
```go
// ユーザーがGETリクエストを送信
GET /api/v1/users/123

// この時点では recordsModified = false
// → Read Replica からデータを取得
```

### 2. データ変更が発生
```go
// ユーザーがPOST/PATCH/DELETEリクエストを送信
POST /api/v1/users
{
  "name": "田中太郎",
  "email": "tanaka@example.com"
}

// この時点で recordsModified = true に設定される
// → 以降のクエリは Primary DB から取得
```

### 3. 変更後の読み取り
```go
// 同じリクエスト内で再度データを取得
GET /api/v1/users/123

// recordsModified = true なので
// → Primary DB から最新データを取得
```

## 実際の使用例

### シナリオ1: ユーザー登録後の確認
```go
func (ctr *UserController) Create(c echo.Context) error {
    // 1. ユーザー登録（POST）
    user := &entity.User{Name: "田中太郎", Email: "tanaka@example.com"}
    err := ctr.userRepository.Insert(ctx, user)  // ← ExecContext が呼ばれ、recordsModified = true に設定
    
    // 2. 登録したユーザー情報を取得（同じリクエスト内）
    createdUser, err := ctr.userRepository.FindByID(ctx, user.ID)  // ← Primary DB から最新データを取得
    
    return c.JSON(http.StatusCreated, createdUser)
}
```

### シナリオ2: 通常の読み取り
```go
func (ctr *UserController) GetAll(c echo.Context) error {
    // 単純なデータ取得
    users, err := ctr.userRepository.FindAll(ctx)  // ← recordsModified = false なので Read Replica から取得
    return c.JSON(http.StatusOK, users)
}
```

### シナリオ3: 複雑な処理での自動最適化
```go
func (ctr *OrderController) ProcessOrder(c echo.Context) error {
    // 1. 初期状態: recordsModified = false
    
    // 2. 在庫チェック（Read Replica から取得）
    stock, err := ctr.stockRepository.FindByProductID(ctx, productID)
    if stock.Quantity < orderQuantity {
        return errors.New("在庫不足")
    }
    
    // 3. 注文作成（Primary DB に書き込み）
    order := &entity.Order{ProductID: productID, Quantity: orderQuantity}
    err = ctr.orderRepository.Insert(ctx, order)  // ← ここで recordsModified = true に設定
    
    // 4. 在庫更新（Primary DB に書き込み）
    stock.Quantity -= orderQuantity
    err = ctr.stockRepository.Update(ctx, stock)
    
    // 5. 更新後の在庫確認（Primary DB から取得）
    updatedStock, err := ctr.stockRepository.FindByProductID(ctx, productID)  // ← Primary DB から取得
    
    return c.JSON(http.StatusCreated, order)
}
```

## HTTPヘッダーの必要性について

### 基本的には不要
- **自動制御**: INSERT/UPDATE/DELETE時に自動的にPrimary DBが選択される
- **フラグ管理**: `recordsModified`フラグにより、以降のクエリは自動的にPrimary DBを使用
- **安全性**: データ整合性が自動的に保証される

### ヘッダーが必要な特殊なケース
- **読み取り処理でも最新データが必要な場合**
- **レプリケーション遅延が問題になる場合**
- **金融取引など、データの即時性が重要な場合**

### 使用例
```go
// 通常のCRUD操作: ヘッダーなし（自動制御に任せる）
GET /api/v1/users     // Read Replica から取得
POST /api/v1/users    // 自動的にPrimary DB使用

// 特殊な読み取り処理: ヘッダーあり
GET /api/v1/users/123?realtime=true  // ヘッダー: X-Record-Modified: on
```

## 重要なポイント

### 1. リクエスト単位での状態管理
- **リクエスト開始時**: `recordsModified = false`
- **最初の書き込み操作**: `recordsModified = true` に設定
- **以降の全クエリ**: Primary DB から取得
- **リクエスト終了時**: Context が破棄されるので状態もリセット

### 2. 自動切り替えのタイミング
- **INSERT/UPDATE/DELETE**: 自動的に `recordsModified = true` に設定
- **SELECT**: `recordsModified` の状態に応じてDBを選択

### 3. ヘッダーによる制御
- **X-Record-Modified: on**: 強制的に Primary DB を使用
- **ヘッダーなし**: 動的にDBを選択（初期は Read Replica）

## 設計の優れた点

### 1. シンプルさ
```go
// たったこれだけの実装で完結
ExecContext() → フラグ立てる
QueryContext() → フラグ見てDB選択

// 結果として
- パフォーマンス向上 ✓
- データ整合性保証 ✓  
- 開発者負荷軽減 ✓
- 保守性向上 ✓
```

### 2. 安全性の確保
```go
// 例：ヘッダーでRead指定しても、実際の更新処理は必ずPrimary DBで実行される
func (ctr *UserController) UpdateUser(c echo.Context) error {
    // ヘッダー: X-Record-Modified: off でも
    user.Name = "新しい名前"
    err = ctr.userRepository.Update(ctx, user)  // ← 必ずPrimary DBで実行される
    
    return c.JSON(http.StatusOK, user)
}
```

### 3. パフォーマンス最適化
```go
// 基本: Read Replica で負荷分散（高速・軽量）
GET /api/v1/users     // Read Replica 1
GET /api/v1/products  // Read Replica 2
GET /api/v1/orders    // Read Replica 3

// 例外: 重要な読み取り処理のみ Primary DB を使用（最新データ）
POST /api/v1/users    // Primary DB
GET /api/v1/balance   // Primary DB（ヘッダーあり）
```

### 4. 自動最適化
```go
// 開発者が意識せずとも適切なDBが選択される
func (ctr *UserController) Create(c echo.Context) error {
    // 1. 初期状態: Read Replica から取得
    existingUser, err := ctr.userRepository.FindByEmail(ctx, email)  // ← Read Replica
    
    // 2. ユーザー作成: 自動的にPrimary DB使用
    user := &entity.User{Email: email}
    err = ctr.userRepository.Insert(ctx, user)  // ← Primary DB
    
    // 3. 以降: 自動的にPrimary DB使用
    createdUser, err := ctr.userRepository.FindByID(ctx, user.ID)  // ← Primary DB
    
    return c.JSON(http.StatusCreated, createdUser)
}
```

## 技術的なポイント

### 1. Context ベース
- リクエスト単位で状態を管理
- リクエスト終了時に自動的にリセット

### 2. 自動検知
- INSERT/UPDATE/DELETE 操作を自動検知
- SQLBoilerのORMが自動的に適切なメソッドを呼び出し

### 3. 手動制御
- HTTPヘッダーでの明示的な制御も可能
- 必要に応じて柔軟に制御可能

### 4. 分散対応
- 複数の Read Replica 間でのロードバランシング
- プライマリDBの負荷軽減

## まとめ

SAPパッケージの仕組みにより、以下のメリットを実現しています：

1. **パフォーマンス向上**: 読み取り専用処理は Read Replica を使用
2. **データ整合性**: 変更後の読み取りは Primary DB から最新データを取得
3. **自動制御**: 開発者が意識せずとも適切なDBが選択される
4. **柔軟性**: 必要に応じてHTTPヘッダーで明示的に制御可能
5. **負荷分散**: 複数のRead Replica間でのロードバランシング

この設計により、**シンプルで効果的な実装**でデータベースの負荷分散とデータ整合性の両方を効率的に実現しています。まさに「**シンプル イズ ベスト**」の良い例です。
