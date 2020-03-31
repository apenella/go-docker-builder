package types

import "strings"

// ResponseBodyStreamMessage contains the ImageBuild's body data from buildResponse
type ResponseBodyStreamMessage struct {
	// Status represents the status value on response body stream message
	Status string `json:"status"`
	// Stream represents the stream value on response body stream message
	Stream      string                                `json:"stream"`
	ErrorDetail *ResponseBodyStreamErrorDetailMessage `json:"errorDetail"`
	// Aux represents the aux value on response body stream message
	Aux *ResponseBodyStreamAuxMessage `json:"aux"`
}

// String return ResponseBodyStreamMessage object as string
func (m *ResponseBodyStreamMessage) String() string {

	if m.Stream != "" {
		return strings.TrimSpace(m.Stream)
	}
	if m.Status != "" {
		return " \u2023 " + strings.TrimSpace(m.Status)
	}
	if m.Aux != nil {
		return m.Aux.String()
	}
	if m.ErrorDetail != nil {
		return m.ErrorDetail.String()
	}

	return ""
}
