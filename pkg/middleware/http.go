package middleware

import (
	"bytes"
	"net/http"

	"go.uber.org/zap"
)

func newResponseWriterLogger(w http.ResponseWriter) *responseWriterLogger {
	return &responseWriterLogger{w: w}
}

type responseWriterLogger struct {
	w    http.ResponseWriter
	body bytes.Buffer
	code int
}

func (w *responseWriterLogger) Write(p []byte) (int, error) {
	w.body.Write(p)
	return w.w.Write(p)
}

func (w *responseWriterLogger) WriteHeader(code int) {
	w.code = code
	w.w.WriteHeader(code)
}

func (w *responseWriterLogger) Header() http.Header {
	return w.w.Header()
}

func NewHTTPLogger(log *zap.SugaredLogger) HTTPMiddleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			writerWrapper := newResponseWriterLogger(w)
			next.ServeHTTP(writerWrapper, r)
			log.Debugf("[%v] %v %v %v", r.Method, r.URL.String(), writerWrapper.code, writerWrapper.body.String())
		})
	}
}
