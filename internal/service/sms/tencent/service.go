package tencent

import (
	"context"
	"fmt"
	"github.com/ecodeclub/ekit"
	"github.com/ecodeclub/ekit/slice"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

type Service struct {
	client   *sms.Client
	appId    *string
	signName *string
}

func (s Service) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	request := sms.NewSendSmsRequest()
	request.SetContext(ctx)
	request.SmsSdkAppId = s.appId
	request.SignName = s.signName
	request.TemplateId = ekit.ToPtr[string](tplId)
	request.TemplateParamSet = s.toPtrSlice(args)
	request.PhoneNumberSet = s.toPtrSlice(numbers)
	response, err := s.client.SendSms(request)
	// handle errors
	if err != nil {
		return err
	}
	for _, statusPtr := range response.Response.SendStatusSet {
		if statusPtr == nil {
			// impossible to enter this route, but sine this is a pointer, need to check if the pointer is nil
			continue
		}
		status := *statusPtr
		if status.Code == nil || *(status.Code) != "Ok" {
			// Fail to Send
			// if status.Code == nil, *status.Code will return error. assuming status.Code != nil
			return fmt.Errorf("Message Failed To Send. code: %s, msg: %s", *status.Code, *status.Message)
		}
	}
	return nil

}

func (s *Service) toPtrSlice(data []string) []*string {
	return slice.Map[string, *string](data,
		func(idx int, src string) *string {
			return &src
		})
}

func NewService(client *sms.Client, appId string, signName string) *Service {
	return &Service{
		client:   client,
		appId:    &appId,
		signName: &signName,
	}
}
