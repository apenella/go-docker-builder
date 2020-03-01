package types

// BuildResponseBodyStreamAuxMessage contains the ImageBuild's aux data from buildResponse
type BuildResponseBodyStreamAuxMessage struct {
	// ID is response body stream aux's id
	ID string `json:"ID"`
}

// String return BuildResponseBodyStreamAuxMessage object as string
func (m *BuildResponseBodyStreamAuxMessage) String() string {

	if m.ID != "" {
		return " \u2023 " + m.ID
	}
	return ""
}
