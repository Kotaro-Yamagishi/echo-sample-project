package repositories

import "echoProject/main/internal/models"

type UserRepositoryIF interface {
    Store(u domain.User) 
    Select() []domain.User
    Delete(id string)
}
