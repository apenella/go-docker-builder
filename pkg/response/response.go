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

	// Print method for DefaultResponse requires to clean each entire line before print
	// printTrans := WithTransformers(
	// 	transformer.Prepend(fmt.Sprintf("\r%s", cursor.ClearEntireLine())),
	// )
	// printTrans(r)

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

	//fmt.Fprintf(r.Writer, "\r%s%s \u2500\u2500 %s\n", cursor.ClearEntireLine(), "prefix", m)
	fmt.Fprintln(r.Writer, str)
}

// Write
// func (d *DefaultResponse) Write(w io.Writer, r io.ReadCloser) error {
// 	scanner := bufio.NewScanner(r)
// 	//prefix := d.Prefix

// 	lineBefore := ""
// 	lines := map[string]string{}
// 	numLayers := 0
// 	for scanner.Scan() {
// 		// fmt.Sprintln(scanner.Text())

// 		streamMessage := &types.ResponseBodyStreamMessage{}
// 		line := scanner.Bytes()
// 		err := json.Unmarshal(line, &streamMessage)
// 		if err != nil {
// 			return errors.New("(responser:Response)", fmt.Sprintf("Error unmarshalling line '%s'", string(line)), err)
// 		}

// 		streamMessageStr := streamMessage.String()

// 		if streamMessageStr != lineBefore && streamMessageStr != "" {
// 			if streamMessage.ID != "" {
// 				// lines[streamMessage.ID] = fmt.Sprintf("%s \u2500\u2500  %s %s", prefix, streamMessage.String(), streamMessage.ProgressString())
// 				lines[streamMessage.ID] = fmt.Sprint(streamMessage.String(), streamMessage.ProgressString())

// 				fmt.Fprintf(w, "%s", cursor.MoveUp(numLayers))

// 				for _, line := range lines {
// 					// fmt.Fprintf(w, "\r%s%s\n", cursor.ClearEntireLine(), line)
// 					d.Fwriteln(w, line)
// 				}

// 				numLayers = len(lines)
// 			} else {
// 				fmt.Fprintf(w, "\n")
// 				lines = map[string]string{}
// 				numLayers = 0

// 				d.Fwriteln(w, fmt.Sprint(streamMessage.String(), streamMessage.ProgressString()))
// 				// fmt.Fprintf(w, "\r%s%s \u2500\u2500  %s %s\n", cursor.ClearEntireLine(), prefix, streamMessage.String(), streamMessage.ProgressString())
// 			}
// 		}

// 		lineBefore = streamMessageStr
// 	}
// 	fmt.Fprintf(w, "\n")

// 	return nil
// }
