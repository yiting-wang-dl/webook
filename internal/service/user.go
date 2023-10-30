package service

import (
	"context"
	"errors"
	"github.com/webook/internal/domain"
	"github.com/webook/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail        = repository.ErrDuplicateEmail
	ErrInvalidUserOrPassword = errors.New("User doesn't exist or password is wrong") // why use erroors.New() instead of repository.ErrUerNotFound
)

type UserService struct {
	repo *repository.CacheUserRepository
}

func NewUserService(repo *repository.CacheUserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (svc *UserService) Signup(ctx context.Context, u domain.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return svc.repo.Create(ctx, u)
}

func (svc *UserService) Login(ctx context.Context, email string, password string) (domain.User, error) { //why return domain.User?
	user, err := svc.repo.FindByEmail(ctx, email)
	if err == repository.ErrUserNotfound { // what does repository.ErrUserNotFound return?
		return domain.User{}, ErrInvalidUserOrPassword
	}
	if err != nil {
		return domain.User{}, err
	}
	// check password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	return user, nil
}

func (svc *UserService) FindById(ctx context.Context, uid int64) (domain.User, error) {
	return svc.repo.FindById(ctx, uid) // why pass uid here and return domain.User?
}

func (svc *UserService) UpdateNonSensitiveInfo(ctx context.Context, user domain.User) error {
	return svc.repo.UpdateNonZeroFields(ctx, user)
}
