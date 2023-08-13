package service

import (
	"boot-camp-homework/homework/week2/webhook/internal/domain"
	"boot-camp-homework/homework/week2/webhook/internal/repository"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserDuplicateEmail = repository.ErrUserDuplicateEmail
	ErrUserNotFound       = repository.ErrUserNotFound
	ErrUserInfoNotFound   = repository.ErrUserInfoNotFound
)
var ErrInvalidUserOrPassword = errors.New("账号/邮箱或密码不对")

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (svc *UserService) SignUp(ctx context.Context, u domain.User) error {
	// 你要考虑加密放在哪里
	bcryptPwd, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(bcryptPwd)
	// 然后就是,存起来
	return svc.repo.Create(ctx, u)
}

func (svc *UserService) Login(ctx *gin.Context, u domain.User) (domain.User, error) {
	//先找用户
	resuser := domain.User{}
	repouser, err := svc.repo.FindByEmail(ctx, u)
	if errors.Is(err, repository.ErrUserNotFound) {
		return resuser, ErrInvalidUserOrPassword
	}
	if err != nil {
		return resuser, err
	}
	//比较密码
	err = bcrypt.CompareHashAndPassword([]byte(repouser.Password), []byte(u.Password))
	if err != nil {
		// DEBUG
		return resuser, ErrInvalidUserOrPassword
	}
	resuser = repouser
	return resuser, nil
}

func (svc *UserService) Edit(ctx context.Context, uinfo domain.UserInfo) error {
	//先找用户
	user := domain.User{
		Id: uinfo.Id,
	}
	finduser, err := svc.repo.FindUserTableById(ctx, user.Id)
	if errors.Is(err, repository.ErrUserNotFound) {
		return ErrUserNotFound
	}
	if err != nil {
		return err
	}

	if finduser.Id == 0 || finduser.Id != uinfo.Id {
		return ErrUserNotFound
	}
	//找到userid,存入info.id中，并将info中其他信息存入数据库,
	_, err = svc.repo.FindUserInfoTableById(ctx, user.Id)
	if errors.Is(err, ErrUserInfoNotFound) {
		err = svc.repo.CreateUserInfo(ctx, uinfo)
		return err
	}
	err = svc.repo.UpdateUserInfo(ctx, uinfo)
	return err
}

func (svc *UserService) Profile(ctx context.Context, id int64) (domain.UserInfo, error) {
	//先找用户
	user := domain.User{
		Id: id,
	}
	finduser, err := svc.repo.FindUserTableById(ctx, user.Id)
	if errors.Is(err, repository.ErrUserNotFound) {
		return domain.UserInfo{}, ErrUserNotFound
	}
	if err != nil {
		return domain.UserInfo{}, err
	}

	if finduser.Id == 0 || finduser.Id != id {
		return domain.UserInfo{}, ErrUserNotFound
	}
	//找到userid,存入info.id中，并将info中其他信息存入数据库,
	uinfo, err := svc.repo.FindUserInfoTableById(ctx, user.Id)
	if err != nil {
		return domain.UserInfo{}, err
	}
	return uinfo, nil
}
