package authy

import (
	"fmt"
	"github.com/gophergala/authy/oauth2"
	"github.com/gophergala/authy/provider"
	"log"
	"net/http"
	"regexp"
)

var providerUrl = regexp.MustCompile("^/connect/([0-9a-z_-]+)$")

func Setup(config Config) error {
	return SetupMux(config, http.DefaultServeMux)
}

// Setup Authy on the given ServeMux, if none given, uses DefaultServeMux
func SetupMux(config Config, mux *http.ServeMux) error {
	var availableProviders = map[string]provider.ProviderConfig{}

	if mux == nil {
		mux = http.DefaultServeMux
	}

	// load all providers
	for providerName, providerConfig := range config.Providers {
		providerData, err := provider.GetProvider(providerName)
		if err != nil {
			return err
		}
		providerConfig.Provider = providerData
		availableProviders[providerName] = providerConfig
	}

	// our handler
	mux.HandleFunc("/connect/", func(rw http.ResponseWriter, r *http.Request) {
		// Parse incoming requests and redirect to the right providers
		if providerUrl.MatchString(r.URL.Path) {
			providerName := providerUrl.FindStringSubmatch(r.URL.Path)[1]
			providerConfig, ok := availableProviders[providerName]

			fmt.Println(providerConfig, ok, providerName)

			if ok != true {
				http.NotFound(rw, r)
				return
			}

			if providerConfig.Provider.OAuth == 2 {
				redirectUrl, err := oauth2.AuthorizeURL(providerConfig)
				if err != nil {
					http.Error(rw, "Internal Server Error", http.StatusInternalServerError)
					log.Fatalln(err)
					return
				}

				// debug
				rw.Header()["Content-Type"] = []string{"text/html"}
				fmt.Fprintf(rw, "redirect url: <a href=\"%s\">%s</a>", redirectUrl, redirectUrl)
			}
		}
	})

	return nil
}
