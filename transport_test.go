package moneyloverapi_test

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/nhatthm/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/nhatthm/moneyloverapi"
	"github.com/nhatthm/moneyloverapi/pkg/testkit/auth"
)

func TestRoundTripper(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario       string
		tripper        func(*testing.T) moneyloverapi.RoundTripperFunc
		expectedHeader string
	}{
		{
			scenario: "token",
			tripper: func(*testing.T) moneyloverapi.RoundTripperFunc {
				tokenProvider := auth.MockTokenProvider(func(p *auth.TokenProvider) {
					p.On("Token", mock.Anything).
						Return("foobaz", nil)
				})(t)

				return moneyloverapi.TokenRoundTripper(tokenProvider, http.DefaultTransport)
			},
			expectedHeader: "AuthJWT foobaz",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			s := httpmock.New(func(s *httpmock.Server) {
				s.ExpectGet("/").
					WithHeader("Authorization", tc.expectedHeader).
					Return("hello world!")
			})(t)

			req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, s.URL(), nil)
			require.NoError(t, err, "could not create a new request")

			client := http.Client{
				Timeout:   time.Second,
				Transport: tc.tripper(t),
			}
			resp, err := client.Do(req)
			require.NoError(t, err, "could not make a request to mocked server")

			respBody, err := ioutil.ReadAll(resp.Body)
			require.NoError(t, err, "could not read response body")

			err = resp.Body.Close()
			require.NoError(t, err, "could not close response body")

			expectedBody := `hello world!`

			assert.Equal(t, http.StatusOK, resp.StatusCode)
			assert.Equal(t, expectedBody, string(respBody))
		})
	}
}

func TestRoundTripper_Error(t *testing.T) {
	t.Parallel()

	p := auth.MockTokenProvider(func(p *auth.TokenProvider) {
		p.On("Token", mock.Anything).
			Return("", errors.New("token error"))
	})(t)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp, err := moneyloverapi.TokenRoundTripper(p, nil)(req) // nolint:bodyclose

	assert.Nil(t, resp)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "token error")
}
