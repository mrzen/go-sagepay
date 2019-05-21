package sagepay

import (
	"context"
	"os"
)

// CredentialsProvider is a source for API credentials.
type CredentialsProvider interface {
	GetCredentials(ctx context.Context) (*Credentials, error)
}

// Credentials are the API credentials for an API client
type Credentials struct {
	Username string
	Password string
}

// StaticCredentialsProvider gets credentials statically from a struct.
type StaticCredentialsProvider struct {
	Credentials
}

// EnvironmentCredentialsProvider gets credentials from the OS environment
type EnvironmentCredentialsProvider struct{}

// GetCredentials gets the current credentials from the environment
func (e EnvironmentCredentialsProvider) GetCredentials(ctx context.Context) (*Credentials, error) {
	return &Credentials{
		Username: os.Getenv("SAGE_USERNAME"),
		Password: os.Getenv("SAGE_PASSWORD"),
	}, nil
}

// StaticCredentials creates a new set of static credentials
func StaticCredentials(username, password string) StaticCredentialsProvider {
	return StaticCredentialsProvider{
		Credentials: Credentials{
			Username: username,
			Password: password,
		},
	}
}

// GetCredentials gets the current credentials
func (s StaticCredentialsProvider) GetCredentials(ctx context.Context) (*Credentials, error) {
	return &s.Credentials, nil
}
