# エラー設計ドキュメント

## 概要

このプロジェクトでは、クリーンアーキテクチャに基づいた階層別エラーハンドリングを実装しています。エラーは各層で詳細を積み上げながら上層に伝播し、最終的にController層でHTTPレスポンスに変換されます。

## エラー設計の基本方針

### 1. エラーの伝播パターン

```
Repository層 → UseCase層 → Controller層
     ↓           ↓           ↓
  基本エラー → 詳細追加 → HTTPレスポンス
```

### 2. 各層の責任

#### Repository層
- データベースエラーの詳細化
- エラーメッセージにコンテキストを追加
- 基本エラーを上層に投げる

#### UseCase層
- ビジネスロジックエラーの追加
- Repository層のエラーに詳細を追加
- エラーの詳細化と上層への伝播

#### Controller層
- エラーをHTTPレスポンスに変換
- 適切なHTTPステータスコードの設定
- クライアント向けのエラーメッセージの生成

## 実装詳細

### 1. 共通レスポンス構造体

#### ErrorResponse
```go
type ErrorResponse struct {
    Success   bool      `json:"success"`
    Error     string    `json:"error"`
    Code      int       `json:"code"`
    Timestamp time.Time `json:"timestamp"`
}
```

#### SuccessResponse
```go
type SuccessResponse struct {
    Success   bool        `json:"success"`
    Data      interface{} `json:"data"`
    Timestamp time.Time   `json:"timestamp"`
}
```

### 2. 各層でのエラーハンドリング

#### Repository層
```go
func (db *CountryImpl) Insert(country entity.Country) error {
    modelCountry := &model.Country{
        Country:    string(country.Country),
        LastUpdate: country.LastUpdate,
    }
    
    if err := db.ds.Insert(modelCountry); err != nil {
        return fmt.Errorf("database error: failed to insert country: %w", err)
    }
    return nil
}
```

#### UseCase層
```go
func (s *CountryImpl) Insert(country entity.Country) error {
    if err := s.repo.Insert(country); err != nil {
        return fmt.Errorf("failed to insert country: %w", err)
    }
    return nil
}
```

#### Controller層
```go
func (c *CountryImpl) Create(ctx echo.Context) error {
    var request struct {
        Country string `json:"country"`
    }
    
    if err := ctx.Bind(&request); err != nil {
        return ctx.JSON(http.StatusBadRequest, output.NewErrorResponse(
            http.StatusBadRequest,
            "invalid request format",
        ))
    }

    country, err := entity.NewValidatedCountry(request.Country)
    if err != nil {
        return ctx.JSON(http.StatusBadRequest, output.NewErrorResponse(
            http.StatusBadRequest,
            err.Error(),
        ))
    }

    if err := c.uc.Insert(country); err != nil {
        return ctx.JSON(http.StatusInternalServerError, output.NewErrorResponse(
            http.StatusInternalServerError,
            "failed to create country",
        ))
    }

    return ctx.JSON(http.StatusCreated, output.NewSuccessResponse(country))
}
```

## HTTPステータスコード

| ステータスコード | 用途 | 例 |
|----------------|------|-----|
| 200 | 成功 | データ取得成功 |
| 201 | 作成成功 | リソース作成成功 |
| 400 | バリデーションエラー | リクエスト形式エラー |
| 500 | サーバーエラー | データベースエラー |

## エラーレスポンス例

### 成功レスポンス
```json
{
    "success": true,
    "data": [
        {
            "country_id": 1,
            "country": "Japan",
            "last_update": "2024-01-01T12:00:00Z"
        }
    ],
    "timestamp": "2024-01-01T12:00:00Z"
}
```

### エラーレスポンス
```json
{
    "success": false,
    "error": "failed to create country",
    "code": 500,
    "timestamp": "2024-01-01T12:00:00Z"
}
```

## エラーの積み上げ例

実際のエラーが発生した場合の伝播例：

```
1. データベース層: "connection refused"
2. Repository層: "database error: failed to insert country: connection refused"
3. UseCase層: "failed to insert country: database error: failed to insert country: connection refused"
4. Controller層: HTTP 500 "failed to create country"
```

## 設計のメリット

### 1. デバッグしやすい
- エラーの発生源が明確
- エラーの伝播経路が追跡可能
- 各層での詳細情報が保持される

### 2. 一貫性
- すべてのエンドポイントで同じレスポンス形式
- 統一されたエラーハンドリングパターン
- 予測可能なAPIレスポンス

### 3. 保守性
- 各層の責任が明確
- エラーハンドリングの変更が容易
- テストしやすい構造

## 今後の改善案

### 1. カスタムエラー型の導入
```go
type AppError struct {
    Code    string
    Message string
    Cause   error
    Layer   string
}
```

### 2. エラーログの強化
- 構造化ログの導入
- エラー詳細のログ出力
- トレーサビリティの向上

### 3. エラーコードの詳細化
- より細かいエラーコードの定義
- クライアント向けエラーコード
- 国際化対応

## ベストプラクティス

### 1. エラーメッセージの設計
- 技術的な詳細は内部で保持
- クライアント向けは一般化
- セキュリティを考慮した情報開示

### 2. ログ出力
- エラーの詳細はログに出力
- 機密情報は適切にマスク
- 構造化ログの活用

### 3. テスト
- エラーケースのテスト
- エラーレスポンスの検証
- エラー伝播のテスト

## まとめ

このエラー設計により、以下の目標を達成しています：

1. **明確な責任分離**: 各層が適切なエラーハンドリングを担当
2. **詳細なエラー情報**: デバッグに必要な情報を保持
3. **一貫したAPI**: クライアントが予測可能なレスポンス
4. **保守性**: 変更しやすく、テストしやすい構造

この設計は、クリーンアーキテクチャの原則に従い、エラーの詳細を積み上げながら上層に伝播し、最終的に適切なHTTPレスポンスに変換するパターンを採用しています。
