package api

type Response struct {
	Status   string      `json:"status"`
	Data     interface{} `json:"data"`
	Metadata interface{} `json:"metadata,omitempty"`
}

// VersionResponse used by version handler
type VersionResponse struct {
	Version string `json:"version" example:"12345`
}

// ErrorResponse general error response
type ErrorResponse struct {
	Message string        `json:"message"`
	Details string        `json:"details"`
	Service string        `json:"service"`
	Type    string        `json:"type"`
	Reason  []ErrorReason `json:"reason"`
}

// ErrorReason used by ErrorResponse to map specific error reasons
type ErrorReason struct {
	Field string      `json:"field"`
	Error string      `json:"error"`
	Value interface{} `json:"value"`
	Type  string      `json:"type"`
}
