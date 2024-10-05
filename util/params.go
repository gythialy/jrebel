package util

import (
	"io"
	"log"
	"net/http"
	"net/url"
)

func UrlParams(request *http.Request) url.Values {
	body, _ := io.ReadAll(request.Body)
	u := request.RequestURI + "?" + string(body)
	log.Printf("[%s] %s", request.Method, u)
	values, _ := url.ParseRequestURI(u)
	return values.Query()
}
