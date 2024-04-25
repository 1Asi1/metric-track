package middleware

import (
	"net"
	"net/http"
	"strings"
)

func CheckSubnetMiddleware(next http.HandlerFunc, trustedSubnet string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if trustedSubnet != "" {
			agentIP := net.ParseIP(strings.TrimSpace(r.Header.Get("X-Real-IP")))
			_, subnet, err := net.ParseCIDR(trustedSubnet)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			if !subnet.Contains(agentIP) {
				w.WriteHeader(http.StatusForbidden)
				return
			}
		}

		next.ServeHTTP(w, r)
	}
}
