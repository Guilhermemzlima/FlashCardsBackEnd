package middleware

import (
	"errors"
	"github.com/Guilhermemzlima/FlashCardsBackEnd/pkg/api/render"
	"net/http"
)

func Header(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		if userId := r.Header.Get("userId"); userId == "" {
			render.ResponseError(w, errors.New("userId is required as a header parameter"), http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
