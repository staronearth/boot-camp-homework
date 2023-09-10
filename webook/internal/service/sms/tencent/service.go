package tencent

import (
	"context"
	"fmt"
	"github.com/ecodeclub/ekit"
	"github.com/ecodeclub/ekit/slice"
	tecentsms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

type Service struct {
	appId    *string
	signName *string
	client   *tecentsms.Client
}

func NewService(client *tecentsms.Client, appId string, signName string) *Service {
	return &Service{
		client:   client,
		appId:    ekit.ToPtr[string](appId),
		signName: ekit.ToPtr[string](signName),
	}
}
func (s *Service) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	req := tecentsms.NewSendSmsRequest()
	req.SmsSdkAppId = s.appId
	req.SignName = s.signName
	req.TemplateId = ekit.ToPtr[string](tplId)
	req.PhoneNumberSet = s.toStringPtrSlice(numbers)
	req.TemplateParamSet = s.toStringPtrSlice(args)

	resp, err := s.client.SendSms(req)
	if err != nil {
		return err
	}
	for _, sendStatus := range resp.Response.SendStatusSet {
		if sendStatus.Code == nil || *(sendStatus.Code) != "OK" {
			return fmt.Errorf("发送短信失败%s,%s", *sendStatus.Code, *sendStatus.Message)
		}
	}
	return nil
}

func (s *Service) toStringPtrSlice(src []string) []*string {
	return slice.Map[string, *string](src, func(idx int, src string) *string {
		return &src
	})
}
