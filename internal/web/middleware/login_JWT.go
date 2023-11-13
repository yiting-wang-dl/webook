package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	ijwt "github.com/webook/internal/web/jwt"
	"net/http"
)

type LoginJWTMiddlewareBuilder struct {
	ijwt.Handler
}

func NewLoginJWTMiddlewareBuilder(hdl ijwt.Handler) *LoginJWTMiddlewareBuilder {
	return &LoginJWTMiddlewareBuilder{
		Handler: hdl,
	}
}

func (m *LoginJWTMiddlewareBuilder) CheckLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		if path == "/users/signup" ||
			path == "/users/login" ||
			path == "/users/login_sms/code/send" ||
			path == "/users/login_sms" ||
			path == "/oauth2/wechat/authurl"||
			path == "/oauth2/wechat/callback" {
			// No need to check validation
			return
		}

		tokenStr := m.ExtractToken(ctx)
		var uc ijwt.UserClaims
		token, err := jwt.ParseWithClaims(tokenStr, &uc, func(token *jwt.Token) (interface{}, error) {
			return ijwt.JWTKey, nil
		}
		if err != nil {
			//println("Incorrect token")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if token == nil || !token.Valid {
			//println("Token is no longer valid")
			// generate a new access_token?

			// token is illegal
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		err = m.CheckSession(ctx, uc.Ssid)
		if err != nil {
			// invaid token or redis issue
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// check redis issue
		//if cnt > 0 {
		//	// invalid token or redis issue
		//	ctx.AbortWithStatus(http.StatusUnauthorized)
		//	return
		//}

		ctx.Set("user", uc)
	}
}
