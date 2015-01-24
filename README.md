Authy
=====

Authy is a go library that acts as an authentication middleware for [net/http](http://golang.org/pkg/net/http),
it aims to provide drop-in support for most OAuth 1 and 2 providers.

Usage
-----

`server.go`
```go
package main

import (
	"encoding/json"
	"github.com/gophergala/authy"
	"net/http"
	"os"
)

type Config struct {
	Authy authy.Config `json:"authy"`
}

func main() {
	f, err := os.Open("config.json")
	if err != nil {
		panic(err)
	}

	decoder := json.NewDecoder(f)

	var config Config
	decoder.Decode(&config)

    // authy.Setup will automatically register routes with http.DefaultServerMux
	authy.Setup(config.Authy)
	http.ListenAndServe(":5000", nil)
}
```

`config.json`
```json
{
	"authy": {
		"providers": {
			"github": {
				"key": "my-key",
				"secret": "my-secret",
				"scope": ["repo", "email"]
			}
		}
	}
}
```