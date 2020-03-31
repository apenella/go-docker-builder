package types

// BuildResponseBodyStreamAuxMessage contains the ImageBuild's aux data from buildResponse
type ResponseBodyStreamAuxMessage struct {
	// ID is response body stream aux's id
	ID string `json:"ID"`
}

// String return BuildResponseBodyStreamAuxMessage object as string
func (m *ResponseBodyStreamAuxMessage) String() string {

	if m.ID != "" {
		return " \u2023 " + m.ID
	}
	return ""
}
