package service

import (
	"boot-camp-homework/webook/internal/repository"
	"boot-camp-homework/webook/internal/service/sms"
	"context"
	"fmt"
	"log"
	"math/rand"
)

var (
	tplId                     = "SMS_462805477"
	ErrCodeSendTooMany        = repository.ErrCodeSendTooMany
	ErrCodeVerifyTooManyTimes = repository.ErrCodeVerifyTooManyTimes
)

type CodeService interface {
	Send(ctx context.Context, biz, phone string) error
	Verify(ctx context.Context, biz, phone, inputcode string) (bool, error)
}
type CachedCodeService struct {
	repo   repository.CodeRepository
	smsSvc sms.Service
}

func NewCodeService(repo repository.CodeRepository, smsSvc sms.Service) CodeService {
	return &CachedCodeService{
		repo:   repo,
		smsSvc: smsSvc,
	}
}

// Send 发送验证码，biz区别是什么业务,phone手机号
func (svc *CachedCodeService) Send(ctx context.Context, biz, phone string) error {
	//两个步骤
	//1.生成一个验证码
	code := svc.generateCode()
	//2.塞进去redis
	err := svc.repo.Store(ctx, biz, phone, code)
	if err != nil {
		//有问题
		return err
	}
	//3.发送验证码
	err = svc.smsSvc.SendV1(ctx, tplId, []sms.NamedArg{
		sms.NamedArg{
			Name: "code",
			Val:  code,
		},
	}, phone)
	if err != nil {
		//这个地方代表，redis有这个验证码，但是不好意思，你没发送成功，用户收不到
		//这个err可能是超时的err,你都不知道发出去了没有
		//在这里重试
		//要重试也是传入一个可重试的service
		//return err
		//记录日志
		log.Println(err)
	}
	return nil
}

func (svc *CachedCodeService) Verify(ctx context.Context, biz, phone, inputcode string) (bool, error) {
	// phone_code:login:188xxxxx
	// code:login:1xxxx
	// user:login:code:152xxx
	return svc.repo.Verify(ctx, biz, phone, inputcode)
}

func (svc *CachedCodeService) generateCode() string {
	num := rand.Intn(1000000)
	//不够六位，加上前导0
	return fmt.Sprintf("%06d", num)
}

func (svc *CachedCodeService) VerifyV1(ctx context.Context, biz, inputcode, phone string) error {
	return nil
}
