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
	red    = "\033[31m"
	green  = "\033[32m"
	yellow = "\033[33m"
	cyan   = "\033[36m"
	reset  = "\033[0m"
)

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

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

func WithLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		statusRecorder := &statusRecorder{
			ResponseWriter: w,
			status:         http.StatusOK,
		}

		next.ServeHTTP(statusRecorder, r)

		log.Printf(
			"%s%s%s %s %s%d%s %v",
			cyan,
			strconv.Quote(r.Method),
			reset,
			strconv.Quote(r.URL.Path),
			statusColor(statusRecorder.status),
			statusRecorder.status,
			reset,
			time.Since(start),
		)
	})
}
