package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"io"
	"net/http"
)

func HMACMiddleware(next http.HandlerFunc, secretKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h := r.Header.Get("HashSHA256")
		if h == "none" || h == "" {
			next.ServeHTTP(w, r)
			return
		}

		body := r.Body
		data, err := io.ReadAll(body)
		r.Body = io.NopCloser(bytes.NewBuffer(data))
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		h1 := hmac.New(sha256.New, []byte(secretKey))
		_, err = h1.Write(data)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		res := h1.Sum(nil)
		if !hmac.Equal([]byte(h), res) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	}
}
