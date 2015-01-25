// This package partially implements OAuth2 for Authy
package oauth2

// see http://tools.ietf.org/html/rfc6749

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/google/go-querystring/query"
	"github.com/gophergala/authy/provider"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type authorizationRequest struct {
	ClientId     string `url:"client_id"`
	ResponseType string `url:"response_type"`
	RedirectURI  string `url:"redirect_uri,omitempty"`
	Scope        string `url:"scope,omitempty"`
	State        string `url:"state,omitempty"`
}

type accessTokenRequest struct {
	ClientId     string `url:"client_id"`
	ClientSecret string `url:"client_secret"`
	GrantType    string `url:"grant_type"`
	Code         string `url:"code"`
	RedirectURI  string `url:"redirect_uri,omitempty"`
}

type Token struct {
	AccessToken string
	Scope       []string
	Type        string
}

// standard oauth2 error (http://tools.ietf.org/html/rfc6749#section-5.2)
type Error struct {
	Code        string
	Description string
	URI         string
}

func (err Error) Error() string {
	msg := err.Code
	if err.Description != "" {
		msg += ": " + err.Description
	}
	if err.URI != "" {
		msg += " (see " + err.URI + ")"
	}
	return msg
}

func genCallbackURL(config provider.ProviderConfig, r *http.Request) string {
	var redirectURI = url.URL{
		Host: r.Host,
		Path: r.URL.Path + "/callback",
	}

	if _, ok := r.Header["X-HTTPS"]; r.TLS != nil || ok == true {
		redirectURI.Scheme = "https"
	} else {
		redirectURI.Scheme = "http"
	}

	return redirectURI.String()
}

// create a new random token for the CSRF check
func NewState() (string, error) {
	rawState := make([]byte, 16)
	_, err := rand.Read(rawState)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(rawState), nil
}

// Generates the proper authorization URL for the given service
func AuthorizeURL(config provider.ProviderConfig, r *http.Request) (dest string, err error) {
	// subdomain support
	baseUrl := config.Provider.AuthorizeURL
	if config.Provider.Subdomain == true {
		if config.Subdomain == "" {
			err = errors.New(fmt.Sprintf("provider %s expects the config to contain your subdomain", config.Provider.Name))
			return
		}
		baseUrl = strings.Replace(baseUrl, "[subdomain]", config.Subdomain, -1)
	}

	authUrl, err := url.Parse(baseUrl)
	if err != nil {
		return
	}

	values, err := query.Values(authorizationRequest{
		ClientId:     config.Key,
		ResponseType: "code",
		RedirectURI:  genCallbackURL(config, r),
		Scope:        strings.Join(config.Scope, config.Provider.ScopeDelimiter),
		State:        config.State,
	})

	// custom parameters
	if len(config.CustomParameters) > 0 {
		for _, name := range config.Provider.CustomParameters {
			if value, ok := config.CustomParameters[name]; ok == true {
				values.Set(name, value)
			}
		}
	}

	if err != nil {
		return
	}

	authUrl.RawQuery = values.Encode()
	dest = authUrl.String()
	return
}

// Query the remote service for an access token
func GetAccessToken(config provider.ProviderConfig, r *http.Request) (token Token, err error) {
	queryValues, err := query.Values(accessTokenRequest{
		ClientId:     config.Key,
		ClientSecret: config.Secret,
		Code:         r.URL.Query().Get("code"),
		GrantType:    "authorization_code",
		RedirectURI:  genCallbackURL(config, r),
	})

	if err != nil {
		return
	}

	resp, err := http.PostForm(config.Provider.AccessURL, queryValues)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	values, err := url.ParseQuery(string(body))
	if err != nil {
		return
	}

	if _, ok := values["error"]; ok == true {
		err = Error{
			Code:        values["error"][0],
			Description: values["error_description"][0],
			URI:         values["error_uri"][0],
		}
		return
	}

	// everything went A-OK!
	token.AccessToken = values["access_token"][0]
	token.Type = values["token_type"][0]

	// TODO: maybe store the scope with the state token and set it there if not returned
	// by the service
	if scope, ok := values["scope"]; ok == true {
		token.Scope = strings.Split(scope[0], config.Provider.ScopeDelimiter)
	}

	return
}
