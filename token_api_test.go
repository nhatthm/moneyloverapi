package moneyloverapi

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	clock "go.nhat.io/clock/mock"

	"github.com/nhatthm/moneyloverapi/pkg/auth"
	"github.com/nhatthm/moneyloverapi/pkg/testkit"
	authMock "github.com/nhatthm/moneyloverapi/pkg/testkit/auth"
)

func TestApiTokenProvider_GetToken(t *testing.T) {
	t.Parallel()

	username := "user@example.org"
	password := "123456"
	cred := Credentials(username, password)
	storageKey := username

	testCases := []struct {
		scenario      string
		mockServer    testkit.ServerMocker
		configure     func(t *testing.T, p *apiTokenProvider)
		expectedError string
	}{
		{
			scenario:   "could not get token from storage",
			mockServer: testkit.MockEmptyServer(),
			configure: func(t *testing.T, p *apiTokenProvider) { //nolint: thelper
				s := authMock.MockTokenStorage(func(s *authMock.TokenStorage) {
					s.On("Get", context.Background(), storageKey).
						Return(auth.OAuthToken{}, errors.New("get token error"))
				})(t)

				p.WithStorage(s)
			},
			expectedError: "could not get token from storage: get token error",
		},
		{
			scenario: "could not get login url",
			mockServer: testkit.MockEmptyServer(
				testkit.WithAuthLoginURLFailure(),
			),
			expectedError: "unexpected response: unexpected response status: 500 Internal Server Error",
		},
		{
			scenario: "could not login",
			mockServer: testkit.MockEmptyServer(
				testkit.WithAuthLoginURLSuccess(),
				testkit.WithAuthTokenFailure(username, password),
			),
			expectedError: "unexpected response: unexpected response status: 401 Unauthorized",
		},
		{
			scenario: "could not set token",
			mockServer: testkit.MockEmptyServer(
				testkit.WithAuthLoginURLSuccess(),
				testkit.WithAuthTokenSuccess(username, password),
			),
			configure: func(t *testing.T, p *apiTokenProvider) { //nolint: thelper
				s := authMock.MockTokenStorage(func(s *authMock.TokenStorage) {
					s.On("Get", context.Background(), storageKey).
						Return(auth.OAuthToken{}, nil)

					s.On("Set", context.Background(), storageKey, mock.Anything).
						Return(errors.New("set token error"))
				})(t)

				p.WithStorage(s)
			},
			expectedError: "could not persist token to storage: set token error",
		},
		{
			scenario: "success",
			mockServer: testkit.MockEmptyServer(
				testkit.WithAuthLoginURLSuccess(),
				testkit.WithAuthTokenSuccess(username, password),
			),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			s := tc.mockServer(t)
			p := newAPITokenProvider(cred).
				WithBaseURL(s.URL()).
				WithTimeout(time.Second)

			if tc.configure != nil {
				tc.configure(t, p)
			}

			token, err := p.Token(context.Background())

			if tc.expectedError == "" {
				assert.Equal(t, s.AccessToken(), token)
				assert.NoError(t, err)
			} else {
				assert.Empty(t, token)
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestApiTokenProvider_GetToken_MissingCredentials(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario      string
		credentials   CredentialsProvider
		expectedError string
	}{
		{
			scenario:      "missing username",
			credentials:   Credentials("", "password"),
			expectedError: "could not get token: missing username",
		},
		{
			scenario:      "missing password",
			credentials:   Credentials("username", ""),
			expectedError: "unexpected response: unexpected response status: 403 Forbidden",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			p := newAPITokenProvider(tc.credentials)

			token, err := p.Token(context.Background())

			assert.Empty(t, token)
			assert.EqualError(t, err, tc.expectedError)
		})
	}
}

func TestApiTokenProvider_GetTokenFromCache(t *testing.T) {
	t.Parallel()

	username := "user@example.org"
	password := "123456"
	cred := Credentials(username, password)
	timestamp := time.Now()

	s := testkit.MockServer(username, password)(t)

	c := clock.Mock(func(c *clock.Clock) {
		c.On("Now").Return(timestamp.Add(4 * time.Minute)).Once()
	})(t)

	p := newAPITokenProvider(cred).
		WithBaseURL(s.URL()).
		WithTimeout(time.Second).
		WithClock(c)

	// 1st try.
	token, err := p.Token(context.Background())

	assert.Equal(t, s.AccessToken(), token)
	assert.NotEmpty(t, string(token))
	assert.NoError(t, err)

	// 2nd try.
	token, err = p.Token(context.Background())

	assert.Equal(t, s.AccessToken(), token)
	assert.NotEmpty(t, string(token))
	assert.NoError(t, err)
}

func TestApiTokenProvider_TokenExpired(t *testing.T) {
	t.Parallel()

	username := "user@example.org"
	password := "123456"
	cred := Credentials(username, password)
	timestamp := time.Now()

	s := testkit.MockEmptyServer(
		testkit.WithAuthSuccess(username, password),
		testkit.WithAuthSuccess(username, password),
	)(t)

	c := clock.Mock(func(c *clock.Clock) {
		c.On("Now").Return(timestamp.AddDate(0, 0, 30)).Once()
	})(t)

	p := newAPITokenProvider(cred).
		WithBaseURL(s.URL()).
		WithTimeout(time.Second).
		WithClock(c)

	// 1st try.
	token1, err := p.Token(context.Background())

	assert.Equal(t, s.AccessToken(), token1)
	assert.NotEmpty(t, string(token1))
	assert.NoError(t, err)

	// 2nd try.
	token2, err := p.Token(context.Background())

	assert.Equal(t, s.AccessToken(), token2)
	assert.NotEqual(t, token1, token2)
	assert.NotEmpty(t, string(token2))
	assert.NoError(t, err)
}
