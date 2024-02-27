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
