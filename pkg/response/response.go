package response

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"

	errors "github.com/apenella/go-common-utils/error"
	"github.com/apenella/go-docker-builder/pkg/types"
)

type DefaultResponse struct {
	Prefix string
}

// Response
func (d *DefaultResponse) Write(w io.Writer, r io.ReadCloser) error {
	scanner := bufio.NewScanner(r)
	prefix := d.Prefix

	lineBefore := ""
	for scanner.Scan() {
		streamMessage := &types.ResponseBodyStreamMessage{}
		line := scanner.Bytes()
		err := json.Unmarshal(line, &streamMessage)
		if err != nil {
			return errors.New("(responser:Response)", fmt.Sprintf("Error unmarshalling line '%s'", string(line)), err)
		}

		streamMessageStr := streamMessage.String()

		if streamMessageStr != lineBefore {
			fmt.Fprintf(w, "\n%s \u2500\u2500  %s %s", prefix, streamMessage.String(), streamMessage.ProgressString())
		} else {
			fmt.Fprintf(w, "\r%s \u2500\u2500  %s %s", prefix, streamMessage.String(), streamMessage.ProgressString())
		}

		lineBefore = streamMessageStr
	}
	// print empty line at the end
	fmt.Println()

	return nil
}
