package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type GzipCompress struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w GzipCompress) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			_, err = io.WriteString(w, err.Error())
			return
		}
		defer func() { err = gz.Close() }()

		w.Header().Set("Content-Encoding", "gzip")
		next.ServeHTTP(GzipCompress{ResponseWriter: w, Writer: gz}, r)
	})
}
