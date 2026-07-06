package services

import (
	"backend/internal/dto"
	"backend/internal/repository"
	"context"
)

type UserService interface {
	Authenticate(ctx context.Context, username, password string)
}

type UserServiceImpl struct {
	repo *repository.UserRepositoryImpl
}

func NewUserService(repository *repository.UserRepositoryImpl) *UserServiceImpl {
	return &UserServiceImpl{
		repo: repository,
	}
}

func (u *UserServiceImpl) Authenticate(ctx context.Context, username, password string) (dto.LoginResponseDTO, error) {

}
