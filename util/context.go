package util

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
)

func GetUrlParams(request *http.Request) url.Values {
	body, _ := io.ReadAll(request.Body)
	log.Printf("{%s} {%s}", request.Method, request.RequestURI+"?"+string(body))
	values, _ := url.ParseRequestURI(request.RequestURI + "?" + string(body))
	return values.Query()
}

func WriteJson(w http.ResponseWriter, jsonObject interface{}) {
	body, _ := json.Marshal(jsonObject)
	w.Header().Add("content-type", "application/json;charset=utf-8")
	_, _ = w.Write(body)
}
