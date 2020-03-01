package types

import "strings"

// BuildResponseBodyStreamMessage contains the ImageBuild's body data from buildResponse
type BuildResponseBodyStreamMessage struct {
	// Status represents the status value on response body stream message
	Status string `json:"status"`
	// Stream represents the stream value on response body stream message
	Stream string `json:"stream"`
	// Aux represents the aux value on response body stream message
	Aux *BuildResponseBodyStreamAuxMessage `json:"aux"`
}

// String return BuildResponseBodyStreamMessage object as string
func (m *BuildResponseBodyStreamMessage) String() string {

	if m.Stream != "" {
		return strings.TrimSpace(m.Stream)
	}
	if m.Status != "" {
		return " \u2023 " + strings.TrimSpace(m.Status)
	}
	if m.Aux != nil {
		return m.Aux.String()
	}
	return ""
}
