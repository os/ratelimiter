package rate

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/pkg/errors"
)

// setResponseHeaders for the given Response.
func setResponseHeaders(r Response, w http.ResponseWriter) {
	h := w.Header()
	suffix := r.Duration.String()
	h.Set(fmt.Sprintf("X-RateLimit-Limit-%s", suffix),
		strconv.FormatUint(uint64(r.Limit), 10))
	h.Set(fmt.Sprintf("X-RateLimit-Remaining-%s", suffix),
		strconv.FormatUint(uint64(r.Remaining), 10))
	h.Set(fmt.Sprintf("X-RateLimit-Reset-%s", suffix),
		strconv.FormatInt(r.Reset.Unix(), 10))
}

// Limit middleware prevents the requests being executed if they exceed the
// constraints of the given limiter. Limiter constraints are applied based on
// the given identifier. Therefore the identifier can be used to apply a certain
// limit for a set of consumers as well as individual consumers.
func Limit(identifier Identifier, limiters ...Limiter) func(http.HandlerFunc) http.HandlerFunc {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			id, err := identifier.Identify(r)
			if err != nil {
				log.Println(errors.Wrap(err,
					"Error occurred during rate limiting"))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			for _, l := range limiters {
				response, err := l.Allow(id)
				if err != nil {
					if err == ErrRateLimitExceeded {
						setResponseHeaders(response, w)
						w.WriteHeader(http.StatusTooManyRequests)
					} else {
						log.Println(errors.Wrap(err,
							"Error occurred during limiting"))
						w.WriteHeader(http.StatusInternalServerError)
					}
					return
				}
				setResponseHeaders(response, w)
			}

			h.ServeHTTP(w, r)
		}
	}
}
