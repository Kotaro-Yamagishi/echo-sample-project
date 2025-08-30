package country

import (
	"errors"
	"testing"
	"time"

	"echoProject/domain/entity"
	"echoProject/domain/types"
	"echoProject/testlib/mocks"
)

// 既存のモック定義を削除（testlib/mocksを使用）

// TestCountryImpl_Select 正常系テスト
func TestCountryImpl_Select(t *testing.T) {
	// テストデータの準備
	mockRepo := mocks.NewMockCountryRepository().
		WithCountries(mocks.CreateTestCountries())

	// UseCaseの作成
	usecase := &CountryImpl{
		repo: mockRepo,
	}

	// テスト実行
	countries, err := usecase.Select()

	// アサーション
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if len(countries) != 3 {
		t.Errorf("expected 3 countries, got %d", len(countries))
	}

	if countries[0].Country != types.CountryName("Japan") {
		t.Errorf("expected first country to be Japan, got %s", countries[0].Country)
	}
}

// TestCountryImpl_Select_Error エラー系テスト
func TestCountryImpl_Select_Error(t *testing.T) {
	// エラーを返すモックリポジトリ
	mockRepo := mocks.NewMockCountryRepository().
		WithSelectError(errors.New("database connection error"))

	usecase := &CountryImpl{
		repo: mockRepo,
	}

	// テスト実行
	countries, err := usecase.Select()

	// アサーション
	if err == nil {
		t.Error("expected error, got nil")
	}

	if countries != nil {
		t.Errorf("expected nil countries, got %v", countries)
	}

	// エラーメッセージの確認
	expectedError := "failed to select countries: database connection error"
	if err.Error() != expectedError {
		t.Errorf("expected error message '%s', got '%s'", expectedError, err.Error())
	}
}

// TestCountryImpl_Insert 正常系テスト
func TestCountryImpl_Insert(t *testing.T) {
	mockRepo := mocks.NewMockCountryRepository()

	usecase := &CountryImpl{
		repo: mockRepo,
	}

	// テストデータ
	country := entity.Country{
		Country:    types.CountryName("Japan"),
		LastUpdate: time.Now(),
	}

	// テスト実行
	err := usecase.Insert(country)

	// アサーション
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	countries := mockRepo.GetCountries()
	if len(countries) != 1 {
		t.Errorf("expected 1 country in repository, got %d", len(countries))
	}

	if countries[0].Country != types.CountryName("Japan") {
		t.Errorf("expected country to be Japan, got %s", countries[0].Country)
	}
}

// TestCountryImpl_Insert_Error エラー系テスト
func TestCountryImpl_Insert_Error(t *testing.T) {
	mockRepo := mocks.NewMockCountryRepository().
		WithInsertError(errors.New("database insert error"))

	usecase := &CountryImpl{
		repo: mockRepo,
	}

	country := entity.Country{
		Country:    types.CountryName("Japan"),
		LastUpdate: time.Now(),
	}

	// テスト実行
	err := usecase.Insert(country)

	// アサーション
	if err == nil {
		t.Error("expected error, got nil")
	}

	expectedError := "failed to insert country: database insert error"
	if err.Error() != expectedError {
		t.Errorf("expected error message '%s', got '%s'", expectedError, err.Error())
	}
}

// TestCountryImpl_Select_Empty 空の結果テスト
func TestCountryImpl_Select_Empty(t *testing.T) {
	mockRepo := mocks.NewMockCountryRepository()

	usecase := &CountryImpl{
		repo: mockRepo,
	}

	// テスト実行
	countries, err := usecase.Select()

	// アサーション
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if len(countries) != 0 {
		t.Errorf("expected 0 countries, got %d", len(countries))
	}
}

// TestCountryImpl_Select_TableDriven テーブル駆動テスト
func TestCountryImpl_Select_TableDriven(t *testing.T) {
	tests := []struct {
		name          string
		mockCountries []entity.Country
		mockError     error
		expectedCount int
		expectedError bool
		expectedFirst string
	}{
		{
			name: "正常系: 複数の国を取得",
			mockCountries: []entity.Country{
				{CountryID: 1, Country: types.CountryName("Japan"), LastUpdate: time.Now()},
				{CountryID: 2, Country: types.CountryName("USA"), LastUpdate: time.Now()},
			},
			mockError:     nil,
			expectedCount: 2,
			expectedError: false,
			expectedFirst: "Japan",
		},
		{
			name:          "正常系: 空の結果",
			mockCountries: []entity.Country{},
			mockError:     nil,
			expectedCount: 0,
			expectedError: false,
			expectedFirst: "",
		},
		{
			name:          "エラー系: データベースエラー",
			mockCountries: nil,
			mockError:     errors.New("database error"),
			expectedCount: 0,
			expectedError: true,
			expectedFirst: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// モックリポジトリの準備
			mockRepo := mocks.NewMockCountryRepository().
				WithCountries(tt.mockCountries).
				WithSelectError(tt.mockError)

			usecase := &CountryImpl{
				repo: mockRepo,
			}

			// テスト実行
			countries, err := usecase.Select()

			// アサーション
			if tt.expectedError {
				if err == nil {
					t.Error("expected error, got nil")
				}
				if countries != nil {
					t.Errorf("expected nil countries, got %v", countries)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
				if len(countries) != tt.expectedCount {
					t.Errorf("expected %d countries, got %d", tt.expectedCount, len(countries))
				}
				if tt.expectedFirst != "" && len(countries) > 0 {
					if string(countries[0].Country) != tt.expectedFirst {
						t.Errorf("expected first country to be %s, got %s", tt.expectedFirst, countries[0].Country)
					}
				}
			}
		})
	}
}

