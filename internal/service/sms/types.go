package sms

import "context"

// Service is the abstract to send msg
// different carriers will require different payload, here we make it as generic as possible but not covering all cases
type Service interface {
	Send(ctx context.Context, tplId string, args []string, numbers ...string) error
}
