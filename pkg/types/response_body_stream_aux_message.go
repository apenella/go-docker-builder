package types

import "fmt"

// BuildResponseBodyStreamAuxMessage contains the ImageBuild's aux data from buildResponse
type ResponseBodyStreamAuxMessage struct {
	// ID is response body stream aux's id
	ID string `json:"ID"`
}

// String return BuildResponseBodyStreamAuxMessage object as string
func (m *ResponseBodyStreamAuxMessage) String() string {

	if m.ID != "" {
		return fmt.Sprintf(" %s %s", LayerMessagePrefix, m.ID)
	}
	return ""
}
