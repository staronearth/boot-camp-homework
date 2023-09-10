package service

import (
	"boot-camp-homework/webook/internal/domain"
	"boot-camp-homework/webook/internal/repository"
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserDuplicate    = repository.ErrUserDuplicate
	ErrUserNotFound     = repository.ErrUserNotFound
	ErrUserInfoNotFound = repository.ErrUserInfoNotFound
)
var ErrInvalidUserOrPassword = errors.New("账号/邮箱或密码不对")

type UserService interface {
	SignUp(ctx context.Context, u domain.User) error
	Login(ctx context.Context, u domain.User) (domain.User, error)
	Edit(ctx context.Context, uinfo domain.UserInfo) error
	Profile(ctx context.Context, id int64) (domain.UserInfo, error)
	FindOrCreate(ctx context.Context, phone string) (domain.User, error)
}
type CachedUserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &CachedUserService{
		repo: repo,
	}
}

func (svc *CachedUserService) SignUp(ctx context.Context, u domain.User) error {
	// 你要考虑加密放在哪里
	bcryptPwd, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(bcryptPwd)
	// 然后就是,存起来
	return svc.repo.Create(ctx, u)
}

func (svc *CachedUserService) Login(ctx context.Context, u domain.User) (domain.User, error) {
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

func (svc *CachedUserService) Edit(ctx context.Context, uinfo domain.UserInfo) error {
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

func (svc *CachedUserService) Profile(ctx context.Context, id int64) (domain.UserInfo, error) {
	//先在缓存中查找数据

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

func (svc *CachedUserService) FindOrCreate(ctx context.Context, phone string) (domain.User, error) {
	//这个是快路径
	domainu, err := svc.repo.FindByPhone(ctx, phone)
	//要判断有没有这个用户

	if !errors.Is(err, repository.ErrUserNotFound) {
		// nil会进来这里
		// 不为ErrUserNotFound也会进来这里
		return domainu, err
	}

	// 在系统资源不足，触发降级之后，不执行慢路径
	if ctx.Value("降级") == "true" {
		return domain.User{}, errors.New("系统降级了")
	}
	//这个是慢路径
	err = svc.repo.Create(ctx, domain.User{Phone: phone})
	if err != nil && !errors.Is(err, repository.ErrUserDuplicate) {
		return domain.User{}, err
	}
	//这里会遇到主从延迟的问题
	return svc.repo.FindByPhone(ctx, phone)
}

func PathDownGrade(ctx context.Context, quick, slow func()) {
	quick()
	if ctx.Value("降级") == "true" {
		return
	}
	slow()
}
