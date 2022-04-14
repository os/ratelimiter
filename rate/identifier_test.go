package rate

import (
	"net/http/httptest"
	"testing"
)

func TestIPIdentifier_Identify(t *testing.T) {
	for _, tc := range []struct {
		Description   string
		XRealIP       string
		XForwardedFor string
		RemoteAddr    string
		ExpectedIP    string
		ExpectError   bool
	}{
		{
			Description:   "Returns X-REAL-IP when available",
			XRealIP:       "10.0.0.1",
			XForwardedFor: "10.0.0.2, 10.0.0.4, 10.0.0.5",
			RemoteAddr:    "10.0.0.3:8000",
			ExpectedIP:    "10.0.0.1",
			ExpectError:   false,
		},
		{
			Description:   "Returns X-FORWARDED-FOR when X-REAL-IP is not available",
			XRealIP:       "",
			XForwardedFor: "10.0.0.2, 10.0.0.4, 10.0.0.5",
			RemoteAddr:    "10.0.0.3:8000",
			ExpectedIP:    "10.0.0.2",
			ExpectError:   false,
		},
		{
			Description:   "Returns RemoteAddr when X-FORWARDED-FOR is not available",
			XRealIP:       "",
			XForwardedFor: "",
			RemoteAddr:    "10.0.0.3:8000",
			ExpectedIP:    "10.0.0.3",
			ExpectError:   false,
		},
		{
			Description:   "Returns an error when IP address could not be identified",
			XRealIP:       "",
			XForwardedFor: "",
			RemoteAddr:    "",
			ExpectedIP:    "",
			ExpectError:   true,
		},
	} {
		t.Run(tc.Description, func(t *testing.T) {
			identifier := NewIPIdentifier()
			request := httptest.NewRequest("GET", "/", nil)
			request.RemoteAddr = tc.RemoteAddr
			request.Header.Set("X-REAL-IP", tc.XRealIP)
			request.Header.Set("X-FORWARDED-FOR", tc.XForwardedFor)

			id, err := identifier.Identify(request)
			if (err != nil) != tc.ExpectError {
				if tc.ExpectError {
					t.Errorf("Expected an error")
				} else {
					t.Errorf("Expected no error")
				}
			}
			if id != tc.ExpectedIP {
				t.Errorf("Expected: %s, got: %s", tc.ExpectedIP, id)
			}
		})
	}
}
