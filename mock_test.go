package authy_test

import (
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"net/url"
)

// a fake session object
type FakeSession struct {
	items map[interface{}]interface{}
}

func (f *FakeSession) Get(key interface{}) interface{} {
	return f.items[key]
}

func (f *FakeSession) Set(key interface{}, val interface{}) {
	f.items[key] = val
}

func (f *FakeSession) Delete(key interface{}) {
	delete(f.items, key)
}

// generate a fake http request
func FakeHttpRequest(requestUrl string) *http.Request {
	parsedUrl, _ := url.Parse(requestUrl)
	return &http.Request{
		URL: parsedUrl,
	}
}

// fake oauth2 service
func FakeOAuthServer() (s *httptest.Server) {
	r := mux.NewRouter()

	r.HandleFunc("/oauth2", func(rw http.ResponseWriter, r *http.Request) {
		values := url.Values{}
		values.Set("access_token", "fakeaccesstoken")
		values.Set("scope", r.URL.Query().Get("scope"))
		values.Set("token_type", "example")
		rw.Write([]byte(values.Encode()))
	})

	s = httptest.NewServer(r)
	return
}
