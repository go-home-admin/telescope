package telescope

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"runtime/debug"
	"time"
)

// Telescope 启用望眼镜
func Telescope() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				log.WithField("url", ctx.Request.URL.Path).Error(fmt.Sprint(r) + "\n" + string(debug.Stack()))
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"code":    1,
					"message": "服务器内部错误",
				})
				ctx.Abort()
			}
			if errorRecord && hasError {
				log.WithContext(ctx).WithFields(log.Fields{"type": "request"}).Error(ctx.Request.URL)
				defer func() {
					hasError = false
				}()
			} else {
				log.WithContext(ctx).WithFields(log.Fields{"type": "request"}).Debug(ctx.Request.URL)
			}
			TelescopeClose()
		}()
		ctx.Set("start", time.Now())
		TelescopeStart()
		ctx.Writer = &TelescopeResponseWriter{
			Body:           bytes.NewBufferString(""),
			ResponseWriter: ctx.Writer,
		}

		data, err := ctx.GetRawData()
		if err == nil {
			ctx.Set("raw", data)
			ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(data))
		}
		ctx.Next()
	}
}
