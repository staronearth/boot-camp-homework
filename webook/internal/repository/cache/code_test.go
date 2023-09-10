package cache

import (
	"boot-camp-homework/webook/internal/repository/cache/redismocks"
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestRedisCodeCache_Set(t *testing.T) {
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) redis.Cmdable

		ctx   context.Context
		biz   string
		phone string
		code  string

		wantErr error
	}{
		{
			name: "redis set成功，res为0",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmdable := redismocks.NewMockCmdable(ctrl)
				res := redis.NewCmd(context.Background())
				res.SetVal(int64(0))
				cmdable.EXPECT().Eval(gomock.Any(), luaSetCode, []string{"phone_code:login:151223344"}, []any{"775633"}).Return(res)
				return cmdable
			},
			biz:     "login",
			phone:   "151223344",
			code:    "775633",
			wantErr: nil,
		},
		{
			name: "发送太多次数，res为-1",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmdable := redismocks.NewMockCmdable(ctrl)
				res := redis.NewCmd(context.Background())
				res.SetVal(int64(-1))
				cmdable.EXPECT().Eval(gomock.Any(), luaSetCode, []string{"phone_code:login:151223344"}, []any{"775633"}).Return(res)
				return cmdable
			},
			biz:     "login",
			phone:   "151223344",
			code:    "775633",
			wantErr: ErrCodeSendTooMany,
		},
		{
			name: "redis set成功，res为-2",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmdable := redismocks.NewMockCmdable(ctrl)
				res := redis.NewCmd(context.Background())
				res.SetVal(int64(-2))
				cmdable.EXPECT().Eval(gomock.Any(), luaSetCode, []string{"phone_code:login:151223344"}, []any{"775633"}).Return(res)
				return cmdable
			},
			biz:     "login",
			phone:   "151223344",
			code:    "775633",
			wantErr: errors.New("系统错误"),
		},
		{
			name: "redis错误",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmdable := redismocks.NewMockCmdable(ctrl)
				res := redis.NewCmd(context.Background())
				res.SetErr(errors.New("mock redis错误"))
				cmdable.EXPECT().Eval(gomock.Any(), luaSetCode, []string{"phone_code:login:151223344"}, []any{"775633"}).Return(res)
				return cmdable
			},
			biz:     "login",
			phone:   "151223344",
			code:    "775633",
			wantErr: errors.New("mock redis错误"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			codecache := NewCodeCache(tc.mock(ctrl))
			err := codecache.Set(context.Background(), tc.biz, tc.phone, tc.code)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
