package middleware

import (
	"github.com/evrone/go-clean-template/internal/user/metrics"
	"github.com/gin-gonic/gin"
	"time"
)

func MetricsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		rw := newResponseWriter(c.Writer)

		c.Next()

		path := c.Request.URL.Path
		statusString := rw.GetStatusString()

		metrics.HttpResponseTime.WithLabelValues(path, statusString, c.Request.Method).Observe(time.Since(start).Seconds())
		metrics.HttpRequestsTotalCollector.WithLabelValues(path, statusString, c.Request.Method).Inc()
	}
}
