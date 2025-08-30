# Context の種類と違い

このプロジェクトでは5種類のContextが使われています。それぞれ何をするか、いつ使うかを説明します。

## 1. `context.Context` (Go言語の基本)

**パッケージ**: `context`
**何をするか**: リクエストの「状況」を管理する

### 具体的に何ができる？
- **キャンセル**: ユーザーがページを離れたら処理を止める
- **タイムアウト**: 5秒で処理が終わらなかったら止める
- **値の受け渡し**: リクエストIDやユーザー情報を関数間で渡す

### いつ使う？
- データベースのクエリ実行時
- 外部APIの呼び出し時
- 長時間の処理を実行する時

```go
// 例：5秒でタイムアウトする処理
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
db.QueryContext(ctx, "SELECT * FROM users")
```

## 2. `echo.Context` (Webフレームワーク)

**パッケージ**: `github.com/labstack/echo/v4`
**何をするか**: Webアプリの「HTTPリクエスト」を処理する

### 具体的に何ができる？
- **URLパラメータ取得**: `/users/123` の `123` を取得
- **フォームデータ取得**: ユーザーが送信した名前やメールを取得
- **JSONレスポンス返却**: データをJSON形式でブラウザに返す
- **HTTPステータス設定**: 200 OK、404 Not Foundなどを設定

### いつ使う？
- Webページのリクエストを受け取る時
- APIのエンドポイントを処理する時
- ブラウザにデータを返す時

```go
// 例：ユーザー情報を取得してJSONで返す
func getUser(c echo.Context) error {
    id := c.Param("id")           // URLからID取得
    user := getUserFromDB(id)     // DBから取得
    return c.JSON(200, user)      // JSONで返す
}
```

## 3. `Context` (カスタム型 - レコード変更追跡)

**パッケージ**: `main/infra/things/sap`
**何をするか**: データベースの「レコードが変更されたかどうか」を追跡する

### 具体的に何ができる？
- **変更フラグ管理**: データベースのレコードが更新されたかどうかを記録
- **ヘッダー読み取り**: ブラウザから送られてきた `X-Record-Modified` ヘッダーを読み取る
- **状態確認**: レコードが変更されたかどうかを確認

### いつ使う？
- データベースの更新処理時
- 変更があった場合のみログを出力したい時
- キャッシュの更新が必要かどうか判断する時

```go
// 例：レコード変更を追跡
func updateUser(ctx *Context, user User) error {
    // データベース更新
    db.Update(user)
    
    // 変更があったかチェック
    if ctx.RecordsModified() {
        log.Println("ユーザー情報が更新されました")
    }
}
```

## 4. `entity.Context` (ドメイン層 - テスト用抽象化)

**パッケージ**: `main/domain/entity`
**何をするか**: テストを簡単にするための「抽象化インターフェース」

### 具体的に何ができる？
- **Echoフレームワークから独立**: Echoを使わなくてもテストできる
- **最小限の機能**: 必要な機能だけを定義
- **モック作成**: テスト用の偽物を作りやすい

### いつ使う？
- ビジネスロジックのテスト時
- Echoフレームワークを使いたくない時
- 単体テストを書く時

```go
// 例：テスト用のモック
type MockContext struct {
    params map[string]string
}

func (m *MockContext) Param(key string) string {
    return m.params[key]
}

// テストで使う
mockCtx := &MockContext{params: map[string]string{"id": "123"}}
```

## 5. `TxContext` (データベーストランザクション用)

**パッケージ**: `main/infra/things/sap`
**何をするか**: データベースの「トランザクション」を管理する

### 具体的に何ができる？
- **トランザクション保持**: データベースの処理をまとめて実行/取り消し
- **一貫性保証**: 複数の処理を全部成功させるか、全部失敗させるか
- **ロールバック**: エラーが起きたら全ての変更を取り消す

### いつ使う？
- 複数のテーブルを同時に更新する時
- お金の送金処理（出金と入金を同時に）
- 在庫管理（在庫減らして売上記録）

