// Package middleware мидлвар для логирования запросов
package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

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

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

// LogMiddleware мидлвар для логирования запросов
func LogMiddleware(h http.Handler) http.Handler {
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
		logrus.Infof("%s %s (%d) %s %d byte", r.Method, r.RequestURI, responseData.status, duration, responseData.size)
	})
}

func ClientInterceptor(ctx context.Context, method string, req interface{},
	reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	// выполняем действия перед вызовом метода
	start := time.Now()
	// вызываем RPC-метод
	err := invoker(ctx, method, req, reply, cc, opts...)
	// выполняем действия после вызова метода
	if err != nil {
		logrus.Printf("[gRPC ERROR] %s, %v", method, err)
	} else {
		logrus.Printf("[gRPC INFO] %s, %v", method, time.Since(start))
	}
	return err
}

func UnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	logrus.Printf("[gRPC INFO] %s", info.FullMethod)
	return handler(ctx, req)
}
