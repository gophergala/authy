package authy

import "github.com/gophergala/authy/provider"

type Config struct {
	Providers map[string]provider.ProviderConfig `json:"providers"`
}
