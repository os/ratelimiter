package rate

import (
	"net"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

// Identifier provides a method to identify the source of the given request.
type Identifier interface {
	Identify(r *http.Request) (string, error)
}

// IPIdentifier identifies the source of the requests based on their IP
// addresses.
type IPIdentifier struct{}

// NewIPIdentifier returns an instance of IPIdentifier.
func NewIPIdentifier() IPIdentifier {
	return IPIdentifier{}
}

// Identify returns the IP address of the given request.
func (i IPIdentifier) Identify(r *http.Request) (string, error) {
	var ip net.IP

	realIP := r.Header.Get("X-REAL-IP")
	ip = net.ParseIP(realIP)
	if ip != nil {
		return ip.String(), nil
	}

	forwardedFor := r.Header.Get("X-FORWARDED-FOR")
	items := strings.Split(forwardedFor, ",")
	for _, item := range items {
		ip = net.ParseIP(item)
		if ip != nil {
			return ip.String(), nil
		}
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", errors.Wrap(err, "Remote address could not be parsed")
	}
	ip = net.ParseIP(host)
	if ip != nil {
		return ip.String(), nil
	}

	return "", errors.New("IP address could not be identified")
}
