package server

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	casbinEnforceAll = promauto.NewCounter(prometheus.CounterOpts{
		Name: "casbin_enforce_all",
		Help: "The total number of enforce requests received",
	})
	casbinEnforceAccessApproved = promauto.NewCounter(prometheus.CounterOpts{
		Name: "casbin_enforce_approved",
		Help: "The total number of enforce requests that returned with true",
	})
	casbinEnforceAccessDenied = promauto.NewCounter(prometheus.CounterOpts{
		Name: "casbin_enforce_denied",
		Help: "The total number of enforce requests that returned with false",
	})
)
