package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/webook/internal/web"
	"log"
	"net/http"
	"strings"
	"time"
)

type LoginJWTMiddlewareBuilder struct {
}

func (m *LoginJWTMiddlewareBuilder) CheckLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		if path == "/users/signup" || path == "/users/login" {
			return
		}
		// token will be in the Authorization Header
		// Bearer XXXX
		authCode := ctx.GetHeader("Authorization")
		if authCode == "" {
			// no token if not logged in. no Authorization Header
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		authSegments := strings.Split(authCode, " ")
		if len(authSegments) != 2 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenStr := authSegments[1]
		println("toeknStr in login_JWT: ", tokenStr)
		var uc web.UserClaims
		token, err := jwt.ParseWithClaims(tokenStr, &uc, func(token *jwt.Token) (interface{}, error) {
			return web.JWTKey, nil
		})
		if err != nil {
			println("Incorrect token")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if token == nil || !token.Valid {
			println("Token is no longer valid")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if uc.UserAgent != ctx.GetHeader("User-Agent") {
			// Pay Attention! This might be an attacker
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		expireTime := uc.ExpiresAt
		// if check every 1 min, need to refresh when remaining time is < 50s
		if expireTime.Sub(time.Now()) < time.Second*50 { //time.Second*6000 {
			uc.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute * 1)) // check every 1 minute
			//uc.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Hour * 15)) // testing purpose, set to 15 Hour
			tokenStr, err = token.SignedString(web.JWTKey)
			ctx.Header("x-jwt-token", tokenStr)
			if err != nil {
				log.Println(err)
			}
		}
		ctx.Set("user", uc)
	}
}
