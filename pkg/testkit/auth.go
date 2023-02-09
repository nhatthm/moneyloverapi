package testkit

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"go.nhat.io/httpmock"
	"go.nhat.io/httpmock/matcher"

	"github.com/nhatthm/moneyloverapi/internal/api"
	"github.com/nhatthm/moneyloverapi/internal/types"
)

func expectAuthLoginURL(s *Server) Expectation {
	return s.Server.ExpectPost("/user/login-url")
}

func expectLogin(s *Server, username, password string) Expectation {
	return s.ExpectPost("/token").
		WithHeader("Authorization", func() matcher.Matcher {
			return httpmock.Exactf("Bearer %s", s.RequestToken())
		}).
		WithHeader("client", client).
		WithHeader("Content-Type", "application/json; charset=utf-8").
		WithBodyJSON(api.TokenRequest{
			ClientInfo: true,
			Email:      username,
			Password:   password,
		})
}

// WithAuthLoginURLFailure expects a request for login URL and returns a 500.
func WithAuthLoginURLFailure() ServerOption {
	return func(s *Server) {
		expectAuthLoginURL(s).
			ReturnCode(http.StatusInternalServerError)
	}
}

// WithAuthLoginURLSuccess expects a request for login URL.
func WithAuthLoginURLSuccess() ServerOption {
	return func(s *Server) {
		expectAuthLoginURL(s).
			Run(func(r *http.Request) ([]byte, error) {
				s.WithRequestToken(requestToken())
				token := string(s.RequestToken())

				loginURL, _ := url.Parse(s.URL() + "/auth") //nolint: errcheck
				loginParams := loginURL.Query()
				loginParams.Set("client", client)
				loginParams.Set("token", token)
				loginParams.Set("callback", s.URL())
				loginURL.RawQuery = loginParams.Encode()

				res := api.LoginResponse{
					Error:  0,
					Msg:    "redirect_url",
					Action: "user_login",
					Data: api.LoginResponseData{
						Status:       true,
						RequestToken: token,
						LoginURL:     loginURL.String(),
					},
				}

				return json.Marshal(res)
			})
	}
}

// WithAuthTokenFailure expects a request for login and returns a 401.
func WithAuthTokenFailure(username, password string) ServerOption {
	return func(s *Server) {
		expectLogin(s, username, password).
			ReturnCode(http.StatusUnauthorized)
	}
}

// WithAuthTokenSuccess expects a request for login.
func WithAuthTokenSuccess(username, password string) ServerOption {
	return func(s *Server) {
		expectLogin(s, username, password).
			Run(func(r *http.Request) ([]byte, error) {
				s.WithAccessToken(accessToken())
				s.WithRefreshToken(refreshToken())

				res := api.TokenResponse{
					Status:       true,
					AccessToken:  string(s.AccessToken()),
					Expire:       types.UnixString(time.Now().Add(time.Hour)),
					RefreshToken: string(s.RefreshToken()),
					ClientInfo:   nil,
				}

				return json.Marshal(res)
			})
	}
}

// WithAuthSuccess expects a success login workflow.
func WithAuthSuccess(username, password string) ServerOption {
	return func(s *Server) {
		WithAuthLoginURLSuccess()(s)
		WithAuthTokenSuccess(username, password)(s)
	}
}
