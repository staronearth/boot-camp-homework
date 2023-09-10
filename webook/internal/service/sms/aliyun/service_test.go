package aliyun

import (
	"boot-camp-homework/webook/internal/service/sms"
	"context"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi "github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/ecodeclub/ekit"
	"os"

	"github.com/stretchr/testify/assert"

	"testing"
)

func TestService_Send(t *testing.T) {
	client, _err := CreateClient(ekit.ToPtr[string](os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_ID")), ekit.ToPtr[string](os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_SECRET")))
	if _err != nil {
		t.Fatal(_err)
	}
	signName := ""
	s := NewService(client, signName)
	testCases := []struct {
		name    string
		tplId   string
		params  []string
		numbers []string
		wantErr error
	}{
		{
			name:   "发送验证码",
			tplId:  "",
			params: []string{""},
			//修改为你的手机号
			numbers: []string{"15829074904"},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			er := s.Send(context.Background(), tc.tplId, tc.params, tc.numbers...)
			assert.Equal(t, er, tc.wantErr)
		})
	}
}
func TestService_SendV1(t *testing.T) {
	client, err := dysmsapi.NewClientWithAccessKey("cn-shanghai", os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_ID"), os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_SECRET"))

	if err != nil {
		t.Fatal(err)
	}
	signName := "我的webook"
	s := NewService(client, signName)
	testCases := []struct {
		name    string
		tplId   string
		params  []sms.NamedArg
		numbers []string
		wantErr error
	}{
		{
			name:  "发送验证码",
			tplId: "SMS_462805477",

			params: []sms.NamedArg{
				sms.NamedArg{
					Name: "code",
					Val:  "654332",
				},
			},
			//修改为你的手机号
			numbers: []string{os.Getenv("PHONE_NUMBER1")},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			er := s.SendV1(context.Background(), tc.tplId, tc.params, tc.numbers...)
			assert.Equal(t, er, tc.wantErr)
		})
	}
}
func CreateClient(accessKeyId *string, accessKeySecret *string) (_result *dysmsapi.Client, _err error) {
	config := &openapi.Config{
		// 必填，您的 AccessKey ID
		AccessKeyId: accessKeyId,
		// 必填，您的 AccessKey Secret
		AccessKeySecret: accessKeySecret,
	}
	// Endpoint 请参考 https://api.aliyun.com/product/Dysmsapi
	config.Endpoint = ekit.ToPtr[string]("dysmsapi.aliyuncs.com")
	_result = &dysmsapi.Client{}
	_result, _err = dysmsapi.NewClient()
	return _result, _err
}
