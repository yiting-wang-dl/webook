package service

import (
	"context"
	"errors"
	"github.com/webook/internal/domain"
	"github.com/webook/internal/repository"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail        = repository.ErrDuplicateUser
	ErrInvalidUserOrPassword = errors.New("User doesn't exist or password is not correct")
)

type UserService interface {
	Signup(ctx context.Context, u domain.User) error
	Login(ctx context.Context, email string, password string) (domain.User, error)
	UpdateNonSensitiveInfo(ctx context.Context, user domain.User) error
	FindById(ctx context.Context, uid int64) (domain.User, error)
	FindOrCreate(ctx context.Context, phone string) (domain.User, error)
	FindOrCreateByWechat(ctx context.Context, info domain.WechatInfo) (domain.User, error)
}

type userService struct {
	repo repository.UserRepository
	//logger *zap.Logger
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
		//logger: zap.L(),
	}
}

func (svc *userService) Signup(ctx context.Context, u domain.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return svc.repo.Create(ctx, u)
}

func (svc *userService) Login(ctx context.Context, email string, password string) (domain.User, error) {
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

func (svc *userService) FindById(ctx context.Context, uid int64) (domain.User, error) {
	return svc.repo.FindById(ctx, uid) // why pass uid here and return domain.User?
}

func (svc *userService) UpdateNonSensitiveInfo(ctx context.Context, user domain.User) error {
	return svc.repo.UpdateNonZeroFields(ctx, user)
}

func (svc *userService) FindOrCreate(ctx context.Context, phone string) (domain.User, error) {
	// Find the user by phone first, most of the users are existing users
	u, err := svc.repo.FindByPhone(ctx, phone)
	if err != repository.ErrUserNotfound {
		// two possibilities
		// err == nil, u is good
		// err != nil, system error
		return u, err
	}
	// Didn't find user
	err = svc.repo.Create(ctx, domain.User{
		Phone: phone,
	})
	// two possibilities,
	// one is that err is the conflict of primary key (phone)
	// the other one is err != nil system error
	if err != nil && err != repository.ErrDuplicateUser {
		return domain.User{}, err
	}
	return svc.repo.FindByPhone(ctx, phone)
}

func (svc *userService) FindOrCreateByWechat(ctx context.Context, wechatInfo domain.WechatInfo) (domain.User, error) {
	u, err := svc.repo.FindByWechat(ctx, wechatInfo.OpenId)
	if err != repository.ErrUserNotfound {
		return u, err
	}
	// new user
	// wechatInfo in JSON
	zap.L().Info("New User", zap.Any("wechatInfo", wechatInfo))
	//svc.logger.Info("New User", zap.Any("wechatInfo", wechatInfo))
	err = svc.repo.Create(ctx, domain.User{
		WechatInfo: wechatInfo,
	})
	if err != nil && err != repository.ErrDuplicateUser {
		return domain.User{}, err
	}
	return svc.repo.FindByWechat(ctx, wechatInfo.OpenId)
}
