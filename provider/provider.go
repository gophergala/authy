package provider

import (
	"errors"
	"fmt"
)

type Provider struct {
	Name           string
	RequestURL     string
	AuthorizeURL   string
	AccessURL      string
	OAuth          int
	ScopeDelimiter string
}

type ProviderConfig struct {
	Provider Provider `json:"-"`
	Key      string   `json:"key" xml:"key,attr"`
	Secret   string   `json:"secret" xml:"secret,attr"`
	Scope    []string `json:"scope" xml:"scope,attr"`
	State    string   `json:"state" xml:"state,attr"`
}

var defaultProviders = map[string]Provider{
	"github": Provider{
		Name:         "github",
		AuthorizeURL: "https://github.com/login/oauth/authorize",
		AccessURL:    "https://github.com/login/oauth/access_token",
		OAuth:        2,
	},
}

var customProviders = map[string]Provider{}

// Get a provider by name
func GetProvider(name string) (Provider, error) {
	if provider, ok := customProviders[name]; ok == true {
		return provider, nil
	}

	if provider, ok := defaultProviders[name]; ok == true {
		return provider, nil
	}

	return Provider{}, errors.New(fmt.Sprintf("unknown provider: %s", name))
}

// Register a custom provider, takes precedence on default providers
func RegisterProvider(provider Provider) error {
	if provider.Name == "" {
		return errors.New("custom provider's name cannot be empty")
	}

	customProviders[provider.Name] = provider
	return nil
}
