package logger

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
)

// RequestLogger — middleware-логер для входящих HTTP-запросов.
func RequestLogger(next http.Handler) http.Handler {
	// получаем Handler приведением типа http.HandlerFunc
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timeStart := time.Now()

		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		defer func() {

			Log.Info("got incoming HTTP request",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Int("status", ww.Status()),
				zap.Int("length", ww.BytesWritten()),
				zap.Duration("time", time.Since(timeStart)),
			)
		}()

		next.ServeHTTP(ww, r)
	})
}
