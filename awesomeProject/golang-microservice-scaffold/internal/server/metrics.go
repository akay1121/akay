package server

import (
	"example/internal/conf"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/metrics"
	"go.opentelemetry.io/otel"
)

func NewMetricsMiddleware(c *conf.Metrics) middleware.Middleware {
	meter := otel.Meter("example-service")
	metricRequests, _ := metrics.DefaultRequestsCounter(meter, metrics.DefaultServerRequestsCounterName)
	metricSeconds, _ := metrics.DefaultSecondsHistogram(meter, metrics.DefaultServerSecondsHistogramName)
	return metrics.Server(
		metrics.WithSeconds(metricSeconds),
		metrics.WithRequests(metricRequests),
	)
}
