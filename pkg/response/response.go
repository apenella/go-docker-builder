package response

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/ahmetb/go-cursor"
	errors "github.com/apenella/go-common-utils/error"
	transformer "github.com/apenella/go-common-utils/transformer/string"
	"github.com/apenella/go-docker-builder/pkg/types"
	orderedmap "github.com/wk8/go-ordered-map"
)

// ExecuteOptions is a function to set executor options
type ResponseOptions func(*DefaultResponse)

type DefaultResponse struct {
	// Writer is where is written the command stdout
	Writer io.Writer

	// Transformers
	Transformers []transformer.TransformerFunc
}

func NewDefaultResponse(opts ...ResponseOptions) *DefaultResponse {
	r := &DefaultResponse{}

	for _, opt := range opts {
		opt(r)
	}

	if r.Writer == nil {
		r.Writer = os.Stdout
	}

	return r
}

// WithWriter set the writer to be used by DefaultExecutor
func WithWriter(w io.Writer) ResponseOptions {
	return func(r *DefaultResponse) {
		r.Writer = w
	}
}

// WithTransformers add trasformes
func WithTransformers(trans ...transformer.TransformerFunc) ResponseOptions {
	return func(r *DefaultResponse) {
		if r.Transformers == nil {
			r.Transformers = trans
		} else {
			r.Transformers = append(r.Transformers, trans...)
		}
	}
}

// Print
func (r *DefaultResponse) Print(reader io.ReadCloser) error {
	scanner := bufio.NewScanner(reader)
	errorContext := "(response::Print)"
	lineBefore := ""
	lines := orderedmap.New()
	numLayers := 0

	for scanner.Scan() {

		streamMessage := &types.ResponseBodyStreamMessage{}
		line := scanner.Bytes()
		err := json.Unmarshal(line, &streamMessage)
		if err != nil {
			return errors.New(errorContext, fmt.Sprintf("Error unmarshalling line '%s'", string(line)), err)
		}

		streamMessageStr := streamMessage.String()

		if streamMessageStr != lineBefore && streamMessageStr != "" {
			if streamMessage.ID != "" {
				// override layer outputs on pull or push messages
				fmt.Fprintf(r.Writer, "%s%s\n", cursor.MoveUp(numLayers+1), cursor.ClearEntireLine())

				lines.Set(streamMessage.ID, fmt.Sprint(streamMessage.String(), streamMessage.ProgressString()))
				for line := lines.Oldest(); line != nil; line = line.Next() {
					r.Fwriteln(fmt.Sprintf("%s%s", line.Value, cursor.ClearLineRight()))
				}

				numLayers = lines.Len()
			} else {
				r.Fwriteln(fmt.Sprintf("%s%s\n", streamMessage.String(), streamMessage.ProgressString()))
				lines = orderedmap.New()
				numLayers = 0
			}
		}

		lineBefore = streamMessageStr
	}

	return nil
}

func (r *DefaultResponse) Fwriteln(m interface{}) {

	str := fmt.Sprint(m)
	if r.Transformers != nil {
		for _, t := range r.Transformers {
			str = t(str)
		}
	}

	fmt.Fprintln(r.Writer, str)
}
