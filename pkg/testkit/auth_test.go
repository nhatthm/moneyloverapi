package testkit_test

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nhatthm/moneyloverapi/pkg/testkit"
)

func TestWithAuthLoginURLFailure(t *testing.T) {
	t.Parallel()

	s := testkit.MockEmptyServer(
		testkit.WithAuthLoginURLFailure(),
	)(t)

	code, _, _, _ := request(t, s.URL(), http.MethodPost, "/user/login-url", nil, nil) // nolint:dogsled

	assert.Equal(t, http.StatusInternalServerError, code)
}

func TestWithAuthLoginURLSuccess(t *testing.T) {
	t.Parallel()

	s := testkit.MockEmptyServer(
		testkit.WithAuthLoginURLSuccess(),
	)(t)

	code, _, body, _ := request(t, s.URL(), http.MethodPost, "/user/login-url", nil, nil)

	expectedBody := fmt.Sprintf(
		`{"msg":"redirect_url","action":"user_login","data":{"status":true,"request_token":%q,"login_url":"%s/auth?callback=%s\u0026client=kHiZbFQOw5LV\u0026token=%s"}}`,
		s.RequestToken(),
		s.URL(),
		url.QueryEscape(s.URL()),
		s.RequestToken(),
	)

	assert.NotEmpty(t, s.RequestToken())
	assert.Equal(t, http.StatusOK, code)
	assert.Equal(t, expectedBody, string(body))
}

func TestWithAuthTokenFailure(t *testing.T) {
	t.Parallel()

	username := "user@example.org"
	password := "123456"
	requestToken := "REQUEST_TOKEN"

	s := testkit.MockEmptyServer(
		func(s *testkit.Server) {
			s.WithRequestToken(requestToken)
		},
		testkit.WithAuthTokenFailure(username, password),
	)(t)

	requestHeader := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", requestToken),
		"Client":        "kHiZbFQOw5LV",
		"Content-Type":  "application/json; charset=utf-8",
	}
	requestBody := fmt.Sprintf(`{"client_info":true,"email":%q,"password":%q}`,
		username, password,
	)

	code, _, _, _ := request(t, s.URL(), http.MethodPost, "/token", requestHeader, []byte(requestBody)) // nolint:dogsled

	assert.Equal(t, http.StatusUnauthorized, code)
}

func TestWithAuthTokenSuccess(t *testing.T) {
	t.Parallel()

	username := "user@example.org"
	password := "123456"
	requestToken := "REQUEST_TOKEN"

	s := testkit.MockEmptyServer(
		func(s *testkit.Server) {
			s.WithRequestToken(requestToken)
		},
		testkit.WithAuthTokenSuccess(username, password),
	)(t)

	requestHeader := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", requestToken),
		"Client":        "kHiZbFQOw5LV",
		"Content-Type":  "application/json; charset=utf-8",
	}
	requestBody := fmt.Sprintf(`{"client_info":true,"email":%q,"password":%q}`,
		username, password,
	)

	code, _, body, _ := request(t, s.URL(), http.MethodPost, "/token", requestHeader, []byte(requestBody))

	expectedBody := fmt.Sprintf(regexp.QuoteMeta(fmt.Sprintf(
		`{"status":true,"access_token":"%s","expire":"%%s","refresh_token":%q}`,
		s.AccessToken(),
		s.RefreshToken(),
	)), "[0-9]+")

	assert.NotEmpty(t, s.AccessToken())
	assert.NotEmpty(t, s.RefreshToken())
	assert.Equal(t, http.StatusOK, code)
	assert.Regexp(t, expectedBody, string(body))
}
