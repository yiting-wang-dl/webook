package failover

import (
	"context"
	"errors"
	"github.com/webook/internal/service/sms"
	"log"
	"sync/atomic"
)

type FailOverSMSService struct {
	svcs []sms.Service

	// v1
	// current provider index
	idx uint64
}

func NewFailOverSMSService(svcs []sms.Service) *FailOverSMSService {
	return &FailOverSMSService{
		svcs: svcs,
	}
}

func (f *FailOverSMSService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	for _, svc := range f.svcs {
		err := svc.Send(ctx, tplId, args, numbers...)
		if err == nil {
			return nil
		}
		log.Println(err)
	}
	return errors.New("round robin all providers, all failed to send")
}

func (f *FailOverSMSService) SendV1(ctx context.Context, tplId string, args []string, numbers ...string) error {
	idx := atomic.AddUint64(&f.idx, 1)
	length := uint64(len(f.svcs))
	// iterate length
	for i := idx; i < idx+length; i++ {
		// use mode to calculate idx
		svc := f.svcs[i%length]
		err := svc.Send(ctx, tplId, args, numbers...)
		switch err {
		case nil:
			return nil
		case context.Canceled, context.DeadlineExceeded:
			return err
		}
		log.Println(err)
	}
	return errors.New("round robin all providers, all failed to send")
}
