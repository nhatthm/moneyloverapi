package testkit

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.nhat.io/httpmock"
	plannerMock "go.nhat.io/httpmock/mock/planner"
	"go.nhat.io/httpmock/planner"

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

	var e planner.Expectation

	p := plannerMock.Mock(func(p *plannerMock.Planner) {
		p.On("Expect", mock.Anything).
			Run(func(args mock.Arguments) {
				e = args[0].(planner.Expectation) //nolint: errcheck
			})

		p.On("IsEmpty").Return(true)
	})(t)

	MockEmptyServer(func(s *Server) {
		s.WithPlanner(p).
			WithAccessToken(accessToken)

		s.Expect(http.MethodGet, "/")
	})(t)

	expectedHeaders := httpmock.Header{
		"Authorization": fmt.Sprintf("AuthJWT %s", accessToken),
	}

	assert.Equal(t, http.MethodGet, e.Method())
	assert.Equal(t, httpmock.Exact("/"), e.URIMatcher())

	requestHeader := e.HeaderMatcher()

	assert.Len(t, requestHeader, 1)

	for key, m := range requestHeader {
		matched, err := m.Match(expectedHeaders[key])

		assert.True(t, matched)
		assert.NoError(t, err)
	}
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

			var e planner.Expectation

			p := plannerMock.Mock(func(p *plannerMock.Planner) {
				p.On("Expect", mock.Anything).
					Run(func(args mock.Arguments) {
						e = args[0].(planner.Expectation) //nolint: errcheck
					})

				p.On("IsEmpty").Return(true)
			})(t)

			s := MockEmptyServer(func(s *Server) {
				s.WithPlanner(p)
			}, tc.mockServer)(t)

			s.WithAccessToken(accessToken)

			expectedHeaders := httpmock.Header{
				"Authorization": fmt.Sprintf("AuthJWT %s", accessToken),
			}

			assert.Equal(t, tc.expectedMethod, e.Method())
			assert.Equal(t, httpmock.Exact("/"), e.URIMatcher())

			requestHeader := e.HeaderMatcher()

			assert.Len(t, requestHeader, 1)

			for key, m := range requestHeader {
				matched, err := m.Match(expectedHeaders[key])

				assert.True(t, matched)
				assert.NoError(t, err)
			}
		})
	}
}
