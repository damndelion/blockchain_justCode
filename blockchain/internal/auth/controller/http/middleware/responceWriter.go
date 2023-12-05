package middleware

import (
	"net/http"
	"strconv"
)

type responseWriter struct {
	http.ResponseWriter
	status int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{
		ResponseWriter: w,
		status:         http.StatusOK,
	}
}

func (r *responseWriter) GetStatusString() string {
	status := r.status
	if r.status == 0 {
		status = http.StatusOK
	}

	return strconv.Itoa(status)
}

// WriteHeader implements http.ResponseWriter and saves status.
func (r *responseWriter) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}
