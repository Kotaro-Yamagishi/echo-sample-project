package repository

import "echoProject/main/domain/entity"

type Country interface {
	Select() ([]entity.Country, error)
	Insert(country entity.Country) error
}
