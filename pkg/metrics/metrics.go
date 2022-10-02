package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"net/http"
)

func Start(logger *zap.SugaredLogger) {
	http.Handle("/metrics", promhttp.Handler())
	logger.Fatal(http.ListenAndServe(":9090", nil))
}

var Pulls = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "pull_total",
		Help: "Pulls from git",
	},
	[]string{"repo"},
)

var Clones = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "clone_total",
		Help: "Clones from git",
	},
	[]string{"repo"},
)
