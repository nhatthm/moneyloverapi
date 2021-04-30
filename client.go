package moneyloverapi

import (
	"net/http"
	"time"

	"github.com/nhatthm/go-clock"
	"github.com/nhatthm/moneyloverapi/internal/api"
	"github.com/nhatthm/moneyloverapi/pkg/auth"
)

const (
	// BaseURL is API Base URL.
	BaseURL = api.DefaultBaseURL
)

// Option configures Client.
type Option func(c *Client)

// Client provides all MoneyLover APIs.
type Client struct {
	api   *api.Client
	token *chainTokenProvider
	clock clock.Clock

	config *config
}

// config is configuration of Client.
type config struct {
	credentials  *chainCredentialsProvider
	tokenStorage auth.TokenStorage
	transport    http.RoundTripper

	baseURL  string
	timeout  time.Duration
	username string
	password string
}

// NewClient initiates a new transaction.Finder.
func NewClient(options ...Option) *Client {
	c := &Client{
		config: &config{
			credentials: chainCredentialsProviders(CredentialsFromEnv()),
			transport:   http.DefaultTransport,

			baseURL: BaseURL,
			timeout: time.Minute,
		},

		token: newChainTokenProvider(),
		clock: clock.New(),
	}

	for _, o := range options {
		o(c)
	}

	c.token.append(initAPITokenProvider(c.config, c.clock))
	c.api = initAPIClient(c.config, c.token)

	return c
}

func initAPITokenProvider(cfg *config, c clock.Clock) auth.TokenProvider {
	cfg.credentials.prepend(Credentials(cfg.username, cfg.password))

	apiToken := newAPITokenProvider(cfg.credentials).
		WithBaseURL(cfg.baseURL).
		WithTimeout(cfg.timeout).
		WithTransport(cfg.transport).
		WithClock(c)

	if cfg.tokenStorage != nil {
		apiToken.WithStorage(cfg.tokenStorage)
	}

	return apiToken
}

func initAPIClient(cfg *config, p auth.TokenProvider) *api.Client {
	c := api.NewClient()
	c.BaseURL = cfg.baseURL
	c.Timeout = cfg.timeout

	c.SetTransport(TokenRoundTripper(p, cfg.transport))

	return c
}
