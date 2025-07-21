package repositories

import "echoProject/main/internal/models"

type UserRepositoryIF interface {
    Store(u models.User) 
    Select() []models.User
    Delete(id string)
}
