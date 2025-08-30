package mocks

import (
	"echoProject/domain/entity"
	"time"
)

// MockCityRepository CityRepositoryのモック実装
type MockCityRepository struct {
	cities    []entity.City
	insertErr error
	selectErr error
}

// NewMockCityRepository 新しいモックリポジトリを作成
func NewMockCityRepository() *MockCityRepository {
	return &MockCityRepository{
		cities: []entity.City{},
	}
}

// WithCities テストデータを設定
func (m *MockCityRepository) WithCities(cities []entity.City) *MockCityRepository {
	m.cities = cities
	return m
}

// WithInsertError 挿入エラーを設定
func (m *MockCityRepository) WithInsertError(err error) *MockCityRepository {
	m.insertErr = err
	return m
}

// WithSelectError 取得エラーを設定
func (m *MockCityRepository) WithSelectError(err error) *MockCityRepository {
	m.selectErr = err
	return m
}

// Select 都市一覧を取得
func (m *MockCityRepository) Select() ([]entity.City, error) {
	if m.selectErr != nil {
		return nil, m.selectErr
	}
	return m.cities, nil
}

// Insert 都市を挿入
func (m *MockCityRepository) Insert(city entity.City) error {
	if m.insertErr != nil {
		return m.insertErr
	}
	m.cities = append(m.cities, city)
	return nil
}

// GetCities 内部のcitiesを取得（テスト用）
func (m *MockCityRepository) GetCities() []entity.City {
	return m.cities
}

// Clear 内部データをクリア
func (m *MockCityRepository) Clear() {
	m.cities = []entity.City{}
}

// CreateTestCity テスト用のCityエンティティを作成
func CreateTestCity(id uint16, name string, countryID uint16) entity.City {
	return entity.City{
		CityID:     id,
		City:       name,
		CountryID:  countryID,
		LastUpdate: time.Now(),
	}
}

// CreateTestCities テスト用のCityエンティティのスライスを作成
func CreateTestCities() []entity.City {
	return []entity.City{
		CreateTestCity(1, "Tokyo", 1),
		CreateTestCity(2, "Osaka", 1),
		CreateTestCity(3, "New York", 2),
		CreateTestCity(4, "London", 3),
	}
}
