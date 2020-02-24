package metrics

import (
	"net/http"
	"reflect"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

type Metrics struct {
	logger         *zap.Logger
	metricsHandler http.Handler
}

type ProcessStartFunc func()
type ProcessFinishFunc func(error)

type loggerError interface {
	Error(args ...interface{})
}

type errorLoggger struct {
	logger loggerError
}

func (el errorLoggger) Println(args ...interface{}) {
	el.logger.Error(args...)
}

// check access and allow if ok
func (m *Metrics) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.metricsHandler.ServeHTTP(w, r)
}

func New(logger *zap.Logger, gatherer prometheus.Gatherer) *Metrics {
	metrics := Metrics{
		logger: logger,
		metricsHandler: promhttp.HandlerFor(gatherer, promhttp.HandlerOpts{
			ErrorLog: errorLoggger{logger: logger.Sugar()},
		}),
	}
	return &metrics
}

func ToefCounters(r prometheus.Registerer, name, description string) (start ProcessStartFunc, finish ProcessFinishFunc) {
	if r == nil || reflect.ValueOf(r).IsNil() {
		return func() {}, func(error) {}
	}

	total := prometheus.NewCounter(prometheus.CounterOpts{Name: name, Help: "Total count of " + description})
	ok := prometheus.NewCounter(prometheus.CounterOpts{Name: name + "_ok", Help: "Ok count of " + description})
	err := prometheus.NewCounter(prometheus.CounterOpts{Name: name + "_err", Help: "Err count of " + description})
	inflight := prometheus.NewGauge(prometheus.GaugeOpts{Name: name + "_inflight", Help: "Err count of " + description})

	r.MustRegister(total, ok, err, inflight)

	start = func() {
		total.Inc()
		inflight.Inc()
	}
	finish = func(error error) {
		if error == nil {
			ok.Inc()
		} else {
			err.Inc()
		}
		inflight.Dec()
	}
	return start, finish
}
