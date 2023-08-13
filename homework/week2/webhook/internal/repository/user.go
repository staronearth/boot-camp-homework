package repository

import (
	"boot-camp-homework/homework/week2/webhook/internal/domain"
	"boot-camp-homework/homework/week2/webhook/internal/repository/dao"
	"context"
)

var (
	ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
	ErrUserNotFound       = dao.ErrUserNotFound
	ErrUserInfoNotFound   = dao.ErrUserInfoNotFound
)

type UserRepository struct {
	dao *dao.UserDAO
}

func NewUserRepository(dao *dao.UserDAO) *UserRepository {
	return &UserRepository{
		dao: dao,
	}
}
func (ur UserRepository) FindByEmail(ctx context.Context, u domain.User) (domain.User, error) {
	daou, err := ur.dao.FindByEmail(ctx, u.Email)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		Id:       daou.Id,
		Email:    daou.Email,
		Password: daou.Password,
	}, nil
}
func (ur *UserRepository) Create(ctx context.Context, u domain.User) error {
	return ur.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})
	//在这里操作缓存
}

func (ur *UserRepository) FindUserTableById(ctx context.Context, idx int64) (domain.User, error) {

	//先从cache里面找
	//再从dao里面找
	daou, err := ur.dao.FindUserTableById(ctx, idx)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		Id:       daou.Id,
		Email:    daou.Email,
		Password: daou.Password,
	}, nil
	//找到了回写cache
}

func (ur *UserRepository) FindUserInfoTableById(ctx context.Context, idx int64) (domain.UserInfo, error) {

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
func (ur UserRepository) CreateUserInfo(ctx context.Context, info domain.UserInfo) error {
	return ur.dao.InsertUserInfo(ctx, dao.UserInfo{
		Id:              info.Id,
		NickName:        info.NickName,
		BrithDays:       info.BrithDays,
		PersonalProfile: info.PersonalProfile,
	})
}

func (ur UserRepository) UpdateUserInfo(ctx context.Context, info domain.UserInfo) error {
	return ur.dao.UpdateUserInfo(ctx, dao.UserInfo{
		Id:              info.Id,
		NickName:        info.NickName,
		BrithDays:       info.BrithDays,
		PersonalProfile: info.PersonalProfile,
	})
}
