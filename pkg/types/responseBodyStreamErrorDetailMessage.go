package types

// BuildResponseBodyStreamAuxMessage contains the ImageBuild's aux data from buildResponse
type ResponseBodyStreamErrorDetailMessage struct {
	// ID is response body stream aux's id
	Message string `json:"message"`
}

// String return BuildResponseBodyStreamAuxMessage object as string
func (m *ResponseBodyStreamErrorDetailMessage) String() string {

	if m.Message != "" {
		return "[ERROR] " + m.Message
	}
	return ""
}
