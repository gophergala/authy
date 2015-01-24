package oauth2

// see http://tools.ietf.org/html/rfc6749

import (
	"github.com/google/go-querystring/query"
	"github.com/gophergala/authy/provider"
	"net/http"
	"net/url"
	"strings"
)

type OAuth2Options struct {
	ClientId     string `url:"client_id"`
	ResponseType string `url:"response_type"`
	RedirectURI  string `url:"redirect_uri,omitempty"`
	Scope        string `url:"scope,omitempty"`
	State        string `url:"state,omitempty"`
}

func genCallbackURL(config provider.ProviderConfig, r *http.Request) string {
	var redirectURI = url.URL{
		Host: r.Host,
		Path: "/connect/" + config.Provider.Name + "/callback",
	}

	if _, ok := r.Header["X-HTTPS"]; r.TLS != nil || ok == true {
		redirectURI.Scheme = "https"
	} else {
		redirectURI.Scheme = "http"
	}

	return redirectURI.String()
}

func AuthorizeURL(config provider.ProviderConfig, r *http.Request) (dest string, err error) {
	authUrl, err := url.Parse(config.Provider.AuthorizeURL)
	if err != nil {
		return
	}

	values, err := query.Values(OAuth2Options{
		ClientId:     config.Key,
		ResponseType: "code",
		RedirectURI:  genCallbackURL(config, r),
		Scope:        strings.Join(config.Scope, ","),
		State:        config.State,
	})

	if err != nil {
		return
	}

	authUrl.RawQuery = values.Encode()
	dest = authUrl.String()
	return
}