```go
// 例：銀行送金処理
func transferMoney(ctx TxContext, from, to string, amount int) error {
    tx := ctx.GetTx()
    
    // 出金処理
    tx.Exec("UPDATE accounts SET balance = balance - ? WHERE id = ?", amount, from)
    
    // 入金処理
    tx.Exec("UPDATE accounts SET balance = balance + ? WHERE id = ?", amount, to)
    
    // 両方成功したらコミット、失敗したらロールバック
    return nil
}
```

## 使い分けのまとめ

| 場面 | 使うContext | 理由 |
|------|-------------|------|
| **Webページのリクエスト処理** | `echo.Context` | HTTPの情報（URL、フォームデータ）を取得できる |
| **データベースのクエリ実行** | `context.Context` | タイムアウトやキャンセル機能が必要 |
| **レコードの変更を追跡** | `Context` (カスタム) | データが更新されたかどうかを記録 |
| **複数のDB処理をまとめて実行** | `TxContext` | 全部成功か全部失敗かの保証が必要 |
| **ビジネスロジックのテスト** | `entity.Context` | Echoフレームワークに依存しない |

## 実際の流れ

```
1. ブラウザからリクエスト → echo.Context
2. レコード変更追跡が必要 → Context (カスタム)
3. データベース処理 → context.Context + TxContext
4. テスト時 → entity.Context
```

## 簡単に言うと

- **`context.Context`**: 処理の「状況」管理（タイムアウト、キャンセル）
- **`echo.Context`**: Webの「リクエスト」処理（URL、フォーム）
- **`Context`**: データの「変更」追跡（更新されたかチェック）
- **`entity.Context`**: テストの「抽象化」（フレームワークに依存しない）
- **`TxContext`**: データベースの「一貫性」保証（全部成功か全部失敗）

## 中身の比較表

| Context | 型 | 持っている情報 | 主な機能 |
|---------|----|----------------|----------|
| **`context.Context`** | インターフェース | タイムアウト、キャンセル、値 | `Deadline()`, `Done()`, `Value()` |
| **`echo.Context`** | インターフェース | HTTPリクエスト/レスポンス | `Param()`, `JSON()`, `Bind()` |
| **`Context`** (カスタム) | 構造体 | `context.Context` + `recordsModified` | `RecordsModified()`, `SetRecordsModified()` |
| **`entity.Context`** | インターフェース | HTTP処理の抽象化 | `Param()`, `Bind()`, `Status()`, `JSON()` |
| **`TxContext`** | インターフェース | `context.Context` + `*sql.Tx` | `GetTx()` |

## 実際の構造比較

### 1. `context.Context` (標準)
```go
type Context interface {
    Deadline() (deadline time.Time, ok bool)
    Done() <-chan struct{}
    Err() error
    Value(key interface{}) interface{}
}
```

### 2. `Context` (カスタム - レコード変更追跡)
```go
type Context struct {
    context.Context        // 標準Contextを埋め込み
    recordsModified bool   // 独自のフラグ
}
```

### 3. `TxContext` (トランザクション)
```go
type TxContext interface {
    context.Context        // 標準Contextを埋め込み
    GetTx() *sql.Tx        // データベーストランザクション
}
```

### 4. `entity.Context` (ドメイン抽象化)
```go
type Context interface {
    Param(string) string           // URLパラメータ取得
    Bind(interface{}) error        // リクエストデータバインド
    Status(int)                    // HTTPステータス設定
    JSON(int, interface{})         // JSONレスポンス
}
```

## なぜ同じ名前なのか？

1. **共通の概念**: すべて「処理の文脈・状況」を表す
2. **Go言語の慣習**: Contextパターンが広く使われている
3. **拡張性**: 基本のContextに機能を追加する設計
4. **型安全性**: 用途ごとに型を分けることで間違いを防ぐ

**結論**: 同じ「Context」という名前でも、用途によって中身が全然違います！
