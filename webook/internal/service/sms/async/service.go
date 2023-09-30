package async

import (
	"boot-camp/webook/internal/repository"
	"boot-camp/webook/internal/service/sms"
	"context"
)

type HealthSMSService struct {
	svc  []sms.Service
	repo repository.SMSAsyncReqRepository
	idx  int64
}

func (h *HealthSMSService) Send(ctx context.Context, biz string, args []string, numbers ...string) error {
	//TODO implement me
	panic("implement me")
}

func (h *HealthSMSService) SendV1(ctx context.Context, biz string, args []sms.NamedArg, numbers ...string) error {
	svc := h.svc[h.idx]
	err := svc.SendV1(ctx, biz, args, numbers...)
	if err != nil {

		//判断运营商是否崩溃了
		//if 崩溃了 {
		//	h.repo.Strore(ctx)
		//}
	}
	return nil
}

func (h *HealthSMSService) StartAsync(ctx context.Context) {
	go func() {
		reqs := h.repo.FindNotSendReq(ctx)
		for _, req := range reqs {
			//在这里发送并且控制重试
			h.SendV1(ctx, req.Biz, req.Args, req.Numbers...)
		}
	}()
}
func NewealthSMSService(svc []sms.Service, repo repository.SMSAsyncReqRepository) sms.Service {
	return &HealthSMSService{
		svc:  svc,
		repo: repo,
	}
}
