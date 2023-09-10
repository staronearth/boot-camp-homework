package retryable

import (
	"boot-camp-homework/webook/internal/service/sms"
	"context"
)

// Service 这个要小心并发问题
type Service struct {
	svc      sms.Service
	retryCnt int
}

func (s *Service) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	//TODO implement me
	panic("implement me")
}

func (s *Service) SendV1(ctx context.Context, tplId string, args []sms.NamedArg, numbers ...string) error {
	err := s.svc.SendV1(ctx, tplId, args, numbers...)
	for err != nil && s.retryCnt < 3 {
		err = s.svc.SendV1(ctx, tplId, args, numbers...)
		s.retryCnt++
	}
	return nil
}
