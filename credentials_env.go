package moneyloverapi

import "os"

const (
	envUsername = "MONEYLOVER_USERNAME"
	envPassword = "MONEYLOVER_PASSWORD"
)

var _ CredentialsProvider = (*envCredentialsProvider)(nil)

// envCredentialsProvider provides username and password from environment variables.
type envCredentialsProvider struct{}

// Username provides a username from moneyloverapi.envUsername variable.
func (p *envCredentialsProvider) Username() string {
	return os.Getenv(envUsername)
}

// Password provides a password from moneyloverapi.envPassword variable.
func (p *envCredentialsProvider) Password() string {
	return os.Getenv(envPassword)
}

// CredentialsFromEnv initiates a new credentials provider.
func CredentialsFromEnv() CredentialsProvider {
	return &envCredentialsProvider{}
}
