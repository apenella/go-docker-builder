package response

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/apenella/go-docker-builder/pkg/types"
)

type ResponseHandler struct {
	Reader io.ReadCloser
	Writer io.Writer
	Prefix string
}

// Response
func (r *ResponseHandler) Run() error {
	scanner := bufio.NewScanner(r.Reader)
	prefix := r.Prefix

	for scanner.Scan() {
		streamMessage := &types.BuildResponseBodyStreamMessage{}
		line := scanner.Bytes()
		err := json.Unmarshal(line, &streamMessage)
		if err != nil {
			return errors.New("(responser:Response) Error unmarshalling line '" + string(line) + "' " + err.Error())
		}

		fmt.Fprintf(r.Writer, "%s \u2500\u2500  %s\n", prefix, streamMessage.String())
	}

	return nil
}

// SetReader
func (r *ResponseHandler) SetReader(reader io.ReadCloser) {
	r.Reader = reader
}

// SetWriter
func (r *ResponseHandler) SetWriter(writer io.Writer) {
	r.Writer = writer
}
