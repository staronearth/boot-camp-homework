package web

import (
	"boot-camp-homework/webook/internal/domain"
	"boot-camp-homework/webook/internal/service"
	svcmocks "boot-camp-homework/webook/internal/service/mocks"
	jwtmocks "boot-camp-homework/webook/internal/web/mocks"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestEncrypt(t *testing.T) {
	bgpgyte, err := bcrypt.GenerateFromPassword([]byte("gaolin@123"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatal(err)
	}
	err = bcrypt.CompareHashAndPassword(bgpgyte, []byte("gaolin@123"))
	assert.NoError(t, err)
}

func TestBirthDays(t *testing.T) {
	birthday := "4000-12-31" // 替换为要验证的生日日期
	if isValidBirthday(birthday) {
		fmt.Println("生日有效")
	} else {
		fmt.Println("生日无效")
	}
}
func isValidBirthday(date string) bool {
	layout := "2006-01-02" // 指定日期格式

	// 解析日期字符串
	birthday, err := time.Parse(layout, date)
	if err != nil {
		return false
	}

	// 获取当前日期
	currentDate := time.Now()

	// 比较生日是否在当前日期之前
	if birthday.Before(currentDate) {
		return true
	}

	return false
}

func TestUserHandler_Signup(t *testing.T) {
	testCases := []struct {
		name string

		mock func(ctrl *gomock.Controller) service.UserService

		reqBody string

		wantCode int
		wantBody string
	}{
		{
			name: "注册成功",
			mock: func(ctrl *gomock.Controller) service.UserService {
				usersvc := svcmocks.NewMockUserService(ctrl)
				usersvc.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "123@qq.com",
					Password: "hello#world123",
				}).Return(nil)
				return usersvc
			},

			reqBody: `
{
"email":"123@qq.com",
"password":"hello#world123",
"confirmPassword":"hello#world123"
}
`,
			wantCode: http.StatusOK,
			wantBody: "注册成功",
		},
		{
			name: "参数不对，bind失败",
			mock: func(ctrl *gomock.Controller) service.UserService {
				usersvc := svcmocks.NewMockUserService(ctrl)
				return usersvc
			},

			reqBody: `
{
"email":"123@qq.com",
"password":"hello#world123",
"confirmPassword":"hello#world123",
}
`,
			wantCode: http.StatusBadRequest,
		},
		{
			name: "邮箱格式不对",
			mock: func(ctrl *gomock.Controller) service.UserService {
				usersvc := svcmocks.NewMockUserService(ctrl)
				return usersvc
			},

			reqBody: `
{
"email":"123qq.com",
"password":"hello#world123",
"confirmPassword":"hello#world123"
}
`,
			wantCode: http.StatusOK,
			wantBody: "你的邮箱格式不对",
		},
		{
			name: "两次输入的密码不一致",
			mock: func(ctrl *gomock.Controller) service.UserService {
				usersvc := svcmocks.NewMockUserService(ctrl)
				return usersvc
			},

			reqBody: `
{
"email":"123@qq.com",
"password":"helloworld123",
"confirmPassword":"hello#world123"
}
`,
			wantCode: http.StatusOK,
			wantBody: "两次输入的密码不一致",
		},
		{
			name: "密码必须大于8位，包含数字、特殊字符",
			mock: func(ctrl *gomock.Controller) service.UserService {
				usersvc := svcmocks.NewMockUserService(ctrl)
				return usersvc
			},

			reqBody: `
{
"email":"123@qq.com",
"password":"hello",
"confirmPassword":"hello"
}
`,
			wantCode: http.StatusOK,
			wantBody: "密码必须大于8位，包含数字、特殊字符",
		},
		{
			name: "邮箱冲突",
			mock: func(ctrl *gomock.Controller) service.UserService {
				usersvc := svcmocks.NewMockUserService(ctrl)
				usersvc.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "123@qq.com",
					Password: "hello#world123",
				}).Return(service.ErrUserDuplicate)
				return usersvc
			},

			reqBody: `
{
"email":"123@qq.com",
"password":"hello#world123",
"confirmPassword":"hello#world123"
}
`,
			wantCode: http.StatusOK,
			wantBody: "邮箱冲突",
		},
		{
			name: "系统异常",
			mock: func(ctrl *gomock.Controller) service.UserService {
				usersvc := svcmocks.NewMockUserService(ctrl)
				usersvc.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "123@qq.com",
					Password: "hello#world123",
				}).Return(errors.New("随便一个error"))
				return usersvc
			},

			reqBody: `
{
"email":"123@qq.com",
"password":"hello#world123",
"confirmPassword":"hello#world123"
}
`,
			wantCode: http.StatusOK,
			wantBody: "系统错误",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			server := gin.Default()
			//用不上codeSvc

			h := NewUserHandler(tc.mock(ctrl), nil, nil)
			h.RegisterRoutes(server)

			req, err := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewBuffer([]byte(tc.reqBody)))
			require.NoError(t, err)
			//数据事json格式的
			req.Header.Set("Content-Type", "application/json")
			// 这里就可以继续使用req

			resp := httptest.NewRecorder()
			//这里就是HTTP请求进入GIN框架的入口
			// 当你这样调用的时候，GIN就会处理这个请求
			//响应写会resp里

			server.ServeHTTP(resp, req)

			assert.Equal(t, tc.wantCode, resp.Code)
			assert.Equal(t, tc.wantBody, resp.Body.String())
		})
	}
}

