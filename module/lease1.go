package module

type Lease1 struct {
	ServerVersion         string      `json:"serverVersion"`
	ServerProtocolVersion string      `json:"serverProtocolVersion"`
	ServerGUID            string      `json:"serverGuid"`
	GroupType             string      `json:"groupType"`
	StatusCode            string      `json:"statusCode"`
	Company               string      `json:"company"`
	Msg                   interface{} `json:"msg"`
	StatusMessage         interface{} `json:"statusMessage"`
	Signature             string      `json:"signature"`
}
