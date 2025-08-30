package testlib

import (
	"errors"
	"testing"

	"echoProject/domain/entity"
	"echoProject/domain/types"
	"echoProject/testlib/mocks"
)

// TestMockUsage モックの使用例テスト
func TestMockUsage(t *testing.T) {
	// 1. 基本的なモック作成
	mockRepo := mocks.NewMockCountryRepository()

	// 2. テストデータの設定
	testCountries := mocks.CreateTestCountries()
	mockRepo.WithCountries(testCountries)

	// 3. 正常系テスト
	countries, err := mockRepo.Select()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if len(countries) != 3 {
		t.Errorf("expected 3 countries, got %d", len(countries))
	}

	// 4. エラー系テスト
	mockRepo.WithSelectError(errors.New("database error"))
	countries, err = mockRepo.Select()
	if err == nil {
		t.Error("expected error, got nil")
	}

	// 5. 挿入テスト
	mockRepo.Clear() // データをクリア
	country := entity.Country{
		Country:    types.CountryName("NewCountry"),
		LastUpdate: CreateTestTime(),
	}

	err = mockRepo.Insert(country)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	insertedCountries := mockRepo.GetCountries()
	if len(insertedCountries) != 1 {
		t.Errorf("expected 1 country, got %d", len(insertedCountries))
	}
}

// TestMockBuilderPattern ビルダーパターンの使用例
func TestMockBuilderPattern(t *testing.T) {
	// ビルダーパターンでモックを作成
	mockRepo := mocks.NewMockCountryRepository().
		WithCountries(mocks.CreateTestCountries()).
		WithInsertError(errors.New("insert failed"))

	// 設定が正しく反映されていることを確認
	countries, err := mockRepo.Select()
	if err != nil {
		t.Errorf("expected no error for select, got %v", err)
	}

	if len(countries) != 3 {
		t.Errorf("expected 3 countries, got %d", len(countries))
	}

	// 挿入エラーが設定されていることを確認
	testCountry := entity.Country{
		Country:    types.CountryName("TestCountry"),
		LastUpdate: CreateTestTime(),
	}

	err = mockRepo.Insert(testCountry)
	if err == nil {
		t.Error("expected insert error, got nil")
	}

	expectedError := "insert failed"
	if err.Error() != expectedError {
		t.Errorf("expected error message '%s', got '%s'", expectedError, err.Error())
	}
}
