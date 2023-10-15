package middleware

import (
	"encoding/gob"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type LoginMiddlewareBuilder struct {
}

func (m *LoginMiddlewareBuilder) CheckLogin() gin.HandlerFunc {
	gob.Register(time.Now())
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		if path == "/users/signup" || path == "/users/login" {
			// no need to valid session
			return
		}

		sess := sessions.Default(ctx)
		userId := sess.Get("userId")
		if userId == nil {
			//abort
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		now := time.Now()

		const updateTimeKey = "update_time"
		val := sess.Get(updateTimeKey)        // get the previous updateTimeKey
		lastUpdateTime, ok := val.(time.Time) // assert last update TimeKey is a time
		if val == nil || !ok || now.Sub(lastUpdateTime) > time.Second*60 {
			sess.Set(updateTimeKey, now)
			sess.Set("userId", userId)
			err := sess.Save()
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
