package ratelimit

import (
	"context"
	"errors"
	"github.com/webook/internal/service/sms"
	"github.com/webook/pkg/limiter"
)

var errLimited = errors.New("Trigger Rate Limit")

var _ sms.Service = &RateLimitSMSService{}

type RateLimitSMSService struct {
	// decorate
	svc     sms.Service
	limiter limiter.Limiter
	key     string
}

type RateLimitSMSServiceV1 struct {
	sms.Service
	limiter limiter.Limiter
	key     string
}

func (r *RateLimitSMSService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	limited, err := r.limiter.Limit(ctx, r.key)
	if err != nil {
		return err
	}
	if limited {
		return errLimited
	}
	return r.svc.Send(ctx, tplId, args, numbers...)
}

func NewRateLimitSMSService(svc sms.Service, l limiter.Limiter) *RateLimitSMSService {
	return &RateLimitSMSService{
		svc:     svc,
		limiter: l,
		key:     "sms-limiter",
	}
}
