package syslog

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	messageCounterDebug = promauto.NewCounter(prometheus.CounterOpts{
		Name:        "prusa_logged_messages_total",
		Help:        "The total number of processed events",
		ConstLabels: prometheus.Labels{"severity": "debug"},
	})
	messageCounterWarn = promauto.NewCounter(prometheus.CounterOpts{
		Name:        "prusa_logged_messages_total",
		Help:        "The total number of processed events",
		ConstLabels: prometheus.Labels{"severity": "warning"},
	})
	messageCounterError = promauto.NewCounter(prometheus.CounterOpts{
		Name:        "prusa_logged_messages_total",
		Help:        "The total number of processed events",
		ConstLabels: prometheus.Labels{"severity": "error"},
	})
	messageCounterCrit = promauto.NewCounter(prometheus.CounterOpts{
		Name:        "prusa_logged_messages_total",
		Help:        "The total number of processed events",
		ConstLabels: prometheus.Labels{"severity": "critical"},
	})
	messageCounterInfo = promauto.NewCounter(prometheus.CounterOpts{
		Name:        "prusa_logged_messages_total",
		Help:        "The total number of processed events",
		ConstLabels: prometheus.Labels{"severity": "informational"},
	})
	messageCounterUnkn = promauto.NewCounter(prometheus.CounterOpts{
		Name:        "prusa_logged_messages_total",
		Help:        "The total number of processed events",
		ConstLabels: prometheus.Labels{"severity": "unknown"},
	})
)
