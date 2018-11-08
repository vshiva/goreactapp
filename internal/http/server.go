package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/vshiva/goreactapp/web"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type serverContext struct {
	Mode string
}

// New Server.
func New(bindAddress string, port int, mode string) (*http.Server, error) {
	gin.SetMode(gin.ReleaseMode)

	h := gin.New()
	h.Use(ginlogger(time.RFC3339, true))

	svrContext := &serverContext{Mode: mode}
	api := h.Group("/api")
	api.GET("/config", configData(svrContext))
	h.Use(static.Serve("/", web.FS()))

	return &http.Server{
		Addr:    fmt.Sprintf("%s:%d", bindAddress, port),
		Handler: h,
	}, nil
}

func ginlogger(timeFormat string, utc bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		// some evil middlewares modify this values
		path := c.Request.URL.Path
		c.Next()

		end := time.Now()
		latency := end.Sub(start)
		if utc {
			end = end.UTC()
		}

		entry := log.WithFields(log.Fields{
			"status":     c.Writer.Status(),
			"method":     c.Request.Method,
			"path":       path,
			"ip":         c.ClientIP(),
			"latency":    latency,
			"user-agent": c.Request.UserAgent(),
			"time":       end.Format(timeFormat),
		})

		if len(c.Errors) > 0 {
			// Append error field if this is an erroneous request.
			entry.Error(c.Errors.String())
		} else {
			entry.Info()
		}
	}
}

func configData(srvCtx *serverContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		response := fmt.Sprintf(`{"mode" : "%s"}`, srvCtx.Mode)
		c.Data(200, "application/json", []byte(response))
	}
}
