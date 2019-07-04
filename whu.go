package main

import (
	"os"
	"strconv"
	"time"

	"github.com/bloom42/rz-go/v2"
	"github.com/bloom42/rz-go/v2/log"

	"github.com/asdine/storm"
	"github.com/labstack/echo"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/zhuharev/talert"
	"github.com/zhuharev/whu/domain/update"
	"github.com/zhuharev/whu/domain/update/delivery"
	"github.com/zhuharev/whu/domain/update/repo"
	whStormDB "github.com/zhuharev/whu/domain/webhook/repo/storm"
)

const version = "0.0.11"

func main() {
	log.Info("start whu", rz.String("version", version))
	if dsn := os.Getenv("TALERT_DSN"); dsn != "" {
		token, chatID, err := talert.ParseDSN(dsn)
		if err != nil {
			panic(err)
		}
		talert.Init(token, chatID)
	}
	talert.Alert("whu started", talert.String("version", version))
	sdb, err := storm.Open(os.Getenv("DB_PATH"))
	if err != nil {
		panic(err)
	}
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
	if err := e.Start(":" + os.Getenv("PORT")); err != nil {
		panic(err)
	}
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
