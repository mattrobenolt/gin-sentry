package recovery

import (
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
)

func Recovery(onlyCrashes bool) gin.HandlerFunc {
	return RecoveryWithClient(raven.DefaultClient, onlyCrashes)
}

func RecoveryWithClient(client *raven.Client, onlyCrashes bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rval := recover(); rval != nil {
				debug.PrintStack()
				rvalStr := fmt.Sprint(rval)
				packet := raven.NewPacket(rvalStr, raven.NewException(errors.New(rvalStr), raven.NewStacktrace(2, 3, nil)), raven.NewHttp(c.Request))
				client.Capture(packet, nil)
				c.Writer.WriteHeader(http.StatusInternalServerError)
			}
			if !onlyCrashes {
				for _, item := range c.Errors {
					msg := item.Err.Error()
					packet := raven.NewPacket(msg, &raven.Message{msg, []interface{}{item.Meta}})
					client.Capture(packet, nil)
				}
			}
		}()
		c.Next()
	}
}