func TestUserHandler_LoginSMS(t *testing.T) {

	testCases := []struct {
		name      string
		usvcmock  func(ctrl *gomock.Controller) service.UserService
		csvcmock  func(ctrl *gomock.Controller) service.CodeService
		tokenmock func(ctrl *gomock.Controller) Token
		reqBody   string

		wantCode int
		wantRes  Result
	}{
		{
			name: "验证码校验通过",
			usvcmock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmocks.NewMockUserService(ctrl)
				userSvc.EXPECT().FindOrCreate(gomock.Any(), "15243994399").Return(domain.User{
					Id:    int64(1),
					Phone: "15243994399",
				}, nil)
				return userSvc
			},
			csvcmock: func(ctrl *gomock.Controller) service.CodeService {
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				codeSvc.EXPECT().Verify(gomock.Any(), "login", "15243994399", "558669").Return(true, nil)
				return codeSvc
			},
			tokenmock: func(ctrl *gomock.Controller) Token {
				jwth := jwtmocks.NewMockToken(ctrl)
				jwth.EXPECT().SetJWTToken(gomock.Any(), int64(1)).Return(nil)
				return jwth
			},
			reqBody: `
		{
			"phone":"15243994399",
			"code":"558669"
		}
		`,
			wantCode: http.StatusOK,
			wantRes: Result{
				Code: 4,
				Msg:  "验证码校验通过",
			},
		},
		{
			name: "设置token失败",
			usvcmock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmocks.NewMockUserService(ctrl)
				userSvc.EXPECT().FindOrCreate(gomock.Any(), "15243994399").Return(domain.User{
					Id:    int64(1),
					Phone: "15243994399",
				}, nil)
				return userSvc
			},
			csvcmock: func(ctrl *gomock.Controller) service.CodeService {
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				codeSvc.EXPECT().Verify(gomock.Any(), "login", "15243994399", "558669").Return(true, nil)
				return codeSvc
			},
			tokenmock: func(ctrl *gomock.Controller) Token {
				jwth := jwtmocks.NewMockToken(ctrl)
				jwth.EXPECT().SetJWTToken(gomock.Any(), int64(1)).Return(errors.New("设置token失败"))
				return jwth
			},
			reqBody: `
		{
			"phone":"15243994399",
			"code":"558669"
		}
		`,
			wantCode: http.StatusOK,
			wantRes: Result{
				Code: 5,
				Msg:  "系统错误",
			},
		},
		{
			name: "findorcreate报错",
			usvcmock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmocks.NewMockUserService(ctrl)
				userSvc.EXPECT().FindOrCreate(gomock.Any(), "15243994399").Return(domain.User{}, errors.New("系统降级了"))
				return userSvc
			},
			csvcmock: func(ctrl *gomock.Controller) service.CodeService {
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				codeSvc.EXPECT().Verify(gomock.Any(), "login", "15243994399", "558669").Return(true, nil)
				return codeSvc
			},
			tokenmock: func(ctrl *gomock.Controller) Token {
				jwth := jwtmocks.NewMockToken(ctrl)
				return jwth
			},
			reqBody: `
		{
			"phone":"15243994399",
			"code":"558669"
		}
		`,
			wantCode: http.StatusOK,
			wantRes: Result{
				Code: 5,
				Msg:  "系统错误",
			},
		},
		{
			name: "验证码有误",
			usvcmock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmocks.NewMockUserService(ctrl)
				return userSvc
			},
			csvcmock: func(ctrl *gomock.Controller) service.CodeService {
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				codeSvc.EXPECT().Verify(gomock.Any(), "login", "15243994399", "558669").Return(false, nil)
				return codeSvc
			},
			tokenmock: func(ctrl *gomock.Controller) Token {
				jwth := jwtmocks.NewMockToken(ctrl)
				return jwth
			},
			reqBody: `
		{
			"phone":"15243994399",
			"code":"558669"
		}
		`,
			wantCode: http.StatusOK,
			wantRes: Result{
				Code: 4,
				Msg:  "验证码有误",
			},
		},
		{
			name: "验证码有误",
			usvcmock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmocks.NewMockUserService(ctrl)
				return userSvc
			},
			csvcmock: func(ctrl *gomock.Controller) service.CodeService {
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				codeSvc.EXPECT().Verify(gomock.Any(), "login", "15243994399", "558669").Return(false, service.ErrCodeVerifyTooManyTimes)
				return codeSvc
			},
			tokenmock: func(ctrl *gomock.Controller) Token {
				jwth := jwtmocks.NewMockToken(ctrl)
				return jwth
			},
			reqBody: `
		{
			"phone":"15243994399",
			"code":"558669"
		}
		`,
			wantCode: http.StatusOK,
			wantRes: Result{
				Code: 5,
				Msg:  "系统错误",
			},
		},
		{
			name: "参数bind错误",
			usvcmock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmocks.NewMockUserService(ctrl)
				return userSvc
			},
			csvcmock: func(ctrl *gomock.Controller) service.CodeService {
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				return codeSvc
			},
			tokenmock: func(ctrl *gomock.Controller) Token {
				jwth := jwtmocks.NewMockToken(ctrl)
				return jwth
			},
			reqBody: `
		{
			"phone":"15243994399",
			"code":"558669".
		}
		`,
			wantCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			server := gin.Default()
			h := NewUserHandler(tc.usvcmock(ctrl), tc.csvcmock(ctrl), tc.tokenmock(ctrl))
			h.RegisterRoutes(server)
			req, err := http.NewRequest(http.MethodPost, "/users/login_sms", bytes.NewBuffer([]byte(tc.reqBody)))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()

			server.ServeHTTP(resp, req)
			assert.Equal(t, tc.wantCode, resp.Code)
			if resp.Code != 200 {
				return
			}
			var webRes Result
			err = json.NewDecoder(resp.Body).Decode(&webRes)
			require.NoError(t, err)
			assert.Equal(t, tc.wantRes, webRes)
		})
	}
}

func TestUserHandler_LoginJWT(t *testing.T) {
}
