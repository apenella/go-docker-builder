package types

import "strings"

// BuildResponseBodyStreamMessage contains the ImageBuild's body data from buildResponse
type BuildResponseBodyStreamMessage struct {
	Status string                             `json:"status"`
	Stream string                             `json:"stream"`
	Aux    *BuildResponseBodyStreamAuxMessage `json:"aux"`
}

// String return an string with BuildResponseBodyStreamMessage content
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
