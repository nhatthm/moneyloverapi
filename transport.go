package moneyloverapi

import (
	"fmt"
	"net/http"

	"github.com/nhatthm/moneyloverapi/pkg/auth"
)

// RoundTripperFunc is an inline http.RoundTripper.
type RoundTripperFunc func(*http.Request) (*http.Response, error)

// RoundTrip satisfies RoundTripperFunc.
func (fn RoundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

// TokenRoundTripper sets Bearer Authorization header to the given request with a token given by a auth.TokenProvider.
func TokenRoundTripper(p auth.TokenProvider, tripper http.RoundTripper) RoundTripperFunc {
	return func(r *http.Request) (*http.Response, error) {
		token, err := p.Token(r.Context())
		if err != nil {
			return nil, err
		}

		r.Header.Add("Authorization", fmt.Sprintf("AuthJWT %s", token))

		return tripper.RoundTrip(r)
	}
}
