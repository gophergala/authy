// Implements several middlewares for using Authy with Martini
package authy

import (
	"github.com/go-martini/martini"
	"github.com/gophergala/authy"
	"github.com/martini-contrib/sessions"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

type Config authy.Config
type Token authy.Token

// Takes an Authy config and returns a middleware to use with martini
// See examples below
func Authy(config Config) martini.Handler {
	baseRoute := "/authy"
	if config.BasePath != "" {
		baseRoute = config.BasePath
	}

	// should be moved in the authy package
	if config.PathLogin == "" {
		config.PathLogin = "/login"
	}

	authRoute := regexp.MustCompile("^" + baseRoute + "/([^/#?]+)")
	callbackRoute := regexp.MustCompile("^" + baseRoute + "/([^/]+)/callback")
	authy, err := authy.NewAuthy(authy.Config(config))

	// due to the way middleware are used, it's the cleanest? way to deal with this?
	if err != nil {
		panic(err)
	}

	return func(s sessions.Session, c martini.Context, w http.ResponseWriter, r *http.Request) {
		c.Map(config)

		// if we are already logged, ignore login route matching
		if tokenValue := s.Get("authy.token.value"); tokenValue != nil {
			c.Map(Token{
				Provider: s.Get("authy.provider").(string),
				Value:    tokenValue.(string),
				Scope:    strings.Split(s.Get("authy.token.scope").(string), ","),
			})
			return
		}

		matches := authRoute.FindStringSubmatch(r.URL.Path)
		if len(matches) > 0 && matches[0] == r.URL.Path {
			redirectUrl, err := authy.Authorize(matches[1], s, r)
			if err != nil {
				panic(err)
			}

			// redirect user to oauth website
			http.Redirect(w, r, redirectUrl, http.StatusFound)
			return
		}

		matches = callbackRoute.FindStringSubmatch(r.URL.Path)
		if len(matches) > 0 && matches[0] == r.URL.Path {
			token, redirectUrl, err := authy.Access(matches[1], s, r)
			if err != nil {
				panic(err)
			}

			// save token in session
			s.Set("authy.provider", matches[1])
			s.Set("authy.token.value", token.Value)
			s.Set("authy.token.scope", strings.Join(token.Scope, ","))

			http.Redirect(w, r, redirectUrl, http.StatusFound)
			return
		}
	}
}

// Use this middleware on the routes where you need the user to be logged in
func LoginRequired() martini.Handler {
	return func(config Config, s sessions.Session, w http.ResponseWriter, r *http.Request) {
		if tokenValue := s.Get("authy.token.value"); tokenValue == nil {
			next := url.QueryEscape(r.URL.RequestURI())
			http.Redirect(w, r, config.PathLogin+"?next="+next, http.StatusFound)
		}
	}
}
