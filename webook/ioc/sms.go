package ioc

import (
	"boot-camp-homework/webook/internal/service/sms"
	"boot-camp-homework/webook/internal/service/sms/aliyun"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"os"
)

func InitSMSService() sms.Service {
	client, err := dysmsapi.NewClientWithAccessKey("cn-shanghai", os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_ID"), os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_SECRET"))

	if err != nil {
		panic(err)
	}

	return aliyun.NewService(client, "我的webook")
}
