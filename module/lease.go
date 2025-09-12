package module

type Lease struct {
	ServerVersion         string        `json:"serverVersion"`
	ServerProtocolVersion string        `json:"serverProtocolVersion"`
	ServerGUID            string        `json:"serverGuid"`
	GroupType             string        `json:"groupType"`
	ID                    int           `json:"id"`
	LicenseType           int           `json:"licenseType"`
	EvaluationLicense     bool          `json:"evaluationLicense"`
	Signature             string        `json:"signature"`
	ServerRandomness      string        `json:"serverRandomness"`
	SeatPoolType          string        `json:"seatPoolType"`
	StatusCode            string        `json:"statusCode"`
	Offline               bool          `json:"offline"`
	ValidFrom             int64         `json:"validFrom"`
	ValidUntil            int64         `json:"validUntil"`
	Company               string        `json:"company"`
	OrderID               string        `json:"orderId"`
	ZeroIds               []interface{} `json:"zeroIds"`
	LicenseValidFrom      int64         `json:"licenseValidFrom"`
	LicenseValidUntil     int64         `json:"licenseValidUntil"`
}
