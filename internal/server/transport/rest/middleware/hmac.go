package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"io"
	"net/http"
)

func HMACMiddleware(next http.HandlerFunc, secretKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h := r.Header.Get("Hash")
		if h == "none" || h == "" {
			next.ServeHTTP(w, r)
			return
		}

		body := r.Body
		data, _ := io.ReadAll(body)

		h1 := hmac.New(sha256.New, []byte(secretKey))
		h1.Write(data)
		res := h1.Sum(nil)
		if !hmac.Equal([]byte(h), res) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	}
}
