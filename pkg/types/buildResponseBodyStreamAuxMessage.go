package types

// BuildResponseBodyStreamAuxMessage contains the ImageBuild's aux data from buildResponse
type BuildResponseBodyStreamAuxMessage struct {
	ID string `json:"ID"`
}

// String return an string with BuildResponseBodyStreamAuxMessage content
func (m *BuildResponseBodyStreamAuxMessage) String() string {

	if m.ID != "" {
		return " \u2023 " + m.ID
	}
	return ""
}
