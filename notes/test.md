# テスト設計

## 概要

このプロジェクトでは、Goのクリーンアーキテクチャに基づいたテスト設計を採用しています。テストは各レイヤー（Controller、UseCase、Repository）で独立して実行され、モックを使用して依存関係を分離しています。

## テスト構造

### ディレクトリ構成

```
main/
├── testlib/                    # テスト用ライブラリ
│   ├── mocks/                  # モック実装
│   │   ├── country_repository.go
│   │   └── city_repository.go
│   ├── helpers.go              # 共通ヘルパー関数
│   └── example_test.go         # 使用例
├── usecase/
│   └── country/
│       └── country_test.go     # UseCase層のテスト
└── domain/
    └── types/
        └── country_test.go     # ドメイン層のテスト
```

## テストライブラリ設計

### 1. モック実装

#### 特徴
- **ビルダーパターン**: メソッドチェーンで設定可能
- **再利用性**: 複数のテストで共通利用
- **型安全性**: コンパイル時にエラー検出

#### 使用例

```go
// 基本的なモック作成
mockRepo := mocks.NewMockCountryRepository()

// テストデータの設定
mockRepo.WithCountries(mocks.CreateTestCountries())

// エラーの設定
mockRepo.WithSelectError(errors.New("database error"))

// ビルダーパターンでの設定
mockRepo := mocks.NewMockCountryRepository().
    WithCountries(testCountries).
    WithInsertError(errors.New("insert failed"))
```

### 2. ヘルパー関数

#### 共通アサーション関数

```go
// エラーの存在をアサート
func AssertError(t *testing.T, err error, expected bool)

// エラーメッセージをアサート
func AssertErrorMessage(t *testing.T, err error, expectedMsg string)

// スライスの長さをアサート
func AssertSliceLength(t *testing.T, slice interface{}, expected int)
```

#### テストデータ作成関数

```go
// テスト用の固定時間を作成
func CreateTestTime() time.Time

// テスト用のエラーを作成
func CreateTestError(message string) error

// テスト用のCountryエンティティを作成
func CreateTestCountry(id uint16, name string) entity.Country
```

## テストパターン

### 1. テーブル駆動テスト

#### 適している場合
- 同じロジックで複数の入力パターンをテスト
- バリデーション関数（正常値、境界値、異常値）
- データ変換関数（様々な入力形式）
- エラーハンドリング（複数のエラーケース）

#### 実装例

```go
func TestValidateCountryName_TableDriven(t *testing.T) {
    tests := []struct {
        name        string
        input       string
        expectError bool
        errorMsg    string
    }{
        {
            name:        "正常系: 有効な国名",
            input:       "Japan",
            expectError: false,
        },
        {
            name:        "エラー系: 空文字",
            input:       "",
            expectError: true,
            errorMsg:    "country name is required",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateCountryName(tt.input)
            // アサーション
        })
    }
}
```

### 2. 個別テスト

#### 適している場合
- テストケースごとに異なるセットアップが必要
- モックの設定が複雑
- テストケースが少ない（1-2個）
- テストケースごとに異なるアサーション

#### 実装例

```go
func TestCountryImpl_Select_ComplexSetup(t *testing.T) {
    // 複雑なモック設定
    testCountry := entity.Country{
        CountryID:  1,
        Country:    types.CountryName("Japan"),
        LastUpdate: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
    }
    mockRepo := mocks.NewMockCountryRepository().
        WithCountries([]entity.Country{testCountry})

    // 詳細なアサーション
    countries, err := usecase.Select()
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    // 複数の条件を同時に検証
}
```

## レイヤー別テスト戦略

### 1. UseCase層のテスト

#### テスト対象
- ビジネスロジック
- エラーハンドリング
- データ変換

#### モック対象
- Repository層

#### テスト例

