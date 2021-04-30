package testkit

import (
	"net/http"
	"sync"

	"github.com/nhatthm/httpmock"
	"github.com/nhatthm/moneyloverapi/pkg/auth"
)

// Request is an alias of httpmock.Request.
type Request = httpmock.Request

// Server is a wrapped httpmock.Server to provide more functionalities for testing MoneyLover APIs.
type Server struct {
	*httpmock.Server

	requestToken auth.Token
	accessToken  auth.Token
	refreshToken auth.Token

	mu sync.Mutex
}

// WithRequestToken sets the accessToken.
func (s *Server) WithRequestToken(token string) *Server {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.requestToken = auth.Token(token)

	return s
}

// RequestToken returns the requestToken.
func (s *Server) RequestToken() auth.Token {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.requestToken
}

// WithAccessToken sets the accessToken.
func (s *Server) WithAccessToken(token string) *Server {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.accessToken = auth.Token(token)

	return s
}

// AccessToken returns the accessToken.
func (s *Server) AccessToken() auth.Token {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.accessToken
}

// WithRefreshToken sets the refreshToken.
func (s *Server) WithRefreshToken(token string) *Server {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.refreshToken = auth.Token(token)

	return s
}

// RefreshToken returns the refreshToken.
func (s *Server) RefreshToken() auth.Token {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.refreshToken
}

// Expect expects a request with Bearer Authorization.
//
//    Server.Expect(http.MethodGet, "/path").
func (s *Server) Expect(method string, requestURI interface{}) *Request {
	return s.Server.Expect(method, requestURI).
		WithHeader("Authorization", func() httpmock.Matcher {
			return httpmock.Exactf("AuthJWT %s", s.AccessToken())
		})
}

// ExpectGet expects a request with Bearer Authorization.
//
//   Server.ExpectGet("/path")
func (s *Server) ExpectGet(requestURI interface{}) *Request {
	return s.Expect(http.MethodGet, requestURI)
}

// ExpectHead expects a request with Bearer Authorization.
//
//   Server.ExpectHead("/path")
func (s *Server) ExpectHead(requestURI interface{}) *Request {
	return s.Expect(http.MethodHead, requestURI)
}

// ExpectPost expects a request with Bearer Authorization.
//
//   Server.ExpectPost("/path")
func (s *Server) ExpectPost(requestURI interface{}) *Request {
	return s.Expect(http.MethodPost, requestURI)
}

// ExpectPut expects a request with Bearer Authorization.
//
//   Server.ExpectPut("/path")
func (s *Server) ExpectPut(requestURI interface{}) *Request {
	return s.Expect(http.MethodPut, requestURI)
}

// ExpectPatch expects a request with Bearer Authorization.
//
//   Server.ExpectPatch("/path")
func (s *Server) ExpectPatch(requestURI interface{}) *Request {
	return s.Expect(http.MethodPatch, requestURI)
}

// ExpectDelete expects a request with Bearer Authorization.
//
//   Server.ExpectDelete("/path")
func (s *Server) ExpectDelete(requestURI interface{}) *Request {
	return s.Expect(http.MethodDelete, requestURI)
}

// NewServer creates a new Server.
func NewServer(t TestingT) *Server {
	s := &Server{
		Server: httpmock.NewServer(t),
	}

	s.WithDefaultResponseHeaders(httpmock.Header{
		"Content-Type": "application/json",
	})

	return s
}
