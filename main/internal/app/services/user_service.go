package services

import (
    "echoProject/main/internal/models"
    "echoProject/main/internal/app/repositories"
)

type UserService struct {
    UserRepository repositories.UserRepositoryIF
}

func NewUserService(repo repositories.UserRepositoryIF) UserServiceIF {
    return &UserService{UserRepository: repo}
}

func (interactor *UserService) Add(u domain.User)  {
    interactor.UserRepository.Store(u)
}

func (interactor *UserService) GetInfo() []domain.User {
    return interactor.UserRepository.Select()
}

func (interactor *UserService) Delete(id string) {
    interactor.UserRepository.Delete(id)
}
