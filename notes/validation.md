# バリデーション設計ドキュメント

## 概要

このプロジェクトでは、クリーンアーキテクチャに基づいた階層別バリデーション設計を採用しています。型エイリアスを活用した型安全性と、各層での適切なバリデーション責任の分離を実現しています。

## バリデーション設計の基本方針

### 1. 階層別バリデーション

```
Controller層 → UseCase層 → Entity層
     ↓           ↓          ↓
  入力値検証   ビジネス検証   ドメイン検証
```

### 2. 型エイリアスによる型安全性

基本のプリミティブ型に対して型エイリアスを定義し、バリデーション機能を組み込んだ型システムを構築しています。

## 実装詳細

### 1. 型エイリアスの定義

#### 基本型エイリアス
```go
// main/domain/types/country.go
type CountryName string

func ValidateCountryName(name string) error {
    if name == "" {
        return errors.New("country name is required")
    }
    if len(name) > 50 {
        return errors.New("country name cannot exceed 50 characters")
    }
    return nil
}

func (c CountryName) Validate() error {
    return ValidateCountryName(string(c))
}

func (c CountryName) String() string {
    return string(c)
}
```

#### 基本型エイリアス
```go
// main/domain/types/city.go
type CityName string

func ValidateCityName(name string) error {
    if name == "" {
        return errors.New("city name is required")
    }
    return nil
}

func (c CityName) Validate() error {
    return ValidateCityName(string(c))
}

func (c CityName) String() string {
    return string(c)
}
```

### 2. Entity層でのバリデーション

#### Entity構造体
```go
// main/domain/entity/country.go
type Country struct {
    CountryID  uint16
    Country    types.CountryName  // 型エイリアスを使用
    LastUpdate time.Time
}

// バリデーション関数
func Validate(countryName string) error {
    if err := types.ValidateCountryName(countryName); err != nil {
        return err
    }
    return nil
}

// ファクトリ関数（バリデーションなし）
func NewCountry(countryName string) Country {
    return Country{
        Country:    types.CountryName(countryName),
        LastUpdate: time.Now(),
    }
}

// バリデーション付きファクトリ関数
func NewValidatedCountry(countryName string) (Country, error) {
    if err := Validate(countryName); err != nil {
        return Country{}, err
    }
    return NewCountry(countryName), nil
}
```

### 3. UseCase層でのバリデーション

```go
// main/usecase/country/country.go
func (s *CountryImpl) CreateCountry(name string) error {
    // バリデーション付きでEntityを作成
    country, err := entity.NewValidatedCountry(name)
    if err != nil {
        return err
    }
    
    // ビジネスロジックのバリデーション
    if err := s.validateBusinessRules(&country); err != nil {
        return err
    }
    
    return s.repo.Insert(country)
}

func (s *CountryImpl) validateBusinessRules(country *entity.Country) error {
    // 重複チェック
    if exists, _ := s.repo.ExistsByName(string(country.Country)); exists {
        return errors.New("country with this name already exists")
    }
    return nil
}
```

### 4. Controller層でのバリデーション

```go
// main/app/controller/country/country.go
func (c *CountryImpl) Create(ctx echo.Context) error {
    var request struct {
        Country string `json:"country" validate:"required"`
    }
    
    // リクエストバインディングのバリデーション
    if err := ctx.Bind(&request); err != nil {
        return ctx.JSON(http.StatusBadRequest, output.NewErrorResponse(
            http.StatusBadRequest,
            "invalid request format",
        ))
    }
    
    // 構造体タグによるバリデーション
    if err := c.validator.Struct(request); err != nil {
        return ctx.JSON(http.StatusBadRequest, output.NewErrorResponse(
            http.StatusBadRequest,
            err.Error(),
        ))
    }
    
    // Entity層でのバリデーション
    country, err := entity.NewValidatedCountry(request.Country)
    if err != nil {
        return ctx.JSON(http.StatusBadRequest, output.NewErrorResponse(
            http.StatusBadRequest,
            err.Error(),
        ))
    }
    
    // UseCase層での処理
    if err := c.uc.Insert(country); err != nil {
        return ctx.JSON(http.StatusInternalServerError, output.NewErrorResponse(
            http.StatusInternalServerError,
            "failed to create country",
        ))
    }
    
    return ctx.JSON(http.StatusCreated, output.NewSuccessResponse(country))
}
```

## バリデーションの階層構造

### 1. Controller層（Presentation層）
- **目的**: HTTPリクエストの形式チェック
- **対象**: リクエストボディ、クエリパラメータ、ヘッダー
- **例**:
  - JSON形式の妥当性
  - 必須フィールドの存在チェック
  - データ型の検証
  - 文字列長の基本チェック

