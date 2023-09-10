package service

import (
	"boot-camp-homework/webook/internal/domain"
	"boot-camp-homework/webook/internal/repository"
	repomocks "boot-camp-homework/webook/internal/repository/mocks"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"testing"
	"time"
)

func TestCachedUserService_Login(t *testing.T) {
	now := time.Now()
	testCases := []struct {
		name string

		mock func(ctrl *gomock.Controller) repository.UserRepository

		ctx       context.Context
		inputUser domain.User

		wantUser domain.User

		wantErr error
	}{
		{
			name: "登录成功",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				userrepo := repomocks.NewMockUserRepository(ctrl)
				userrepo.EXPECT().FindByEmail(gomock.Any(), domain.User{
					Email:    "123@qq.com",
					Password: "hello#world123",
				}).Return(domain.User{
					Email:    "123@qq.com",
					Password: "$2a$10$sufBzJzpFEDazKwZ1JRIvexRHDzVTL9zpLaJW4vMKp9P17RxicgxK",
					Phone:    "13333322233",
					Ctime:    now,
				}, nil)
				return userrepo
			},

			inputUser: domain.User{
				Email:    "123@qq.com",
				Password: "hello#world123",
			},

			wantUser: domain.User{
				Email:    "123@qq.com",
				Password: "$2a$10$sufBzJzpFEDazKwZ1JRIvexRHDzVTL9zpLaJW4vMKp9P17RxicgxK",
				Phone:    "13333322233",
				Ctime:    now,
			},
			wantErr: nil,
		},
		{
			name: "用户不存在",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				userrepo := repomocks.NewMockUserRepository(ctrl)
				userrepo.EXPECT().FindByEmail(gomock.Any(), domain.User{
					Email:    "123@qq.com",
					Password: "hello#world123",
				}).Return(domain.User{}, repository.ErrUserNotFound)
				return userrepo
			},

			inputUser: domain.User{
				Email:    "123@qq.com",
				Password: "hello#world123",
			},

			wantUser: domain.User{},
			wantErr:  ErrInvalidUserOrPassword,
		},
		{
			name: "DB的err",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				userrepo := repomocks.NewMockUserRepository(ctrl)
				userrepo.EXPECT().FindByEmail(gomock.Any(), domain.User{
					Email:    "123@qq.com",
					Password: "hello#world123",
				}).Return(domain.User{}, errors.New("db的err"))
				return userrepo
			},

			inputUser: domain.User{
				Email:    "123@qq.com",
				Password: "hello#world123",
			},

			wantUser: domain.User{},
			wantErr:  errors.New("db的err"),
		},
		{
			name: "密码不匹配",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				userrepo := repomocks.NewMockUserRepository(ctrl)
				userrepo.EXPECT().FindByEmail(gomock.Any(), domain.User{
					Email:    "123@qq.com",
					Password: "hello#world1234",
				}).Return(domain.User{
					Email:    "123@qq.com",
					Password: "$2a$10$sufBzJzpFEDazKwZ1JRIvexRHDzVTL9zpLaJW4vMKp9P17RxicgxK",
					Phone:    "13333322233",
					Ctime:    now,
				}, nil)
				return userrepo
			},

			inputUser: domain.User{
				Email:    "123@qq.com",
				Password: "hello#world1234",
			},

			wantUser: domain.User{},
			wantErr:  ErrInvalidUserOrPassword,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			usvc := NewUserService(tc.mock(ctrl))
			resuser, err := usvc.Login(context.Background(), tc.inputUser)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantUser, resuser)
		})
	}
}

func TestBcrpty(t *testing.T) {
	passwd, err := bcrypt.GenerateFromPassword([]byte("hello#world123"), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	println(string(passwd))
}
