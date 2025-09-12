package util

import (
	"io"
	"log"
	"net/http"
	"net/url"
)

func UrlParamsFromBody(r *http.Request) (url.Values, error) {
	var body []byte
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	content := string(body)
	u := url.URL{
		Scheme:     "",
		Opaque:     "",
		User:       nil,
		Host:       "",
		Path:       "",
		RawPath:    "",
		ForceQuery: false,
		RawQuery:   content,
		Fragment:   "",
	}
	log.Printf("[%s%s]: %s", r.Method, r.URL.Path, content)
	return u.Query(), err
}
