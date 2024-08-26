package prom

import (
	"github.com/dapr-platform/common"
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

var liveGauge prometheus.Gauge

func init() {
	register()
	go loopWriteLive()
}

func register() {
	common.Logger.Debugln("register live gauge")
	name := "live"
	liveGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "iot",
		Subsystem: "service",
		Name:      name,
		Help:      "",
	})

	prometheus.MustRegister(liveGauge)
	common.Logger.Debugln("register live gauge,send counter success")
}

func loopWriteLive() {
	for {

		time.Sleep(time.Second * 15)
		liveGauge.Set(1)
	}
}
