// Package middleware мидлвар для проверки ip адреса
package middleware

import (
	"net"
	"net/http"

	"github.com/sirupsen/logrus"
)

// TrustedSubnetMiddleware проверка на принятие данных от доверенного ip адреса
func TrustedSubnetMiddleware(cidr string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("X-Real-IP") != "" {
				ip := net.ParseIP(r.Header.Get("X-Real-IP"))
				_, subnetIpNet, err := net.ParseCIDR(cidr)
				if err != nil {
					logrus.Errorf("Ошибка при парсинге CIDR: %v", err)
					return
				}
				if subnetIpNet.Contains(ip) {
					logrus.Printf("%s входит в сеть %s\n", ip, cidr)
					h.ServeHTTP(w, r)
				} else {
					logrus.Printf("%s не входит в сеть %s\n", ip, cidr)
					w.WriteHeader(http.StatusForbidden)
				}
			} else {
				w.WriteHeader(http.StatusForbidden)
			}
		})
	}
}
