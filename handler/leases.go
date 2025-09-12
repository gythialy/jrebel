package handler

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"strconv"
	"time"

	"github.com/gythialy/jrebel/module"
	"github.com/gythialy/jrebel/util"
)

const (
	defaultLease = `{
    "serverVersion": "3.2.4",
    "serverProtocolVersion": "1.1",
    "serverGuid": "a1b4aea8-b031-4302-b602-670a990272cb",
    "groupType": "managed",
    "id": 1,
    "licenseType": 1,
    "evaluationLicense": false,
    "signature": "OJE9wGg2xncSb+VgnYT+9HGCFaLOk28tneMFhCbpVMKoC/Iq4LuaDKPirBjG4o394/UjCDGgTBpIrzcXNPdVxVr8PnQzpy7ZSToGO8wv/KIWZT9/ba7bDbA8/RZ4B37YkCeXhjaixpmoyz/CIZMnei4q7oWR7DYUOlOcEWDQhiY=",
    "serverRandomness": "H2ulzLlh7E0=",
    "seatPoolType": "standalone",
    "statusCode": "SUCCESS",
    "offline": false,
    "validFrom": null,
    "validUntil": null,
    "company": "Administrator",
    "orderId": "",
    "zeroIds": [],
    "licenseValidFrom": 1490544001000,
    "licenseValidUntil": 1891839999000
	}`

	defaultValidateConnection = `{
    "serverVersion": "3.2.4",
    "serverProtocolVersion": "1.1",
    "serverGuid": "a1b4aea8-b031-4302-b602-670a990272cb",
    "groupType": "managed",
    "statusCode": "SUCCESS",
    "company": "Administrator",
    "canGetLease": true,
    "licenseType": 1,
    "evaluationLicense": false,
    "seatPoolType": "standalone"
	}`

	defaultLease1 = `{
    "serverVersion": "3.2.4",
    "serverProtocolVersion": "1.1",
    "serverGuid": "a1b4aea8-b031-4302-b602-670a990272cb",
    "groupType": "managed",
    "statusCode": "SUCCESS",
    "msg": null,
    "statusMessage": null,
	"signature": "dGVzdA=="
}`

	randCharset       = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	allRandCharsetLen = len(randCharset)
)

func serverRandomness() (serverRandomness string) {
	b := make([]byte, 11)
	for i := 0; i < 11; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(allRandCharsetLen)))
		if err != nil {
			log.Fatalln("crypto/rand failed:", err)
		}
		b[i] = randCharset[num.Int64()]
	}
	return string(b) + "="
}

type LeaseHandler struct {
	Lease        module.Lease
	Lease1       module.Lease1
	ValidateConn module.ValidateConnection
	singer       *util.Signer
}

func NewHandler() *LeaseHandler {
	handler := LeaseHandler{
		singer: util.NewSigner(),
	}
	if lease, err := util.FromJson[module.Lease](defaultLease); err == nil {
		lease.LicenseValidUntil = time.Now().AddDate(3, 0, 0).UnixMilli()
		handler.Lease = lease
	} else {
		log.Fatalln(err)
	}

	if lease1, err := util.FromJson[module.Lease1](defaultLease1); err == nil {
		handler.Lease1 = lease1
	} else {
		log.Fatalln(err)
	}

	if validateConn, err := util.FromJson[module.ValidateConnection](defaultValidateConnection); err == nil {
		handler.ValidateConn = validateConn
	} else {
		log.Fatalln(err)
	}
	return &handler
}

func (l *LeaseHandler) Leases(w http.ResponseWriter, request *http.Request) {
	params, err := util.UrlParamsFromBody(request)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		_, _ = fmt.Fprintf(w, "%s\n", err)
		return
	}

	clientRandomness := params.Get("randomness")
	username := params.Get("username")
	guid := params.Get("guid")
	if clientRandomness == "" || username == "" || guid == "" {
		w.WriteHeader(http.StatusForbidden)
		_, _ = fmt.Fprint(w)
		return
	}

	offline := false
	if parsed, err := strconv.ParseBool(params.Get("offline")); err == nil {
		offline = parsed
	}

	validFrom := ""
	validUntil := ""
	lease := l.Lease

	if offline {
		clientTime := params.Get("clientTime")
		offlineDays := params.Get("offlineDays")

		startTimeInt, err := strconv.ParseInt(clientTime, 10, 64)
		if err != nil {
			startTimeInt = int64(time.Now().Second()) * 1000
		}

		offlineDaysInt, err := strconv.ParseInt(offlineDays, 10, 64)
		if err != nil {
			offlineDaysInt = 90
		}

		expireTime := startTimeInt + (offlineDaysInt * 24 * 60 * 60 * 1000)
		lease.Offline = offline
		lease.ValidFrom = startTimeInt
		lease.ValidUntil = expireTime

		validFrom = clientTime
		validUntil = strconv.FormatInt(expireTime, 10)
	}

	srvRandomness := serverRandomness()
	if signature, err := l.singer.SignLease(clientRandomness, srvRandomness, guid, offline, validFrom, validUntil); err == nil {
		lease.Signature = signature
		lease.Company = username
		lease.ServerRandomness = srvRandomness
		response(w, lease)
	} else {
		w.WriteHeader(http.StatusForbidden)
		_, _ = fmt.Fprintf(w, "%s\n", err)
	}
}

func (l *LeaseHandler) ValidateConnection(w http.ResponseWriter, _ *http.Request) {
	response(w, l.ValidateConn)
}

func (l *LeaseHandler) Leases1(w http.ResponseWriter, request *http.Request) {
	values, _ := util.UrlParamsFromBody(request)
	company := values.Get("username")
	lease1 := l.Lease1
	if company != "" {
		lease1.Company = company
	}
	response(w, lease1)
}

func response(w http.ResponseWriter, body interface{}) {
	w.Header().Set("content-type", "application/json; charset=utf-8")
	bodyData, err := util.ToJson(body)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		_, _ = fmt.Fprintf(w, "%s\n", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprintf(w, "%s\n", string(bodyData))
}
