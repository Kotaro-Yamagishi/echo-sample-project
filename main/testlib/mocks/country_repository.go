package mocks

import (
	"echoProject/main/domain/entity"
	"echoProject/main/domain/types"
	"time"
)

// MockCountryRepository CountryRepositoryのモック実装
type MockCountryRepository struct {
	countries []entity.Country
	insertErr error
	selectErr error
}

// NewMockCountryRepository 新しいモックリポジトリを作成
func NewMockCountryRepository() *MockCountryRepository {
	return &MockCountryRepository{
		countries: []entity.Country{},
	}
}

// WithCountries テストデータを設定
func (m *MockCountryRepository) WithCountries(countries []entity.Country) *MockCountryRepository {
	m.countries = countries
	return m
}

// WithInsertError 挿入エラーを設定
func (m *MockCountryRepository) WithInsertError(err error) *MockCountryRepository {
	m.insertErr = err
	return m
}

// WithSelectError 取得エラーを設定
func (m *MockCountryRepository) WithSelectError(err error) *MockCountryRepository {
	m.selectErr = err
	return m
}

// Select 国一覧を取得
func (m *MockCountryRepository) Select() ([]entity.Country, error) {
	if m.selectErr != nil {
		return nil, m.selectErr
	}
	return m.countries, nil
}

// Insert 国を挿入
func (m *MockCountryRepository) Insert(country entity.Country) error {
	if m.insertErr != nil {
		return m.insertErr
	}
	m.countries = append(m.countries, country)
	return nil
}

// GetCountries 内部のcountriesを取得（テスト用）
func (m *MockCountryRepository) GetCountries() []entity.Country {
	return m.countries
}

// Clear 内部データをクリア
func (m *MockCountryRepository) Clear() {
	m.countries = []entity.Country{}
}

// CreateTestCountry テスト用のCountryエンティティを作成
func CreateTestCountry(id uint16, name string) entity.Country {
	return entity.Country{
		CountryID:  id,
		Country:    types.CountryName(name),
		LastUpdate: time.Now(),
	}
}

// CreateTestCountries テスト用のCountryエンティティのスライスを作成
func CreateTestCountries() []entity.Country {
	return []entity.Country{
		CreateTestCountry(1, "Japan"),
		CreateTestCountry(2, "USA"),
		CreateTestCountry(3, "UK"),
	}
}
