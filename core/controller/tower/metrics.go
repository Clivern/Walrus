// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package tower

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
)

var (
	workersCount = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "walrus",
			Name:      "workers_count",
			Help:      "Number of Async Workers",
		})

	queueCapacity = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "walrus",
			Name:      "workers_queue_capacity",
			Help:      "The maximum number of messages queue can process",
		})
)

func init() {
	prometheus.MustRegister(workersCount)
	prometheus.MustRegister(queueCapacity)
}

// Metrics controller
func Metrics() http.Handler {
	workersCount.Set(float64(viper.GetInt(fmt.Sprintf(
		"%s.broker.native.workers",
		viper.GetString("role"),
	))))

	queueCapacity.Set(float64(viper.GetInt(fmt.Sprintf(
		"%s.broker.native.capacity",
		viper.GetString("role"),
	))))

	return promhttp.Handler()
}
