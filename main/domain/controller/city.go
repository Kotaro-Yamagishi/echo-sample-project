package ctrIF

import (
	"echoProject/main/domain/entity"
)

type City interface {
	GetCity() []entity.City
}
