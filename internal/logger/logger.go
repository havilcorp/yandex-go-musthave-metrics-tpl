package logger

import (
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// type ResponseWriter interface {
// 	Header() Header
// 	Write([]byte) (int, error)
// 	WriteHeader(statusCode int)
// }

type (
	responseData struct {
		status int
		size   int
	}
)

type (
	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

// type ResponseWriter interface {
// 	Header() http.Header
// 	Write([]byte) (int, error)
// 	WriteHeader(statusCode int)
// }

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	fmt.Println("Write")
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	fmt.Println("WriteHeader")
	fmt.Println("status", statusCode)
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

func WithLogging(h http.Handler) http.Handler {

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	sugar := *logger.Sugar()
	sugar.Infoln()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		responseData := &responseData{
			status: 200,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}
		h.ServeHTTP(&lw, r)

		duration := time.Since(start)

		sugar.Infoln(
			"uri", r.RequestURI,
			"method", r.Method,
			"status", responseData.status, // получаем перехваченный код статуса ответа
			"duration", duration,
			"size", responseData.size, // получаем перехваченный размер ответа
		)
	})
}
