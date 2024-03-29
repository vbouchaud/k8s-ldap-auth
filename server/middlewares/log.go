package middlewares

import (
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

// ProxyResponseWriter is a workaround for getting HTTP Response information in the access logs
// With the default ResponseWriter interface we only have methods to write the status code and parts of
// the response content. But there's no properties to tell the current status of the response.
// This implementation proxy the default ResponseWriter methods to the usual one, but keeps the
// current state of the request in it's properties
type ProxyResponseWriter struct {
	code   int
	length int
	parent http.ResponseWriter
}

// NewProxyResponseWriter create a new ProxyResponseWriter that wrap API calls to another ResponseWriter
func NewProxyResponseWriter(parent http.ResponseWriter) *ProxyResponseWriter {
	return &ProxyResponseWriter{
		code:   200,
		length: 0,
		parent: parent,
	}
}

// Header return the inner ResponseWriter Header
func (brs *ProxyResponseWriter) Header() http.Header {
	return brs.parent.Header()
}

// Write a portion of the response content to the inner ResponseWriter, and keep track of the byte length added
func (brs *ProxyResponseWriter) Write(content []byte) (int, error) {
	length, err := brs.parent.Write(content)
	brs.length += length
	return length, err
}

// WriteHeader to the inner ResponseWriter, and keep track of the current response status
func (brs *ProxyResponseWriter) WriteHeader(code int) {
	brs.code = code
	brs.parent.WriteHeader(code)
}

// AccessLog provide an HTTP server middleware for logging access to the server
// It follows tha common Apache Access Log format, except for the %l and %u values
// that are not implemented yet (%l should probably be ignore anymay).
// You can find more information about this format here : https://httpd.apache.org/docs/2.4/logs.html
func AccessLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		received := time.Now()
		wrapper := NewProxyResponseWriter(res)
		next.ServeHTTP(wrapper, req)
		elapsed := time.Now().Sub(received)

		log.Info().
			Str("remoteaddr", req.RemoteAddr).
			Str("method", req.Method).
			Str("url", req.URL.String()).
			Str("proto", req.Proto).
			Int("code", wrapper.code).
			Int("length", wrapper.length).
			Str("referer", req.Header.Get("Referer")).
			Str("useragent", req.Header.Get("User-Agent")).
			Int64("elapsed", elapsed.Microseconds()).
			Send()
	})
}