### 2. UseCase層（Application層）
- **目的**: ビジネスロジックに基づく検証
- **対象**: 複数エンティティ間の関係性、ビジネスルール
- **例**:
  - 外部キー制約の検証
  - 重複チェック
  - ビジネス条件の満足性
  - 権限チェック

### 3. Entity層（Domain層）
- **目的**: ドメイン固有の制約チェック
- **対象**: エンティティ自体の整合性
- **例**:
  - ドメイン固有の値制約
  - エンティティの状態整合性
  - 不変条件の検証

## 型エイリアスのメリット

### 1. 型安全性
```go
// コンパイル時に型チェックが行われる
func NewCountry(countryName string) Country {
    return Country{
        Country: types.CountryName(countryName),  // 型エイリアスを使用
        LastUpdate: time.Now(),
    }
}

// 型の不一致でコンパイルエラー
// country := entity.NewCountry(123)  // エラー
```

### 2. バリデーションの一元管理
```go
// 一度定義すれば、すべての場所で同じルールが適用される
func ValidateCountryName(name string) error {
    if name == "" {
        return errors.New("country name is required")
    }
    if len(name) > 50 {
        return errors.New("country name cannot exceed 50 characters")
    }
    return nil
}
```

### 3. 再利用性
```go
// 異なる場所で同じバリデーションを使用
func (c CountryName) Validate() error {
    return ValidateCountryName(string(c))
}

func Validate(countryName string) error {
    return types.ValidateCountryName(countryName)
}
```

## バリデーションの実行タイミング

### 1. 作成時バリデーション
```go
// バリデーション付きで作成
country, err := entity.NewValidatedCountry("Japan")
if err != nil {
    // バリデーションエラーを処理
    return err
}
```

### 2. 更新時バリデーション
```go
// 更新時も同様のバリデーションを実行
func (c *Country) SetCountryName(name string) error {
    if err := Validate(name); err != nil {
        return err
    }
    c.Country = types.CountryName(name)
    c.LastUpdate = time.Now()
    return nil
}
```

### 3. データベースからの取得時
```go
// データベースから取得したデータは既にバリデーション済み
func (r *countryRepository) GetByID(id uint16) (*entity.Country, error) {
    dbCountry, err := model.Countries(...).One(...)
    if err != nil {
        return nil, err
    }
    
    // バリデーションなしでEntityを作成
    country := entity.NewCountryFromDB(
        dbCountry.CountryID,
        dbCountry.Country,
        dbCountry.LastUpdate,
    )
    return &country, nil
}
```

## 設計のメリット

### 1. 責任の分離
- 各層が適切なバリデーション責任を持つ
- 関心の分離が明確
- テストしやすい構造

### 2. 型安全性
- コンパイル時の型チェック
- 実行時エラーの削減
- IDEのサポート

### 3. 保守性
- バリデーションルールの一元管理
- 変更時の影響範囲が明確
- 拡張しやすい構造

### 4. 再利用性
- バリデーション関数の再利用
- 型エイリアスの組み合わせ
- 共通化しやすい設計

## 今後の改善案

### 1. カスタムバリデーター
```go
// より高度なバリデーションルール
type CountryValidator interface {
    ValidateCreate(country *entity.Country) error
    ValidateUpdate(country *entity.Country) error
    ValidateName(name string) error
}
```

### 2. バリデーションエラーの構造化
```go
type ValidationError struct {
    Field   string
    Message string
    Code    string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation error on field %s: %s", e.Field, e.Message)
}
```

### 3. 国際化対応
```go
// 多言語対応のバリデーションメッセージ
func ValidateCountryName(name string, lang string) error {
    if name == "" {
        return getLocalizedError("country_name_required", lang)
    }
    return nil
}
```

## ベストプラクティス

### 1. バリデーションの順序
1. 入力値の基本チェック（Controller層）
2. ドメイン制約のチェック（Entity層）
3. ビジネスルールのチェック（UseCase層）

### 2. エラーメッセージの設計
- ユーザーフレンドリーなメッセージ
- 技術的な詳細は内部で保持
- 一貫したメッセージ形式

### 3. パフォーマンスの考慮
- 不要なバリデーションの回避
- 早期リターンの活用
- キャッシュの活用

## まとめ

このバリデーション設計により、以下の目標を達成しています：

1. **型安全性**: 型エイリアスによるコンパイル時チェック
2. **責任の分離**: 各層での適切なバリデーション責任
3. **再利用性**: バリデーション関数の共通化
4. **保守性**: 変更しやすく、拡張しやすい構造
5. **一貫性**: 統一されたバリデーションパターン

この設計は、クリーンアーキテクチャの原則に従い、型安全性とバリデーション機能を組み合わせた堅牢なシステムを実現しています。
