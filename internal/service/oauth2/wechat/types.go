package wechat

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/webook/internal/domain"
	"github.com/webook/pkg/logger"
	"net/http"
	"net/url"
)

type Service interface {
	AuthURL(ctx context.Context, state string) (string, error)
	VerifyCode(ctx context.Context, code string) (domain.WechatInfo, error)
}

var redirectURL = url.PathEscape("https://mydomain.com/oauth2/wechat/callback")

type service struct {
	appID     string
	appSecret string
	client    *http.Client
	l         logger.LoggerV1
}

func NewService(appID string, appSecret string, l logger.LoggerV1) Service {
	return &service{
		appID:     appID,
		appSecret: appSecret,
		client:    http.DefaultClient,
	}
}
func (s *service) VerifyCode(ctx context.Context,
	code string) (domain.WechatInfo, error) {
	accessTokenUrl := fmt.Sprintf(`https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code`,
		s.appID, s.appSecret, code)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, accessTokenUrl, nil)
	if err != nil {
		return domain.WechatInfo{}, err
	}
	httpResp, err := s.client.Do(req)
	if err != nil {
		return domain.WechatInfo{}, err
	}

	var res Result
	err = json.NewDecoder(httpResp.Body).Decode(&res)
	if err != nil {
		// Failed to transfer JSON to struct
		return domain.WechatInfo{}, err
	}
	if res.ErrCode != 0 {
		return domain.WechatInfo{},
			fmt.Errorf("failed to call wechat interface errcode %d, errmsg %s", res.ErrCode, res.ErrMsg)
	}
	return domain.WechatInfo{
		UnionId: res.UnionId,
		OpenId:  res.OpenId,
	}, nil
}

func (s *service) AuthURL(ctx context.Context, state string) (string, error) {
	const authURLPattern = `https://open.weixin.qq.com/connect/qrconnect?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_login&state=%s#wechat_redirect`
	return fmt.Sprintf(authURLPattern, s.appID, redirectURL, state), nil
}

type Result struct {
	AccessToken string `json:"access_token"`
	// access_token interface calling certificate timed out，unit（second）
	ExpiresIn int64 `json:"expires_in"`
	// user refreh access_token
	RefreshToken string `json:"refresh_token"`
	// authorize user unique id
	OpenId string `json:"openid"`
	// user authorized domain，seperated by (,)
	Scope string `json:"scope"`
	// only when this website is authorized by userinfo
	UnionId string `json:"unionid"`

	// error results
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}
