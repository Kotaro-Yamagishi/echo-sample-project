package usecase

import "echoProject/main/src/domain"

// repository の構成

type UserRepository interface {
	Store(domain.User)
	Select() []domain.User
	Delete(id string)
}