// TestCountryImpl_Insert_TableDriven テーブル駆動テスト
func TestCountryImpl_Insert_TableDriven(t *testing.T) {
	tests := []struct {
		name          string
		country       entity.Country
		mockError     error
		expectedError bool
		expectedMsg   string
	}{
		{
			name: "正常系: 国を挿入",
			country: entity.Country{
				Country:    types.CountryName("Japan"),
				LastUpdate: time.Now(),
			},
			mockError:     nil,
			expectedError: false,
		},
		{
			name: "エラー系: データベース挿入エラー",
			country: entity.Country{
				Country:    types.CountryName("Japan"),
				LastUpdate: time.Now(),
			},
			mockError:     errors.New("insert error"),
			expectedError: true,
			expectedMsg:   "failed to insert country: insert error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// モックリポジトリの準備
			mockRepo := mocks.NewMockCountryRepository().
				WithInsertError(tt.mockError)

			usecase := &CountryImpl{
				repo: mockRepo,
			}

			// テスト実行
			err := usecase.Insert(tt.country)

			// アサーション
			if tt.expectedError {
				if err == nil {
					t.Error("expected error, got nil")
				}
				if tt.expectedMsg != "" && err.Error() != tt.expectedMsg {
					t.Errorf("expected error message '%s', got '%s'", tt.expectedMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
				countries := mockRepo.GetCountries()
				if len(countries) != 1 {
					t.Errorf("expected 1 country in repository, got %d", len(countries))
				}
			}
		})
	}
}

// TestCountryImpl_Select_ComplexSetup 複雑なセットアップが必要なテスト
func TestCountryImpl_Select_ComplexSetup(t *testing.T) {
	// 複雑なモック設定
	testCountry := entity.Country{
		CountryID:  1,
		Country:    types.CountryName("Japan"),
		LastUpdate: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	mockRepo := mocks.NewMockCountryRepository().
		WithCountries([]entity.Country{testCountry})

	// 特別な設定が必要なUseCase
	usecase := &CountryImpl{
		repo: mockRepo,
	}

	// 複雑なアサーション
	countries, err := usecase.Select()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 詳細な検証
	if len(countries) != 1 {
		t.Fatalf("expected 1 country, got %d", len(countries))
	}

	country := countries[0]
	if country.CountryID != 1 {
		t.Errorf("expected CountryID 1, got %d", country.CountryID)
	}

	if country.Country != types.CountryName("Japan") {
		t.Errorf("expected Country Japan, got %s", country.Country)
	}

	// 日時の詳細検証
	expectedTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	if !country.LastUpdate.Equal(expectedTime) {
		t.Errorf("expected LastUpdate %v, got %v", expectedTime, country.LastUpdate)
	}
}

// TestCountryImpl_Insert_WithValidation バリデーション付きの挿入テスト
func TestCountryImpl_Insert_WithValidation(t *testing.T) {
	mockRepo := mocks.NewMockCountryRepository()

	usecase := &CountryImpl{
		repo: mockRepo,
	}

	// バリデーション済みのCountryを作成
	country, err := entity.NewCountry("ValidCountry")
	if err != nil {
		t.Fatalf("failed to create valid country: %v", err)
	}

	// 挿入テスト
	err = usecase.Insert(country)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// 挿入されたデータの検証
	countries := mockRepo.GetCountries()
	if len(countries) != 1 {
		t.Errorf("expected 1 country in repository, got %d", len(countries))
	}

	insertedCountry := countries[0]
	if insertedCountry.Country != country.Country {
		t.Errorf("expected country %s, got %s", country.Country, insertedCountry.Country)
	}

	// バリデーションが正しく動作することを確認
	if err := entity.Validate(string(insertedCountry.Country)); err != nil {
		t.Errorf("inserted country should be valid: %v", err)
	}
}
