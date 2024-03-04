package middleware

import (
	"encoding/gob"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type LoginMiddlewareBuilder struct {
	paths []string
}

func NewLoginMiddlewareBuilder() *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{}
}

func (l *LoginMiddlewareBuilder) IgnorePaths(path string) *LoginMiddlewareBuilder {
	l.paths = append(l.paths, path)
	return l
}

func (l *LoginMiddlewareBuilder) Build() gin.HandlerFunc {
	gob.Register(time.Now())
	return func(ctx *gin.Context) {
		// 登录和注册不需要登录校验
		for _, path := range l.paths {
			if ctx.Request.URL.Path == path {
				return
			}
		}
		// 这一步相当于是从ctx中拿到DefaultKey对应的值，也就是session
		sess := sessions.Default(ctx)
		id := sess.Get("userID")
		// 如果没有登录
		if id == nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 刷新session
		updateTime := sess.Get("update_time")
		now := time.Now()
		if updateTime == nil {
			sess.Set("update_time", now)
			sess.Set("userId", id)
			sess.Options(sessions.Options{
				MaxAge: 60,
			})
			// 因为main里面已经确定了store是redis，所以这里实际上是把session中的数据保存到redis中去
			if err := sess.Save(); err != nil {
				panic(err)
			}
		} else {
			updateTimeVal, _ := updateTime.(time.Time)
			if now.Sub(updateTimeVal) > time.Second*10 {
				sess.Set("update_time", now)
				sess.Set("userId", id)
				sess.Options(sessions.Options{
					MaxAge: 60,
				})
				if err := sess.Save(); err != nil {
					panic(err)
				}
			}
		}
	}
}
