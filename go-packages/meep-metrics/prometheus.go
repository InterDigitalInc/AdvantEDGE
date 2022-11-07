/*
 * Copyright (c) 2022  The AdvantEDGE Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package metrics

import (
	"net/http"
	"strconv"
	"time"

	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics
var (
	metricHttpRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "metrics_http_request_duration",
		Help:    "A histogram of http request durations",
		Buckets: prometheus.LinearBuckets(10, 10, 5),
	}, []string{"sbox", "svc", "path", "method", "status"})

	metricHttpNotificationDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "metrics_http_notification_duration",
		Help:    "A histogram of http notification durations",
		Buckets: prometheus.LinearBuckets(10, 10, 5),
	}, []string{"sbox", "svc", "notif", "url", "method", "status"})
)

// ResponseWriter wrapper to capture status code
type MetricsResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewMetricsResponseWriter(w http.ResponseWriter) *MetricsResponseWriter {
	return &MetricsResponseWriter{w, http.StatusOK}
}

func (mrw *MetricsResponseWriter) WriteHeader(code int) {
	mrw.statusCode = code
	mrw.ResponseWriter.WriteHeader(code)
}

// MetricsHandler - REST API Metrics Handler
func MetricsHandler(inner http.Handler, sbox string, svc string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Cache request data
		start := time.Now()
		route := mux.CurrentRoute(r)
		path, _ := route.GetPathTemplate()
		method := r.Method

		// Create metrics response writer wrapper to capture status code
		mrw := NewMetricsResponseWriter(w)

		// Process request
		inner.ServeHTTP(mrw, r)

		// Get request duration & response status code
		procTime := float64(time.Since(start).Microseconds()) / 1000.0
		status := strconv.Itoa(mrw.statusCode)

		log.Info("sbox: ", sbox)
		log.Info("svc: ", svc)
		log.Info("path: ", path)
		log.Info("method: ", method)
		log.Info("status: ", status)
		log.Info("procTime: ", procTime)

		// Store HTTP Request metrics
		metricHttpRequestDuration.WithLabelValues(sbox, svc, path, method, status).Observe(procTime)
	})
}

func ObserveNotification(sbox string, svc string, notif string, url string, resp *http.Response, duration float64) {
	var status string
	if resp != nil {
		status = strconv.Itoa(resp.StatusCode)
	} else {
		status = strconv.Itoa(http.StatusInternalServerError)
	}

	// Store HTTP Notification metrics
	metricHttpNotificationDuration.WithLabelValues(sbox, svc, notif, url, "POST", status).Observe(duration)
}
