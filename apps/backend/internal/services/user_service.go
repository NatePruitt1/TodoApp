package services

import (
	"backend/internal/config"
	"backend/internal/dto"
	"backend/internal/models"
	"backend/internal/repository"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var ErrUsernameTaken = errors.New("Username taken.")

// dummyPasswordHash is a valid bcrypt hash (of an arbitrary password) used to
// perform a dummy comparison when a username lookup fails. This ensures the
// "user not found" path takes roughly the same amount of time as the
// "wrong password" path, mitigating username-enumeration via timing attacks.
const dummyPasswordHash = "$2a$10$iFc4Bhptv3kbZ7jFJq0EV.pCjZRva7nw1FD6ZyI55sNwWiMyRm9re"

type UserService interface {
	CreateAccount(ctx context.Context, username, password string) (dto.LoginResponseDTO, string, error)
	AuthenticateAccount(ctx context.Context, username, password string) (dto.LoginResponseDTO, string, error)
	RefreshAccount(ctx context.Context, userId uuid.UUID) (dto.LoginResponseDTO, string, error)
}

type UserServiceImpl struct {
	repo   repository.UserRepository
	config *config.Config
}

func NewUserService(repository *repository.UserRepositoryImpl, cfg *config.Config) UserService {
	return &UserServiceImpl{
		repo:   repository,
		config: cfg,
	}
}

func (u *UserServiceImpl) AuthenticateAccount(ctx context.Context, username, password string) (dto.LoginResponseDTO, string, error) {
	uret, err := u.repo.GetByUsername(username)
	if err != nil {
		bcrypt.CompareHashAndPassword([]byte(dummyPasswordHash), []byte(password))
		return dto.LoginResponseDTO{}, "", err
	}

	if e := bcrypt.CompareHashAndPassword([]byte(uret.PasswordHash), []byte(password)); e != nil {
		return dto.LoginResponseDTO{}, "", e
	}

	token, e := GenerateToken(uret, u.config.JwtSecret)
	if e != nil {
		return dto.LoginResponseDTO{}, "", e
	}

	curr := time.Now()
	u.repo.UpdateLastLogin(uret.ID, curr)

	return dto.LoginResponseDTO{
		ID:        uret.ID,
		Username:  username,
		CreatedAt: uret.CreatedAt,
		LastLogin: &curr,
	}, token, nil
}

func (u *UserServiceImpl) RefreshAccount(ctx context.Context, userId uuid.UUID) (dto.LoginResponseDTO, string, error) {
	uret, err := u.repo.GetByID(userId)
	if err != nil {
		return dto.LoginResponseDTO{}, "", err
	}

	token, e := GenerateToken(uret, u.config.JwtSecret)
	if e != nil {
		return dto.LoginResponseDTO{}, "", e
	}

	curr := time.Now()
	u.repo.UpdateLastLogin(uret.ID, curr)

	return dto.LoginResponseDTO{
		ID:        uret.ID,
		Username:  uret.Username,
		CreatedAt: uret.CreatedAt,
		LastLogin: &curr,
	}, token, nil
}

func (u *UserServiceImpl) CreateAccount(ctx context.Context, username, password string) (dto.LoginResponseDTO, string, error) {
	// Get if there is a user with the name.
	_, e := u.repo.GetByUsername(username)
	if e == nil {
		return dto.LoginResponseDTO{}, "", ErrUsernameTaken
	}

	bytes, e := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if e != nil {
		return dto.LoginResponseDTO{}, "", e
	}

	passwordHash := string(bytes)
	user := models.User{
		ID:           uuid.New(),
		Username:     username,
		PasswordHash: passwordHash,
		CreatedAt:    time.Now(),
		LastLogin:    nil,
	}

	e = u.repo.Create(user)
	if e != nil {
		return dto.LoginResponseDTO{}, "", e
	}

	token, e := GenerateToken(user, u.config.JwtSecret)
	if e != nil {
		return dto.LoginResponseDTO{}, "", e
	}

	return dto.LoginResponseDTO{
		ID:        user.ID,
		Username:  username,
		CreatedAt: user.CreatedAt,
		LastLogin: user.LastLogin,
	}, token, nil
}
