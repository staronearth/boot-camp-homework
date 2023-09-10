package memory

import (
	"boot-camp-homework/webook/internal/service/sms"
	"context"
	"fmt"
)

type Service struct {
}

func (s Service) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	fmt.Println(args)
	return nil
}

func (s Service) SendV1(ctx context.Context, tplId string, args []sms.NamedArg, numbers ...string) error {
	return nil
}
