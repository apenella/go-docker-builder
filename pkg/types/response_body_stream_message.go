package types

import (
	"fmt"
	"strings"
)

// ResponseBodyStreamMessage contains the ImageBuild's body data from buildResponse
type ResponseBodyStreamMessage struct {
	// Aux represents the aux value on response body stream message
	Aux *ResponseBodyStreamAuxMessage `json:"aux"`
	// ErrorDetail
	ErrorDetail *ResponseBodyStreamErrorDetailMessage `json:"errorDetail"`
	// ID identify layer
	ID string `json:"id"`
	// Progress contains the progress bar
	Progress string `json:"progress"`
	// Status represents the status value on response body stream message
	Status string `json:"status"`
	// Stream represents the stream value on response body stream message
	Stream string `json:"stream"`
}

// String return ResponseBodyStreamMessage object as string
func (m *ResponseBodyStreamMessage) String() string {

	if m.Status != "" {
		str := fmt.Sprintf("%s ", separator)
		if m.ID != "" {
			str = fmt.Sprintf("%s %s: ", str, strings.TrimSpace(m.ID))
		}
		str = fmt.Sprintf("%s %s ", str, strings.TrimSpace(m.Status))
		return str
	}
	if m.Stream != "" {
		return strings.TrimSpace(m.Stream)
	}
	if m.Aux != nil {
		return m.Aux.String()
	}
	if m.ErrorDetail != nil {
		return m.ErrorDetail.String()
	}

	return ""
}

// ProgressString returns progress bar
func (m *ResponseBodyStreamMessage) ProgressString() string {
	str := ""

	if m.Progress != "" {
		return strings.TrimSpace(m.Progress)
	}

	return str
}
