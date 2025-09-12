package handler

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"html"
	"net/http"

	"github.com/gythialy/jrebel/util"
)

const asmPrivateKey = `
-----BEGIN PRIVATE KEY-----
MIIBVAIBADANBgkqhkiG9w0BAQEFAASCAT4wggE6AgEAAkEAt5yrcHAAjhglnCEn
6yecMWPeUXcMyo0+itXrLlkpcKIIyqPw546bGThhlb1ppX1ySX/OUA4jSakHekNP
5eWPawIDAQABAkBbr9pUPTmpuxkcy9m5LYBrkWk02PQEOV/fyE62SEPPP+GRhv4Q
Fgsu+V2GCwPQ69E3LzKHPsSNpSosIHSO4g3hAiEA54JCn41fF8GZ90b9L5dtFQB2
/yIcGX4Xo7bCvl8DaPMCIQDLCUN8YiXppydqQ+uYkTQgvyq+47cW2wcGumRS46dd
qQIhAKp2v5e8AMj9ROFO5B6m4SsVrIkwFICw17c0WzDRxTEBAiAYDmftk990GLcF
0zhV4lZvztasuWRXE+p4NJtwasLIyQIgVKzknJe8VOt5a3shCMOyysoNEg+YAt02
O98RPCU0nJg=
-----END PRIVATE KEY-----`

func PingHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("content-type", "text/html; charset=utf-8")
	params, _ := util.UrlParamsFromBody(req)

	salt := params.Get("salt")
	if salt == "" {
		w.WriteHeader(403)
		_, _ = fmt.Fprint(w)
	} else {
		xmlContent := "<PingResponse><message></message><responseCode>OK</responseCode><salt>" + html.EscapeString(salt) + "</salt></PingResponse>"
		signature, err := signWithMd5([]byte(xmlContent))
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
		} else {
			body := "<!-- " + hex.EncodeToString(signature) + " -->\n" + xmlContent
			_, _ = w.Write([]byte(body))
		}
	}
}

func ObtainTicketHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("content-type", "application/json; charset=utf-8")
	params, _ := util.UrlParamsFromBody(req)
	salt := params.Get("salt")
	username := params.Get("username")
	prolongationPeriod := "607875500"

	if salt == "" || username == "" {
		w.WriteHeader(http.StatusForbidden)
	} else {
		w.WriteHeader(http.StatusOK)
		xmlContent := "<ObtainTicketResponse><message></message><prolongationPeriod>" + prolongationPeriod + "</prolongationPeriod><responseCode>OK</responseCode><salt>" + html.EscapeString(salt) + "</salt><ticketId>1</ticketId><ticketProperties>licensee=" + html.EscapeString(username) + "\tlicenseType=0\t</ticketProperties></ObtainTicketResponse>"
		signature, err := signWithMd5([]byte(xmlContent))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.Header().Add("content-type", "application/xml; charset=utf-8")
			body := "<!-- " + hex.EncodeToString(signature) + " -->\n" + xmlContent
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(body))
		}
	}
}

func ReleaseTicketHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("content-type", "text/html; charset=utf-8")
	params, _ := util.UrlParamsFromBody(req)

	salt := params.Get("salt")
	if salt == "" {
		w.WriteHeader(http.StatusForbidden)
	} else {
		xmlContent := "<ReleaseTicketResponse><message></message><responseCode>OK</responseCode><salt>" + salt + "</salt></ReleaseTicketResponse>"
		signature, err := signWithMd5([]byte(xmlContent))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprintf(w, "%s\n", err)
		} else {
			w.Header().Add("content-type", "application/xml; charset=utf-8")
			body := "<!-- " + hex.EncodeToString(signature) + " -->\n" + xmlContent
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(body))
		}
	}
}

func signWithMd5(data []byte) ([]byte, error) {
	asnPrivateKeyBlock, _ := pem.Decode([]byte(asmPrivateKey))
	privateKey, err := x509.ParsePKCS8PrivateKey(asnPrivateKeyBlock.Bytes)
	if err != nil {
		return nil, err
	}
	h := crypto.Hash.New(crypto.MD5)
	h.Write(data)
	hashed := h.Sum(nil)

	return rsa.SignPKCS1v15(rand.Reader, privateKey.(*rsa.PrivateKey), crypto.MD5, hashed)
}
