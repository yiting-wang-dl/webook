package ratelimit

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/webook/pkg/limiter"
	"log"
	"net/http"
)

type Builder struct {
	prefix  string
	limiter limiter.Limiter
}

func NewBuilder(l limiter.Limiter) *Builder {
	return &Builder{
		prefix:  "ip-limiter",
		limiter: l,
	}
}

func (b *Builder) Prefix(prefix string) *Builder {
	b.prefix = prefix
	return b
}

func (b *Builder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		if ctx.GetHeader("x-stress") == "true" {
			// use context.Context to save this
			newCtx := context.WithValue(ctx, "x-stress", true)
			ctx.Request = ctx.Request.Clone(newCtx)
			ctx.Next()
			return
		}

		limited, err := b.limiter.Limit(ctx, fmt.Sprintf("%s:%s", b.prefix, ctx.ClientIP()))
		if err != nil {
			log.Println(err)
			// Here we implement a logic that if redis crashed, in order to save the system, limit immediately
			ctx.AbortWithStatus(http.StatusInternalServerError)
			// another way, if redis crashed, still serve the normal users, do not limit rate
			// ctx.Next()
			return
		}
		if limited {
			log.Println(err)
			ctx.AbortWithStatus(http.StatusTooManyRequests)
			return
		}
		ctx.Next()
	}
}
