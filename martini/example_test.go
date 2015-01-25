package authy_test

import (
	"github.com/go-martini/martini"
	"github.com/gophergala/authy/martini"
	"github.com/gophergala/authy/provider"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
)

var config = authy.Config{
	PathLogin: "/login",
	Callback:  "/login/success",
	Providers: map[string]provider.ProviderConfig{
		"github": provider.ProviderConfig{
			Key:    "my-key",
			Secret: "my-secret",
			Scope:  []string{"repo", "user:mail"},
		},
	},
}

func ExampleAuthy() {
	m := martini.Classic()

	// the session need to be set for the CSRF token system to work
	m.Use(sessions.Sessions("authy", sessions.NewCookieStore([]byte("no one will guess this passphrase"))))
	m.Use(authy.Authy(config))
}

func ExampleLoginRequired() {
	m := martini.Classic()

	// the session need to be set for the CSRF token system to work
	m.Use(sessions.Sessions("authy", sessions.NewCookieStore([]byte("no one will guess this passphrase"))))
	m.Use(authy.Authy(config))

	// use authy.LoginRequired to redirect the user to the login page if not logged in
	m.Get("/profile", authy.LoginRequired(), func(token authy.Token, r render.Render) {
		r.HTML(200, "callback", token)
	})
}
