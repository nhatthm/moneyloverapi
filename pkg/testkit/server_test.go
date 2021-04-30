package testkit

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/nhatthm/httpmock"
	"github.com/stretchr/testify/assert"

	"github.com/nhatthm/moneyloverapi/pkg/auth"
)

func TestServer_WithAccessToken(t *testing.T) {
	t.Parallel()

	expected := auth.Token(`ACCESS_TOKEN`)

	s := (&Server{}).WithAccessToken(`ACCESS_TOKEN`)

	assert.Equal(t, expected, s.AccessToken())
}

func TestServer_Expect(t *testing.T) {
	t.Parallel()

	accessToken := `ACCESS_TOKEN`

	s := MockEmptyServer(func(s *Server) {
		s.WithAccessToken(accessToken)
		s.Expect(http.MethodGet, "/")
	})(t)

	expectedHeaders := httpmock.Header{
		"Authorization": fmt.Sprintf("AuthJWT %s", accessToken),
	}

	assert.Equal(t, http.MethodGet, s.ExpectedRequests[0].Method)
	assert.Equal(t, httpmock.Exact("/"), s.ExpectedRequests[0].RequestURI)

	for key, matcher := range s.ExpectedRequests[0].RequestHeader {
		assert.True(t, matcher.Match(expectedHeaders[key]))
	}

	s.ResetExpectations()
}

func TestServer_ExpectAliases(t *testing.T) {
	t.Parallel()

	accessToken := `ACCESS_TOKEN`

	testCases := []struct {
		scenario       string
		mockServer     func(s *Server)
		expectedMethod string
	}{
		{
			scenario: "GET",
			mockServer: func(s *Server) {
				s.ExpectGet("/")
			},
			expectedMethod: http.MethodGet,
		},
		{
			scenario: "HEAD",
			mockServer: func(s *Server) {
				s.ExpectHead("/")
			},
			expectedMethod: http.MethodHead,
		},
		{
			scenario: "POST",
			mockServer: func(s *Server) {
				s.ExpectPost("/")
			},
			expectedMethod: http.MethodPost,
		},
		{
			scenario: "PUT",
			mockServer: func(s *Server) {
				s.ExpectPut("/")
			},
			expectedMethod: http.MethodPut,
		},
		{
			scenario: "PATCH",
			mockServer: func(s *Server) {
				s.ExpectPatch("/")
			},
			expectedMethod: http.MethodPatch,
		},
		{
			scenario: "DELETE",
			mockServer: func(s *Server) {
				s.ExpectDelete("/")
			},
			expectedMethod: http.MethodDelete,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			s := MockEmptyServer(tc.mockServer)(t)
			s.WithAccessToken(accessToken)

			expectedHeaders := httpmock.Header{
				"Authorization": fmt.Sprintf("AuthJWT %s", accessToken),
			}

			assert.Equal(t, tc.expectedMethod, s.ExpectedRequests[0].Method)
			assert.Equal(t, httpmock.Exact("/"), s.ExpectedRequests[0].RequestURI)

			for key, matcher := range s.ExpectedRequests[0].RequestHeader {
				assert.True(t, matcher.Match(expectedHeaders[key]))
			}

			s.ResetExpectations()
		})
	}
}
