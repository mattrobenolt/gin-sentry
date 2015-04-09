package recovery

import (
	"errors"
	"fmt"
	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"runtime/debug"
)

func Recovery(client *raven.Client, onlyCrashes bool) gin.HandlerFunc {
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
					packet := raven.NewPacket(item.Err, &raven.Message{item.Err, []interface{}{item.Meta}})
					client.Capture(packet, nil)
				}
			}
		}()
		c.Next()
	}
}
