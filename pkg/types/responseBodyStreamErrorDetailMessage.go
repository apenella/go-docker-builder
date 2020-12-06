package types

import "fmt"

const (
	errorPrefix = "[ERROR]"
)

// BuildResponseBodyStreamAuxMessage contains the ImageBuild's aux data from buildResponse
type ResponseBodyStreamErrorDetailMessage struct {
	// ID is response body stream aux's id
	Message string `json:"message"`
}

// String return BuildResponseBodyStreamAuxMessage object as string
func (m *ResponseBodyStreamErrorDetailMessage) String() string {

	if m.Message != "" {
		return fmt.Sprintf("%s %s", errorPrefix, m.Message)
	}

	return ""
}
