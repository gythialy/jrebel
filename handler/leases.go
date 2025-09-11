package handler

import (
	"log"
	"math/rand"
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
    "signature": "",
    "serverRandomness": "",
    "seatPoolType": "standalone",
    "statusCode": "SUCCESS",
    "offline": false,
    "validFrom": null,
    "validUntil": null,
    "company": "Administrator",
    "orderId": "",
    "zeroIds": [

    ],
    "licenseValidFrom": 1490544001000,
    "licenseValidUntil": 1691839999000
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
    "company": ""
	}`

	randCharset       = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	allRandCharsetLen = len(randCharset)
)

func serverRandomness() (serverRandomness string) {
	b := make([]byte, 11)
	for i := 0; i < 11; i++ {
		b[i] = randCharset[rand.Intn(allRandCharsetLen)]
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
	params := util.UrlParams(request)
	randomness := params.Get("randomness")
	username := params.Get("username")
	guid := params.Get("guid")

	offline := false
	if offlineParam := params.Get("offline"); offlineParam != "" {
		if parsed, err := strconv.ParseBool(offlineParam); err == nil {
			offline = parsed
		}
	}

	validFrom := "null"
	validUntil := "null"
	lease := l.Lease

	if offline {
		clientTime := params.Get("clientTime")
		offlineDays := params.Get("offlineDays")

		startTimeInt, err := strconv.ParseInt(clientTime, 10, 64)
		if err != nil {
			startTimeInt = time.Now().UnixMilli()
		}

		offlineDaysInt, err := strconv.ParseInt(offlineDays, 10, 64)
		if err != nil {
			offlineDaysInt = int64(90)
		}

		if offlineDaysInt < 1 {
			offlineDaysInt = 1
		}

		expireTime := startTimeInt + (offlineDaysInt * 24 * 60 * 60 * 1000)
		lease.Offline = offline
		lease.ValidFrom = startTimeInt
		lease.ValidUntil = expireTime

		validFrom = strconv.FormatInt(startTimeInt, 10)
		validUntil = strconv.FormatInt(expireTime, 10)
	}

	if randomness == "" || username == "" || guid == "" {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	serverRandomness := serverRandomness()
	if signature, err := l.singer.SignLease(randomness, serverRandomness, guid, offline, validFrom, validUntil); err == nil {
		lease.Signature = signature
		lease.Company = username
		lease.ServerRandomness = serverRandomness
		response(w, lease)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (l *LeaseHandler) ValidateConnection(w http.ResponseWriter, _ *http.Request) {
	response(w, l.ValidateConn)
}

func (l *LeaseHandler) Leases1(w http.ResponseWriter, request *http.Request) {
	values := util.UrlParams(request)
	company := values.Get("username")
	lease1 := l.Lease1
	if company != "" {
		lease1.Company = company
	}
	response(w, lease1)
}

func response(w http.ResponseWriter, body interface{}) {
	if content, err := util.ToJson(body); err == nil {
		w.WriteHeader(http.StatusOK)
		w.Header().Add("content-type", "application/json;charset=utf-8")
		_, _ = w.Write([]byte(content))
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
