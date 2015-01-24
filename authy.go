package authy

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/gophergala/authy/oauth2"
	"github.com/gophergala/authy/provider"
	"log"
	"net/http"
	"regexp"
	"time"
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
				// setup State to prevent any forgery (TODO: break this up in it's own function that his framework aware)
				rawState := make([]byte, 16)
				_, err := rand.Read(rawState)
				if err != nil {
					http.Error(rw, "Internal Server Error", http.StatusInternalServerError)
					log.Println("error", err)
					return
				}

				sum := md5.Sum(rawState)
				stateSum := hex.EncodeToString(sum[0:16])
				http.SetCookie(rw, &http.Cookie{
					Name:    "state",
					Value:   stateSum,
					Expires: time.Now().Add(10 * time.Minute),
				})
				providerConfig.State = hex.EncodeToString(rawState)

				// generate authorisation URL
				redirectUrl, err := oauth2.AuthorizeURL(providerConfig, r)

				if err != nil {
					http.Error(rw, "Internal Server Error", http.StatusInternalServerError)
					log.Println("error", err)
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
