package module

type Lease1 struct {
	ServerVersion         string `json:"serverVersion"`
	ServerProtocolVersion string `json:"serverProtocolVersion"`
	ServerGuid            string `json:"serverGuid"`
	GroupType             string `json:"groupType"`
	StatusCode            string `json:"statusCode"`
	Msg                   string `json:"msg,omitempty"`
	StatusMessage         string `json:"statusMessage,omitempty"`
	Company               string `json:"company"`
}
