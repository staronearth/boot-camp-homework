package integration

import (
	"boot-camp-homework/webook/internal/web"
	"boot-camp-homework/webook/ioc"
	"bytes"
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestUserHandler_e2e_SendLoginSMSCode(t *testing.T) {
	server := InitWebServer()
	rdb := ioc.InitRedis()
	testCases := []struct {
		name string

		before  func(t *testing.T)
		after   func(t *testing.T)
		reqBody string

		wantCode int
		wantBody web.Result
	}{
		{
			name: "发送成功",
			before: func(t *testing.T) {
				//不需要，也就是redis什么数据也没有

			},
			after: func(t *testing.T) {
				ctx, cancle := context.WithTimeout(context.Background(), 3*time.Second)
				//你要清理数据
				//"phone_code:%s:%s"
				val, err := rdb.GetDel(ctx, "phone_code:login:15829074904").Result()
				cancle()
				assert.NoError(t, err)
				//你的验证码是6位
				assert.True(t, len(val) == 6)
			},
			reqBody: `
{
	"phone":"15829074904"
}
`,
			wantCode: http.StatusOK,
			wantBody: web.Result{
				Msg: "发送成功",
			},
		},
		{
			name: "发送成功",
			before: func(t *testing.T) {
				//不需要，也就是redis什么数据也没有

			},
			after: func(t *testing.T) {
				ctx, cancle := context.WithTimeout(context.Background(), 3*time.Second)
				//你要清理数据
				//"phone_code:%s:%s"
				val, err := rdb.GetDel(ctx, "phone_code:login:15829074904").Result()
				cancle()
				assert.NoError(t, err)
				//你的验证码是6位
				assert.True(t, len(val) == 6)
			},
			reqBody: `
{
	"phone":"15829074904"
}
`,
			wantCode: http.StatusOK,
			wantBody: web.Result{
				Msg: "发送成功",
			},
		},
		{
			name: "发送太频繁,请稍后再试",
			before: func(t *testing.T) {
				//不需要，也就是redis什么数据也没有
				ctx, cancle := context.WithTimeout(context.Background(), 3*time.Second)
				//你要清理数据
				//"phone_code:%s:%s"
				_, err := rdb.Set(ctx, "phone_code:login:15829074904", "123456", 570*time.Second).Result()
				cancle()
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancle := context.WithTimeout(context.Background(), 3*time.Second)
				//你要清理数据
				//"phone_code:%s:%s"
				val, err := rdb.GetDel(ctx, "phone_code:login:15829074904").Result()
				cancle()
				assert.NoError(t, err)
				//你的验证码是6位
				assert.Equal(t, "123456", val)
			},
			reqBody: `
{
	"phone":"15829074904"
}
`,
			wantCode: http.StatusOK,
			wantBody: web.Result{
				Code: 5,
				Msg:  "发送太频繁,请稍后再试",
			},
		},
		{
			name: "系统错误",
			before: func(t *testing.T) {
				//不需要，也就是redis什么数据也没有
				ctx, cancle := context.WithTimeout(context.Background(), 3*time.Second)
				//你要清理数据
				//"phone_code:%s:%s"
				_, err := rdb.Set(ctx, "phone_code:login:15829074904", "123456", 0).Result()
				cancle()
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancle := context.WithTimeout(context.Background(), 3*time.Second)
				//你要清理数据
				//"phone_code:%s:%s"
				val, err := rdb.GetDel(ctx, "phone_code:login:15829074904").Result()
				cancle()
				assert.NoError(t, err)
				//你的验证码是6位
				assert.Equal(t, "123456", val)
			},
			reqBody: `
{
	"phone":"15829074904"
}
`,
			wantCode: http.StatusOK,
			wantBody: web.Result{
				Code: 5,
				Msg:  "系统错误",
			},
		},
		{
			name: "输入有误",
			before: func(t *testing.T) {
			},
			after: func(t *testing.T) {
			},
			reqBody: `
{
	"phone":""
}
`,
			wantCode: http.StatusOK,
			wantBody: web.Result{
				Code: 4,
				Msg:  "输入有误",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			req, err := http.NewRequest(http.MethodPost, "/users/login_sms/code/send", bytes.NewBuffer([]byte(tc.reqBody)))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()

			server.ServeHTTP(resp, req)
			assert.Equal(t, tc.wantCode, resp.Code)
			if resp.Code != 200 {
				return
			}
			var webRes web.Result
			err = json.NewDecoder(resp.Body).Decode(&webRes)
			require.NoError(t, err)
			assert.Equal(t, tc.wantBody, webRes)
			tc.after(t)
		})
	}
}
