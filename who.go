package main

import (
	"os"
	"strconv"
	"time"

	"github.com/asdine/storm"
	"github.com/labstack/echo"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/zhuharev/whu/domain/update"
	"github.com/zhuharev/whu/domain/update/delivery"
	"github.com/zhuharev/whu/domain/update/repo"
	whStormDB "github.com/zhuharev/whu/domain/webhook/repo/storm"
)

func main() {
	sdb, err := storm.Open(os.Getenv("DB_PATH"))
	defer sdb.Close()
	repo, err := repo.New(sdb)
	if err != nil {
		panic(err)
	}
	whRepo := whStormDB.New(sdb)
	updUC := update.New(whRepo, repo)
	e := echo.New()
	e.Use(HTTPMetrics("who"))
	// metrics
	e.Any("/metrics", echo.WrapHandler(promhttp.Handler()))
	delivery.New(e, updUC, whRepo)
	e.Start(":" + os.Getenv("PORT"))
}

// HTTPMetrics is the middleware function that logs duration of responses.
func HTTPMetrics(appName string) echo.MiddlewareFunc {
	labels := []string{"method", "uri", "code"}

	echoRequests := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: appName,
		Subsystem: "http",
		Name:      "requests_count",
		Help:      "Requests count by method/path/status.",
	}, labels)

	echoDurations := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: appName,
		Subsystem: "http",
		Name:      "responses_duration_seconds",
		Help:      "Response time by method/path/status.",
	}, labels)

	prometheus.MustRegister(echoRequests, echoDurations)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			if err := next(c); err != nil {
				c.Error(err)
			}

			metrics := []string{c.Request().Method, c.Path(), strconv.Itoa(c.Response().Status)}

			echoDurations.WithLabelValues(metrics...).Observe(time.Since(start).Seconds())
			echoRequests.WithLabelValues(metrics...).Inc()

			return nil
		}
	}
}
