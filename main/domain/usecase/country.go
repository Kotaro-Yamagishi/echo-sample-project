package usecase

import "echoProject/domain/entity"

type Country interface {
	Select() ([]entity.Country, error)
	Insert(country entity.Country) error
}
