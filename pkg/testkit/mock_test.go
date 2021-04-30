package testkit_test

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"testing"
	"time"

	"github.com/nhatthm/httpmock"
	"github.com/nhatthm/moneyloverapi/pkg/testkit"
	"github.com/stretchr/testify/assert"
)

func TestMockServer(t *testing.T) {
	t.Parallel()

	username := "user@example.org"
	password := "123456"

	s := testkit.MockServer(username, password)(t)

	// 1st step: get login url.
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

	// 2nd step: login.
	requestHeader := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", s.RequestToken()),
		"Client":        "kHiZbFQOw5LV",
		"Content-Type":  "application/json; charset=utf-8",
	}
	requestBody := fmt.Sprintf(`{"client_info":true,"email":%q,"password":%q}`,
		username, password,
	)

	code, _, body, _ = request(t, s.URL(), http.MethodPost, "/token", requestHeader, []byte(requestBody))

	expectedBody = fmt.Sprintf(regexp.QuoteMeta(fmt.Sprintf(
		`{"status":true,"access_token":"%s","expire":"%%s","refresh_token":%q}`,
		s.AccessToken(),
		s.RefreshToken(),
	)), "[0-9]+")

	assert.NotEmpty(t, s.AccessToken())
	assert.NotEmpty(t, s.RefreshToken())
	assert.Equal(t, http.StatusOK, code)
	assert.Regexp(t, expectedBody, string(body))
}

func TestMockEmptyServer(t *testing.T) {
	t.Parallel()

	accessToken := "ACCESS_TOKEN"

	s := testkit.MockEmptyServer(func(s *testkit.Server) {
		s.WithAccessToken(accessToken)
		s.ExpectGet("/").Return(`{}`)
	})(t)

	requestHeader := map[string]string{
		"Authorization": fmt.Sprintf("AuthJWT %s", s.AccessToken()),
	}

	code, headers, _, _ := request(t, s.URL(), http.MethodGet, "/", requestHeader, nil)

	expectedHeaders := map[string]string{
		"Content-Type": "application/json",
	}

	assert.Equal(t, http.StatusOK, code)
	httpmock.AssertHeaderContains(t, headers, expectedHeaders)
}

// nolint:thelper
func request(
	t *testing.T,
	baseURL string,
	method, uri string,
	headers map[string]string,
	body []byte,
) (int, map[string]string, []byte, time.Duration) {
	return httpmock.DoRequest(t,
		method, baseURL+uri,
		headers, body,
	)
}
