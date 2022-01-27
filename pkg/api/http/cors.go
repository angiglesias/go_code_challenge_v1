package http

import (
	"net/http"

	"challenge/pkg/logging"
)

// Very simple, free pass cors implementation to accept preflight requests and third party domains API comsumption on browsers
func SetupCorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(rw http.ResponseWriter, r *http.Request) {
			// respond preflight OPTIONS request
			if r.Method == http.MethodOptions {
				logging.Debugf("[CORS] Processing preflight request")
				// write CORS headers to response
				rw.Header().Set("Access-Control-Allow-Origin", "*")                                                                                   // wildcard origin
				rw.Header().Set("Access-Control-Allow-Credentials", "true")                                                                           //
				rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")                                                    // allow most common http methods
				rw.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization") // common headers used in REST API requests
				rw.Header().Set("Access-Control-Max-Age", "3600")                                                                                     // very large time, depends in real production requirements
				return
			}
			// send request to next http handler
			next.ServeHTTP(rw, r)
		},
	)
}
