package middleware

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw
		c.Next()

		statusCode := c.Writer.Status()
		reqUri := c.Request.RequestURI
		internalError := blw.body.String()
		internalError = strings.ReplaceAll(internalError, "\"", "")
		clientIP := c.ClientIP()
		reqMethod := c.Request.Method

		logMsg, _ := json.Marshal(map[string]interface{}{
			"statusCode": statusCode,
			"clientIP":   clientIP,
			"reqMethod":  reqMethod,
			"reqUri":     reqUri,
			"detail":     internalError,
		})

		if statusCode >= 400 {
			log.Error(string(logMsg))
		} else {
			log.Info(string(logMsg))
		}
	}
}
