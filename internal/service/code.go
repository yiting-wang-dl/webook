package service

import (
	"context"
	"github.com/webook/internal/repository"
	"github.com/webook/internal/service/sms"
)

var ErrCodeSendTooMany = repository.ErrCodeSendTooMany

type CodeService interface {
}

type codeService struct {
	repo repository.CodeRepository
	sms  sms.Service
}

func (svc *CodeService) Send(ctx context.Context, biz, phone string) error {
	return nil
}

func (svc *CodeService) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	return ok, err
}
