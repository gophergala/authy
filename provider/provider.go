// The provider package contains the list of service providers and their base configuration
package provider

import (
	"errors"
	"fmt"
)

// Contains implementation details to be used by Authy
type Provider struct {
	Name             string
	RequestURL       string
	AuthorizeURL     string
	AccessURL        string
	OAuth            int
	ScopeDelimiter   string
	Subdomain        bool
	CustomParameters []string
}

// Those keys are imported from your config, set the proper ones based on your provider's oauth information
type ProviderConfig struct {
	Provider         Provider          `json:"-"`
	Key              string            `json:"key"`
	Secret           string            `json:"secret"`
	Scope            []string          `json:"scope"`
	State            string            `json:"state"`
	Callback         string            `json:"callback"`
	Subdomain        string            `json:"subdomain"`
	CustomParameters map[string]string `json:"custom_parameters"`
}

var customProviders = map[string]Provider{}

// Get a provider by name
func GetProvider(name string) (Provider, error) {
	provider, ok := customProviders[name]
	if ok != true {
		provider, ok = defaultProviders[name]
	}

	if ok != true {
		return Provider{}, errors.New(fmt.Sprintf("unknown provider: %s", name))
	}

	if provider.ScopeDelimiter == "" {
		provider.ScopeDelimiter = ","
	}
	return provider, nil
}

// Register a custom provider, takes precedence on default providers
func RegisterProvider(provider Provider) error {
	if provider.Name == "" {
		return errors.New("custom provider's name cannot be empty")
	}

	customProviders[provider.Name] = provider
	return nil
}
