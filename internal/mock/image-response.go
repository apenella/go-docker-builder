package mock

import (
	"context"
	"io"
	"iter"

	"github.com/moby/moby/api/types/jsonstream"
)

// ImageResponse is a test fake implementing client.ImagePullResponse and
// client.ImagePushResponse (both embed io.ReadCloser plus JSONMessages/Wait).
type ImageResponse struct {
	io.ReadCloser
}

// NewImageResponse wraps an io.ReadCloser as an ImagePull/Push response fake
func NewImageResponse(rc io.ReadCloser) ImageResponse {
	return ImageResponse{ReadCloser: rc}
}

// JSONMessages implements client.ImagePullResponse / client.ImagePushResponse
func (r ImageResponse) JSONMessages(ctx context.Context) iter.Seq2[jsonstream.Message, error] {
	return func(yield func(jsonstream.Message, error) bool) {}
}

// Wait implements client.ImagePullResponse / client.ImagePushResponse
func (r ImageResponse) Wait(ctx context.Context) error {
	return nil
}