```go
func TestCountryImpl_Select(t *testing.T) {
    // モックリポジトリの準備
    mockRepo := mocks.NewMockCountryRepository().
        WithCountries(mocks.CreateTestCountries())

    usecase := &CountryImpl{repo: mockRepo}

    // テスト実行
    countries, err := usecase.Select()

    // アサーション
    if err != nil {
        t.Errorf("expected no error, got %v", err)
    }
    if len(countries) != 3 {
        t.Errorf("expected 3 countries, got %d", len(countries))
    }
}
```

### 2. ドメイン層のテスト

#### テスト対象
- バリデーション関数
- エンティティのメソッド
- 型エイリアスの動作

#### テスト例

```go
func TestValidateCountryName(t *testing.T) {
    // 境界値テスト
    err := ValidateCountryName("")
    if err == nil {
        t.Error("expected error for empty string")
    }
    
    // 正常系テスト
    err = ValidateCountryName("Japan")
    if err != nil {
        t.Errorf("expected no error, got %v", err)
    }
}
```

### 3. Controller層のテスト

#### テスト対象
- HTTPリクエスト/レスポンス
- バインディング処理
- エラーレスポンスの形式

#### モック対象
- UseCase層

## エラーハンドリングのテスト

### 1. エラー伝播のテスト

```go
func TestCountryImpl_Select_Error(t *testing.T) {
    mockRepo := mocks.NewMockCountryRepository().
        WithSelectError(errors.New("database connection error"))

    usecase := &CountryImpl{repo: mockRepo}

    countries, err := usecase.Select()

    // エラーの存在確認
    if err == nil {
        t.Error("expected error, got nil")
    }

    // エラーメッセージの確認
    expectedError := "failed to select countries: database connection error"
    if err.Error() != expectedError {
        t.Errorf("expected error message '%s', got '%s'", expectedError, err.Error())
    }
}
```

### 2. エラーレスポンスのテスト

```go
func TestCountryController_Create_ValidationError(t *testing.T) {
    // バリデーションエラーのテスト
    request := map[string]interface{}{
        "country": "", // 無効な値
    }

    // エラーレスポンスの形式確認
    // - HTTPステータスコード: 400
    // - エラーメッセージの内容
    // - レスポンス構造体の形式
}
```

## テスト実行

### コマンド

```bash
# 特定のパッケージのテスト
go test ./main/usecase/country/ -v

# カバレッジ付きテスト
go test ./main/usecase/country/ -v -cover

# 全パッケージのテスト
go test ./main/... -v

# ベンチマークテスト
go test ./main/usecase/country/ -bench=.

# テストライブラリのテスト
go test ./main/testlib/ -v
```

### テストカバレッジ

```bash
# カバレッジレポートの生成
go test ./main/... -coverprofile=coverage.out

# HTMLレポートの生成
go tool cover -html=coverage.out -o coverage.html
```

## ベストプラクティス

### 1. テストの命名規則

- **関数名**: `Test[構造体名]_[メソッド名]_[シナリオ]`
- **テーブル駆動**: `Test[関数名]_TableDriven`
- **エラー系**: `Test[関数名]_Error`

### 2. アサーションの原則

- **明確性**: 期待値と実際の値を明確に比較
- **詳細性**: エラーメッセージに十分な情報を含める
- **一貫性**: 同じパターンでアサーションを記述

### 3. モックの設計原則

- **単一責任**: 1つのモックは1つの責任を持つ
- **設定の柔軟性**: 様々なテストケースに対応可能
- **状態の管理**: テスト間で状態を適切に管理

### 4. テストデータの管理

- **再利用性**: 共通のテストデータを作成
- **一貫性**: テストデータの形式を統一
- **保守性**: テストデータの更新を容易にする

## 今後の拡張予定

### 1. 統合テスト
- 実際のデータベースを使用したテスト
- エンドツーエンドテスト
- APIテスト

### 2. パフォーマンステスト
- ベンチマークテストの追加
- 負荷テストの実装

### 3. テスト自動化
- CI/CDパイプラインでの自動テスト実行
- テストレポートの自動生成
- カバレッジ閾値の設定
