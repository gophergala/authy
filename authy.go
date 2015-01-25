// Package authy implements the base methods for implementing oauth providers
// It is recommended instead to use one of the middlewares already provided by the package
//
// Middlewares:
//
// * Martini: https://github.com/gophergala/authy/martini
//
// For a full guide visit https://github.com/gophergala/authy
package authy

import (
	"errors"
	"fmt"
	"github.com/gophergala/authy/oauth2"
	"github.com/gophergala/authy/provider"
	"net/http"
)

// Authy represents the current configuration and cached provider data
type Authy struct {
	config    Config
	providers map[string]provider.ProviderConfig
}

// Token is returned on a successful auth
type Token struct {
	// The actual value of the token
	Value string
	// The scopes returned by the provider, some providers may allow the user to change the scope of an auth request
	// Make sure to check the available scopes before doing queries on their webservices
	Scope []string
}

// Parse the configuration and build the list of providers, return an Authy instance
func NewAuthy(config Config) (Authy, error) {
	var availableProviders = map[string]provider.ProviderConfig{}

	// load all providers
	for providerName, providerConfig := range config.Providers {
		providerData, err := provider.GetProvider(providerName)
		if err != nil {
			return Authy{}, err
		}
		providerConfig.Provider = providerData
		availableProviders[providerName] = providerConfig
	}

	return Authy{
		config:    config,
		providers: availableProviders,
	}, nil
}

// Generate a CSRF token and store it in the provided session object, return the authorisation URL
// It should be noted that the session object should prevent the user from seeing the sum generated
func (a Authy) Authorize(providerName string, session Session, r *http.Request) (string, error) {
	providerConfig, ok := a.providers[providerName]
	if ok != true {
		return "", errors.New(fmt.Sprintf("unknown provider %s", providerName))
	}

	if providerConfig.Provider.OAuth == 2 {
		state, err := oauth2.NewState()
		if err != nil {
			return "", err
		}

		// save authentication state in session
		session.Set("authy."+providerName+".state", state)
		providerConfig.State = state

		// generate authorisation URL
		redirectUrl, err := oauth2.AuthorizeURL(providerConfig, r)

		if err != nil {
			return "", err
		}

		return redirectUrl, nil
	}

	return "", errors.New("Not Implemented")
}

// Check the CSRF token then query the distant provider for an access token using the code that was provided by the
// authorization API
func (a Authy) Access(providerName string, session Session, r *http.Request) (Token, string, error) {
	providerConfig, ok := a.providers[providerName]
	if ok != true {
		return Token{}, "", errors.New(fmt.Sprintf("unknown provider %s", providerName))
	}

	if providerConfig.Provider.OAuth == 2 {
		// check the state parameter against CSRF
		state := session.Get("authy." + providerName + ".state")
		if state == nil {
			return Token{}, "", errors.New("state token is not set in session, possible CSRF")
		}

		stateParam := r.URL.Query().Get("state")
		if stateParam != state.(string) {
			return Token{}, "", errors.New("invalid state param provided, possible CSRF")
		}
		// we don't need it anymore
		session.Delete("authy." + providerName + ".state")

		code := r.URL.Query().Get("code")
		if code == "" {
			return Token{}, "", errors.New("code was not found in the query parameters")
		}

		// retrieve access token from provider
		token, err := oauth2.GetAccessToken(providerConfig, r)
		if err != nil {
			return Token{}, "", err
		}

		// provide the proper callback URL
		redirectUrl := a.config.Callback
		if providerConfig.Callback != "" {
			redirectUrl = providerConfig.Callback
		}

		// return the token
		return Token{
			Value: token.AccessToken,
			Scope: token.Scope,
		}, redirectUrl, nil
	}

	return Token{}, "", errors.New("Not Implemented")
}
