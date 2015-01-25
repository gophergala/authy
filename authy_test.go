package authy_test

import (
	"github.com/gophergala/authy"
	"github.com/gophergala/authy/provider"
	. "github.com/smartystreets/goconvey/convey"
	"net/url"
	"testing"
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

var badConfig = authy.Config{
	Providers: map[string]provider.ProviderConfig{
		"invalid": provider.ProviderConfig{
			Key: "my-key",
		},
	},
}

func TestAuthy(t *testing.T) {
	Convey("Invalid configuration", t, func() {
		_, err := authy.NewAuthy(badConfig)
		So(err, ShouldNotEqual, nil)
	})

	Convey("Instanciate Authy", t, func() {
		a, err := authy.NewAuthy(config)
		So(err, ShouldEqual, nil)

		// create a fake session
		session := &FakeSession{
			items: map[interface{}]interface{}{},
		}

		// and a fake oauth
		server := FakeOAuthServer()

		// and even a fake github
		github, _ := provider.GetProvider("github")
		github.AuthorizeURL = server.URL + "/oauth2"
		github.AccessURL = server.URL + "/oauth2"
		provider.RegisterProvider(github)

		Convey("Try to get url for an invalid provider", func() {
			_, err := a.Authorize("bitbucket", session, FakeHttpRequest("http://localhost:2000/authy/bitbucket"))
			So(err, ShouldNotEqual, nil)
		})

		Convey("Get authorization url for GitHub", func() {
			_, err := a.Authorize("github", session, FakeHttpRequest("http://localhost:2000/authy/github"))
			So(err, ShouldEqual, nil)

			Convey("Get access token from an invalid provider", func() {
				_, _, err := a.Access("someone", session, FakeHttpRequest("http://localhost:2000/"))
				So(err, ShouldNotEqual, nil)
			})

			Convey("Get access token from GitHub", func() {
				_, _, err := a.Access("github", session, FakeHttpRequest("http://localhost:2000/authy/github/callback?code=foo&state="+url.QueryEscape(session.Get("authy.github.state").(string))))
				So(err, ShouldEqual, nil)

				Convey("State was deleted so second call should fail", func() {
					_, _, err := a.Access("github", session, FakeHttpRequest("http://localhost:2000/"))
					So(err, ShouldNotEqual, nil)
					server.Close()
				})
			})

		})

	})
}
