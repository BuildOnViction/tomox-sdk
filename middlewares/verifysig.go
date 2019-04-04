package middlewares

import (
	"net/http"
)

func VerifySignature(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// TODO: Logic to verify the signature of user here

		next.ServeHTTP(w, r)
	})
}
