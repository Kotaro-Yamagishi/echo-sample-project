package models

// Entity の定義
// どの層からも呼び出し可能

type User struct {
	ID   int    `json:"id" gorm:"primary_key"`
	Name string `json:"name"`
}
