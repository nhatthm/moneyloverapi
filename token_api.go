package moneyloverapi

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/bool64/ctxd"
	"go.nhat.io/clock"

	"github.com/nhatthm/moneyloverapi/internal/api"
	"github.com/nhatthm/moneyloverapi/pkg/auth"
)

var (
	// ErrUsernameIsEmpty indicates that the username is empty.
	ErrUsernameIsEmpty = errors.New("missing username")
	// ErrPasswordIsEmpty indicates that the username is empty.
	ErrPasswordIsEmpty = errors.New("missing password")
)

var _ auth.TokenProvider = (*apiTokenProvider)(nil)

var emptyToken = auth.OAuthToken{}

type apiTokenProvider struct {
	api         *api.Client
	transport   http.RoundTripper
	credentials CredentialsProvider
	storage     auth.TokenStorage
	clock       clock.Clock

	mu sync.Mutex
}

func (p *apiTokenProvider) oauthClient(baseURL string) *api.Client {
	c := api.NewClient()
	c.BaseURL = baseURL
	c.Timeout = p.api.Timeout
	c.InstrumentCtxFunc = p.api.InstrumentCtxFunc
	c.SetTransport(p.transport)

	return c
}

func (p *apiTokenProvider) getToken(ctx context.Context, key string) (auth.OAuthToken, error) {
	return p.storage.Get(ctx, key)
}

func (p *apiTokenProvider) setToken(ctx context.Context, key string, res api.TokenResponse) (auth.OAuthToken, error) {
	token := auth.OAuthToken{
		AccessToken: auth.Token(res.AccessToken),
		ExpiresAt:   res.Expire.Time(),
	}

	if err := p.storage.Set(ctx, key, token); err != nil {
		return auth.OAuthToken{}, err
	}

	return token, nil
}

func (p *apiTokenProvider) getLoginURL(ctx context.Context) (*url.URL, error) {
	res, err := p.api.PostUserLoginURL(ctx, api.PostUserLoginURLRequest{})
	if err != nil {
		return nil, ctxd.WrapError(ctx, err, "unexpected response")
	}

	return url.Parse(res.ValueOK.Data.LoginURL)
}

func (p *apiTokenProvider) login(ctx context.Context, loginURL url.URL) (*api.TokenResponse, error) {
	loginParams := loginURL.Query()
	authorization := loginParams.Get("token")
	client := loginParams.Get("client")

	password := p.credentials.Password()
	if password == "" {
		return nil, ctxd.WrapError(ctx, ErrPasswordIsEmpty, "could not get token")
	}

	oauth := p.oauthClient(fmt.Sprintf("%s://%s", loginURL.Scheme, loginURL.Host))

	res, err := oauth.PostToken(ctx, api.PostTokenRequest{
		Authorization: fmt.Sprintf("Bearer %s", authorization),
		Client:        &client,
		Body: &api.TokenRequest{
			ClientInfo: true,
			Email:      p.credentials.Username(),
			Password:   password,
		},
	})
	if err != nil {
		return nil, ctxd.WrapError(ctx, err, "unexpected response")
	}

	return res.ValueOK, err
}

func (p *apiTokenProvider) get(ctx context.Context, key string) (auth.Token, error) {
	loginURL, err := p.getLoginURL(ctx)
	if err != nil {
		return "", err
	}

	res, err := p.login(ctx, *loginURL)
	if err != nil {
		return "", err
	}

	token, err := p.setToken(ctx, key, *res)
	if err != nil {
		return "", ctxd.WrapError(ctx, err, "could not persist token to storage")
	}

	return token.AccessToken, nil
}

func (p *apiTokenProvider) WithBaseURL(baseURL string) *apiTokenProvider {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.api.BaseURL = baseURL

	return p
}

func (p *apiTokenProvider) WithTimeout(timeout time.Duration) *apiTokenProvider {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.api.Timeout = timeout

	return p
}

//nolint:unparam
func (p *apiTokenProvider) WithStorage(storage auth.TokenStorage) *apiTokenProvider {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.storage = storage

	return p
}

func (p *apiTokenProvider) WithTransport(transport http.RoundTripper) *apiTokenProvider {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.transport = transport
	p.api.SetTransport(p.transport)

	return p
}

func (p *apiTokenProvider) WithClock(clock clock.Clock) *apiTokenProvider {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.clock = clock

	return p
}

func (p *apiTokenProvider) Token(ctx context.Context) (auth.Token, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	username := p.credentials.Username()
	if username == "" {
		return "", ctxd.WrapError(ctx, ErrUsernameIsEmpty, "could not get token")
	}

	key := username

	token, err := p.getToken(ctx, key)
	if err != nil {
		return "", ctxd.WrapError(ctx, err, "could not get token from storage")
	}

	if token != emptyToken && !token.IsExpired(p.clock.Now()) {
		return token.AccessToken, nil
	}

	return p.get(ctx, key)
}

func newAPITokenProvider(
	credentials CredentialsProvider,
) *apiTokenProvider {
	c := api.NewClient()
	c.BaseURL = BaseURL
	c.Timeout = time.Minute

	p := &apiTokenProvider{
		api:         c,
		credentials: credentials,
		storage:     NewInMemoryTokenStorage(),
		clock:       clock.New(),
	}

	p.WithTransport(http.DefaultTransport)

	return p
}
