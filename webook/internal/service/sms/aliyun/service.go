package aliyun

import (
	"boot-camp-homework/webook/internal/service/sms"
	"context"
	"encoding/json"
	"fmt"
	dysmsapi "github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"strings"
)

type Service struct {
	client   *dysmsapi.Client
	signName string
}

func NewService(client *dysmsapi.Client, signName string) *Service {
	return &Service{
		client:   client,
		signName: signName,
	}
}
func (s *Service) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	return nil
	//sendSmsRequest := &aliyunsms.SendSmsRequest{
	//	SignName:      tea.String(s.signName),
	//	TemplateCode:  tea.String("SMS_462805477"),
	//	PhoneNumbers:  tea.String("15829074904"),
	//	TemplateParam: tea.String("{\"code\":\"988653\"}"),
	//}
	//runtime := &util.RuntimeOptions{}
	//tryErr := func() (_e error) {
	//	defer func() {
	//		if r := tea.Recover(recover()); r != nil {
	//			_e = r
	//		}
	//	}()
	//	// 复制代码运行请自行打印 API 的返回值
	//	_, _err := s.client.SendSmsWithOptions(sendSmsRequest, runtime)
	//	if _err != nil {
	//		return _err
	//	}
	//
	//	return nil
	//}()
	//
	//if tryErr != nil {
	//	var error = &tea.SDKError{}
	//	if _t, ok := tryErr.(*tea.SDKError); ok {
	//		error = _t
	//	} else {
	//		error.Message = tea.String(tryErr.Error())
	//	}
	//	// 如有需要，请打印 error
	//	_, _err := util.AssertAsString(error.Message)
	//	if _err != nil {
	//		return _err
	//	}
	//}
	//return nil
}

func (s *Service) SendV1(ctx context.Context, tplId string, args []sms.NamedArg, numbers ...string) error {
	req := dysmsapi.CreateSendSmsRequest()
	req.Scheme = "https"
	req.SignName = s.signName
	req.PhoneNumbers = strings.Join(numbers, ",")

	//传入的是Json
	argsMap := make(map[string]string, len(args))
	for _, arg := range args {
		argsMap[arg.Name] = arg.Val
	}

	val, err := json.Marshal(argsMap)
	if err != nil {
		return err
	}
	req.TemplateParam = string(val)
	req.TemplateCode = tplId

	resp, err := s.client.SendSms(req)
	if err != nil {
		return err
	}
	if resp.Code != "OK" {
		return fmt.Errorf("发送失败, code:%s,原因%s", resp.Code, resp.Message)
	}
	return nil
}
