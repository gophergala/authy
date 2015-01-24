package oauth2

// see http://tools.ietf.org/html/rfc6749

import (
	"fmt"
	"github.com/google/go-querystring/query"
	"github.com/gophergala/authy/provider"
	"net/url"
)

type OAuth2Options struct {
	ClientId     string   `url:"client_id"`
	ResponseType string   `url:"response_type"`
	RedirectURI  string   `url:"redirect_uri,omitempty"`
	Scope        []string `url:"scope,omitempty"`
	State        string   `url:"state,omitempty"`
}

func AuthorizeURL(config provider.ProviderConfig) (dest string, err error) {
	// TODO: load from passed provider
	url, err := url.Parse(config.Provider.AuthorizeURL)
	if err != nil {
		return
	}

	values, err := query.Values(OAuth2Options{
		ClientId:     "foobar",
		ResponseType: "code",
	})

	if err != nil {
		return
	}

	url.RawQuery = values.Encode()

	fmt.Println(url.String())
	dest = url.String()
	return
}
