package repository

import (
	"boot-camp-homework/webook/internal/domain"
	"boot-camp-homework/webook/internal/repository/cache"
	cachemocks "boot-camp-homework/webook/internal/repository/cache/mocks"
	"boot-camp-homework/webook/internal/repository/dao"
	daomocks "boot-camp-homework/webook/internal/repository/dao/mocks"
	"context"
	"database/sql"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func TestCachedUserRepository_FindUserTableById(t *testing.T) {
	now := time.UnixMilli(time.Now().UnixMilli())
	testCases := []struct {
		name string

		mock func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache)

		ctx context.Context
		idx int64

		wantUser domain.User
		wantErr  error
	}{
		{
			name: "缓存未命中,数据库加载返回数值",

			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				daomock := daomocks.NewMockUserDAO(ctrl)
				cachemock := cachemocks.NewMockUserCache(ctrl)

				cachemock.EXPECT().Get(gomock.Any(), int64(1)).Return(domain.User{}, errors.New("mock cache出错的情况"))
				daomock.EXPECT().FindUserTableById(gomock.Any(), int64(1)).Return(dao.User{
					Id:       int64(1),
					Password: "123456",
					Email: sql.NullString{
						String: "123@qq.com",
						Valid:  true,
					},
					Phone: sql.NullString{
						String: "123223456676",
						Valid:  true,
					},
					Ctime: now.UnixMilli(),
					Utime: now.UnixMilli(),
				}, nil)
				cachemock.EXPECT().Set(gomock.Any(), domain.User{
					Id:       int64(1),
					Password: "123456",
					Email:    "123@qq.com",
					Phone:    "123223456676",
					Ctime:    now,
				}).Return(nil)
				return daomock, cachemock
			},
			idx: int64(1),
			wantUser: domain.User{
				Id:       int64(1),
				Password: "123456",
				Email:    "123@qq.com",
				Phone:    "123223456676",
				Ctime:    now,
			},
			wantErr: nil,
		},
		{
			name: "缓存命中",

			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				daomock := daomocks.NewMockUserDAO(ctrl)
				cachemock := cachemocks.NewMockUserCache(ctrl)

				cachemock.EXPECT().Get(gomock.Any(), int64(1)).Return(domain.User{
					Id:       int64(1),
					Password: "123456",
					Email:    "123@qq.com",
					Phone:    "123223456676",
					Ctime:    now,
				}, nil)
				return daomock, cachemock
			},
			idx: int64(1),
			wantUser: domain.User{
				Id:       int64(1),
				Password: "123456",
				Email:    "123@qq.com",
				Phone:    "123223456676",
				Ctime:    now,
			},
			wantErr: nil,
		},
		{
			name: "缓存未命中,查询失败",

			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				daomock := daomocks.NewMockUserDAO(ctrl)
				cachemock := cachemocks.NewMockUserCache(ctrl)

				cachemock.EXPECT().Get(gomock.Any(), int64(1)).Return(domain.User{}, errors.New("mock cache出错的情况"))
				daomock.EXPECT().FindUserTableById(gomock.Any(), int64(1)).Return(dao.User{}, errors.New("mock db出错"))
				return daomock, cachemock
			},
			idx:      int64(1),
			wantUser: domain.User{},
			wantErr:  errors.New("mock db出错"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			urepo := NewUserRepository(tc.mock(ctrl))
			resuser, err := urepo.FindUserTableById(context.Background(), tc.idx)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantUser, resuser)
		})
	}
}
