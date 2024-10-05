package module

type ValidateConnection struct {
	ServerVersion         string `json:"serverVersion"`
	ServerProtocolVersion string `json:"serverProtocolVersion"`
	ServerGuid            string `json:"serverGuid"`
	GroupType             string `json:"groupType"`
	StatusCode            string `json:"statusCode"`
	Company               string `json:"company"`
	CanGetLease           bool   `json:"canGetLease"`
	LicenseType           int    `json:"licenseType"`
	EvaluationLicense     bool   `json:"evaluationLicense"`
	SeatPoolType          string `json:"seatPoolType"`
}
