package repository

import (
	"boot-camp-homework/webook/internal/domain"
	"boot-camp-homework/webook/internal/repository/cache"
	"boot-camp-homework/webook/internal/repository/dao"
	"context"
	"database/sql"
	"log"
	"time"
)

var (
	ErrUserDuplicate    = dao.ErrUserDuplicate
	ErrUserNotFound     = dao.ErrUserNotFound
	ErrUserInfoNotFound = dao.ErrUserInfoNotFound
)

// UserRepository 是核心，它有不同的实现。丹斯 Factory 本身如果只是初始化一下
// 那么它不是你的核心
type UserRepository interface {
	FindByEmail(ctx context.Context, u domain.User) (domain.User, error)
	FindByPhone(ctx context.Context, phone string) (domain.User, error)
	Create(ctx context.Context, u domain.User) error
	FindUserTableById(ctx context.Context, idx int64) (domain.User, error)
	FindUserInfoTableById(ctx context.Context, idx int64) (domain.UserInfo, error)
	CreateUserInfo(ctx context.Context, info domain.UserInfo) error
	UpdateUserInfo(ctx context.Context, info domain.UserInfo) error
}
type CachedUserRepository struct {
	dao   dao.UserDAO
	cache cache.UserCache
}

func NewUserRepository(dao dao.UserDAO, c cache.UserCache) UserRepository {
	return &CachedUserRepository{
		dao:   dao,
		cache: c,
	}
}
func (ur *CachedUserRepository) FindByEmail(ctx context.Context, u domain.User) (domain.User, error) {
	daou, err := ur.dao.FindByEmail(ctx, u.Email)
	if err != nil {
		return domain.User{}, err
	}
	return ur.entityToDomain(daou), nil
}

func (ur *CachedUserRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	daou, err := ur.dao.FindByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}
	return ur.entityToDomain(daou), nil
}
func (ur *CachedUserRepository) Create(ctx context.Context, u domain.User) error {
	return ur.dao.Insert(ctx, ur.domainToEntity(u))
	//在这里操作缓存
}

func (ur *CachedUserRepository) FindUserTableById(ctx context.Context, idx int64) (domain.User, error) {

	//先从cache里面找
	//第一种，缓存里面有数据
	//第二种，缓存里面没有数据
	//第三种，缓存出错了，你也不知道有没有数据
	cacheu, err := ur.cache.Get(ctx, idx)
	if err == nil {
		//缓存必然有数据
		return cacheu, nil
	}
	//if errors.Is(err, cache.ErrKeyNotExist) {
	//	//去数据中加载
	//}
	//选加载-做好兜底，万一redis真的崩了，你要保护住你的数据库
	//我数据库限流

	//选不加载-用户体验差
	//再从dao里面找
	daou, err := ur.dao.FindUserTableById(ctx, idx)
	if err != nil {
		return domain.User{}, err
	}
	domainu := ur.entityToDomain(daou)
	//找到了回写cache
	err = ur.cache.Set(ctx, domainu)
	if err != nil {
		//做好监控
		log.Println("缓存设置失败")
	}
	//go func() {
	//
	//}()

	return domainu, nil
}

func (ur *CachedUserRepository) FindUserInfoTableById(ctx context.Context, idx int64) (domain.UserInfo, error) {

	//先从cache里面找
	//再从dao里面找
	daoui, err := ur.dao.FindUserInfoTableById(ctx, idx)
	if err != nil {
		return domain.UserInfo{}, ErrUserInfoNotFound
	}
	return domain.UserInfo{
		Id:              daoui.Id,
		NickName:        daoui.NickName,
		BrithDays:       daoui.BrithDays,
		PersonalProfile: daoui.PersonalProfile,
	}, nil
	//找到了回写cache
}
func (ur *CachedUserRepository) CreateUserInfo(ctx context.Context, info domain.UserInfo) error {
	return ur.dao.InsertUserInfo(ctx, dao.UserInfo{
		Id:              info.Id,
		NickName:        info.NickName,
		BrithDays:       info.BrithDays,
		PersonalProfile: info.PersonalProfile,
	})
}

func (ur *CachedUserRepository) UpdateUserInfo(ctx context.Context, info domain.UserInfo) error {
	return ur.dao.UpdateUserInfo(ctx, dao.UserInfo{
		Id:              info.Id,
		NickName:        info.NickName,
		BrithDays:       info.BrithDays,
		PersonalProfile: info.PersonalProfile,
	})
}
func (ur *CachedUserRepository) domainToEntity(u domain.User) dao.User {
	return dao.User{
		Id: u.Id,
		Email: sql.NullString{
			String: u.Email,
			Valid:  u.Email != "",
		},
		Password: u.Password,
		Phone: sql.NullString{
			String: u.Phone,
			Valid:  u.Phone != "",
		},
		Ctime: u.Ctime.UnixMilli(),
	}
}
func (ur *CachedUserRepository) entityToDomain(u dao.User) domain.User {
	return domain.User{
		Id:       u.Id,
		Email:    u.Email.String,
		Password: u.Password,
		Phone:    u.Phone.String,
		Ctime:    time.UnixMilli(u.Ctime),
	}
}
