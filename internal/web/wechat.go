package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	uuid "github.com/lithammer/shortuuid/v4"
	"github.com/webook/internal/service"
	"github.com/webook/internal/service/oauth2/wechat"
	ijwt "github.com/webook/internal/web/jwt"
	"net/http"
)

type OAuth2WechatHandler struct {
	svc     wechat.Service
	userSvc service.UserService
	ijwt.Handler
	key             []byte
	stateCookieName string
}

func NewOAuth2WechatHandler(svc wechat.Service, hdl ijwt.Handler, userSvc service.UserService) *OAuth2WechatHandler {
	return &OAuth2WechatHandler{
		svc:             svc,
		userSvc:         userSvc,
		key:             []byte("k6CswdUm77WKcbM68UQUuxVsHSpTCwgB"),
		stateCookieName: "jwt-state",
		Handler:         hdl,
	}
}

func (o *OAuth2WechatHandler) RegisterRoutes(server *gin.Engine) {
	g := server.Group("/oahth2/wechat")
	g.GET("/authurl", o.Auth2URL)
	g.Any("/callback", o.Callback)
}

func (o *OAuth2WechatHandler) Auth2URL(ctx *gin.Context) {
	state := uuid.New()
	val, err := o.svc.AuthURL(ctx, state)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Msg:  "Failed to construct redirect url",
			Code: 5,
		})
		return
	}
	err = o.setStateCookie(ctx, state)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Msg:  "Server Error",
			Code: 5,
		})
	}
	ctx.JSON(http.StatusOK, Result{
		Data: val,
	})
}

func (o *OAuth2WechatHandler) Callback(ctx *gin.Context) {
	err := o.verifyState(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Msg:  "illegal request",
			Code: 4,
		})
		return
	}
	code := ctx.Query("code")
	// state := ctx.Query("state")
	wechatInfo, err := o.svc.VerifyCode(ctx, code)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Msg:  "Wrong Authorization Code",
			Code: 4,
		})
		return
	}
	u, err := o.userSvc.FindOrCreateByWechat(ctx, wechatInfo)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Msg:  "System Error",
			Code: 5,
		})
		return
	}
	err = o.SetLoginToken(ctx, u.Id)
	if err != nil {
		ctx.String(http.StatusOK, "system error")
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "OK",
	})
	return
}

func (o *OAuth2WechatHandler) verifyState(ctx *gin.Context) error {
	state := ctx.Query("state")
	ck, err := ctx.Cookie(o.stateCookieName)
	if err != nil {
		return fmt.Errorf("Not able to obtain cookie %w", err)
	}
	var sc StateClaims
	_, err = jwt.ParseWithClaims(ck, &sc, func(token *jwt.Token) (interface{}, error) {
		return o.key, nil
	})
	if err != nil {
		return fmt.Errorf("Failed to decrypt token %w", err)
	}
	if state != sc.State {
		// state not matching
		return fmt.Errorf("state not matching")
	}
	return nil
}

func (o *OAuth2WechatHandler) setStateCookie(ctx *gin.Context, state string) error {
	claims := StateClaims{
		State: state,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString(o.key)
	if err != nil {
		return err
	}
	ctx.SetCookie(o.stateCookieName, tokenStr, 600, "/oauth2/wechat/callback", "", false, true)
	return nil
}

type StateClaims struct {
	jwt.RegisteredClaims
	State string
}
