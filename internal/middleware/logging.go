package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

type statusRecorder struct {
	http.ResponseWriter

	status int
}

const (
	magenta = "\u001b[35m"
	red    = "\033[31m"
	green  = "\033[32m"
	yellow = "\033[33m"
	cyan   = "\033[36m"
	reset  = "\033[0m"
)

// WriteHeader sends a HTTP header response and saves the status to r.status
func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

// Write writes b as an HTTP reply
func (r *statusRecorder) Write(b []byte) (int, error) {
	if r.status == 0 {
		r.status = http.StatusOK
	}

	n, err := r.ResponseWriter.Write(b)
	if err != nil {
		return n, fmt.Errorf("failed to write response: %w", err)
	}

	return n, nil
}

// statusColor determines the statusColor based on the status code
func statusColor(status int) string {
	switch {
	case status >= http.StatusInternalServerError:
		return red
	case status >= http.StatusBadRequest:
		return yellow
	case status >= http.StatusMultipleChoices:
		return cyan
	case status >= http.StatusOK:
		return green
	default:
		return reset
	}
}

// Logging Middleware, Logs Request Method, Method URL Path, URL and time taken to satisfy request
func WithLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		statusRecorder := &statusRecorder{
			ResponseWriter: w,
			status:         http.StatusOK,
		}

		next.ServeHTTP(statusRecorder, r)

		log.Printf(
			"%s%s%s %s%s%s %s %s%d%s %v",
			cyan,
			strconv.Quote(r.Method),
			reset,
			strconv.Quote(r.URL.Path),
			magenta,
			strconv.Quote(r.URL.String()),
			reset,
			statusColor(statusRecorder.status),
			statusRecorder.status,
			reset,
			time.Since(start),
		)
	})
}
