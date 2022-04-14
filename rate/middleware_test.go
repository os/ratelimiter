package rate

import (
	"github.com/pkg/errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestLimit(t *testing.T) {
	for _, tc := range []struct {
		Description             string
		LimiterError            error
		LimiterResponse         Response
		ExpectedStatusCode      int
		ExpectedResponseHeaders map[string]string
		ExpectHandlerToBeCalled bool
	}{
		{
			Description:  "Handlers are called if rate limits are not exceeded",
			LimiterError: nil,
			LimiterResponse: Response{
				Limit:     10,
				Duration:  time.Second,
				Remaining: 9,
				Reset:     time.Unix(1607090218, 0),
			},
			ExpectedStatusCode: 200,
			ExpectedResponseHeaders: map[string]string{
				"X-RateLimit-Limit-1s":     "10",
				"X-RateLimit-Remaining-1s": "9",
				"X-RateLimit-Reset-1s":     "1607090218",
			},
			ExpectHandlerToBeCalled: true,
		},
		{
			Description:  "Handlers are not called if a rate limit is exceeded",
			LimiterError: ErrRateLimitExceeded,
			LimiterResponse: Response{
				Limit:     10,
				Duration:  time.Second,
				Remaining: 0,
				Reset:     time.Unix(1607090218, 0),
			},
			ExpectedStatusCode: http.StatusTooManyRequests,
			ExpectedResponseHeaders: map[string]string{
				"X-RateLimit-Limit-1s":     "10",
				"X-RateLimit-Remaining-1s": "0",
				"X-RateLimit-Reset-1s":     "1607090218",
			},
			ExpectHandlerToBeCalled: false,
		},
		{
			Description:  "Handlers are not called if there is an error",
			LimiterError: errors.New("Connection error"),
			LimiterResponse: Response{
				Limit:     10,
				Duration:  time.Second,
				Remaining: 9,
				Reset:     time.Unix(1607090218, 0),
			},
			ExpectedStatusCode:      http.StatusInternalServerError,
			ExpectHandlerToBeCalled: false,
		},
	} {
		t.Run(tc.Description, func(t *testing.T) {
			handlerCalled := false
			identifier := FakeIdentifier{
				FnIdentify: func(r *http.Request) (string, error) {
					return "default", nil
				},
			}
			limiter := FakeLimiter{
				FnAllow: func(id string) (Response, error) {
					return tc.LimiterResponse, tc.LimiterError
				},
			}
			handler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				handlerCalled = true
			})
			w := httptest.NewRecorder()
			r, _ := http.NewRequest("GET", "", nil)
			limit := Limit(identifier, limiter)
			limit(handler).ServeHTTP(w, r)

			if w.Code != tc.ExpectedStatusCode {
				t.Errorf("Expected status code: %d, got: %d", tc.ExpectedStatusCode, w.Code)
			}
			if handlerCalled != tc.ExpectHandlerToBeCalled {
				if tc.ExpectHandlerToBeCalled {
					t.Errorf("Expected handler to be called")
				} else {
					t.Errorf("Expected handler not to be called")
				}
			}
			headers := w.Header()
			for k, v := range tc.ExpectedResponseHeaders {
				hv := headers.Get(k)
				if hv != v {
					t.Errorf("Expected %s header to be set to %s, got: %s", k, v, hv)
				}
			}
		})
	}
}
