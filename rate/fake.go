package rate

import "net/http"

type FakeIdentifier struct {
	FnIdentify func(r *http.Request) (string, error)
}

func (f FakeIdentifier) Identify(r *http.Request) (string, error) {
	return f.FnIdentify(r)
}

type FakeLimiter struct {
	FnAllow func(id string) (Response, error)
}

func (f FakeLimiter) Allow(id string) (Response, error) {
	return f.FnAllow(id)
}
