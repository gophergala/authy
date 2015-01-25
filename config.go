package authy

import (
	"github.com/gophergala/authy/provider"
)

// Configuration for authy, is already mapped for being parsed by encoding/json
//
// Example JSON file:
//
//   {
//     "callback": "/login/success",
//     "providers": {
//       "github": {
//         "key": "be148a4abf2796b3a8e1",
//         "secret": "1bbf884bbf79ef21fef03410389eb451300abd84",
//         "scope": ["repo", "email"]
//       }
//    }
//  }
type Config struct {
	// Where to redirect the user for login if supported by the middleware (defaults to /login)
	PathLogin string `json:"login"`
	// Which base route to use (defaults to /authy)
	BasePath string `json:"base_path"`
	// Where the user is redirected by default after a successful auth
	Callback string `json:"callback"`
	// A list of providers
	Providers map[string]provider.ProviderConfig `json:"providers"`
}
