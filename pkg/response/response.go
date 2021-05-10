package response

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"

	"github.com/ahmetb/go-cursor"
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
	lines := map[string]string{}
	numLayers := 0
	for scanner.Scan() {
		// fmt.Sprintln(scanner.Text())

		streamMessage := &types.ResponseBodyStreamMessage{}
		line := scanner.Bytes()
		err := json.Unmarshal(line, &streamMessage)
		if err != nil {
			return errors.New("(responser:Response)", fmt.Sprintf("Error unmarshalling line '%s'", string(line)), err)
		}

		streamMessageStr := streamMessage.String()

		if streamMessageStr != lineBefore && streamMessageStr != "" {
			if streamMessage.ID != "" {
				lines[streamMessage.ID] = fmt.Sprintf("%s \u2500\u2500  %s %s", prefix, streamMessage.String(), streamMessage.ProgressString())

				fmt.Fprintf(w, "%s", cursor.MoveUp(numLayers))

				for _, line := range lines {
					fmt.Fprintf(w, "\r%s%s\n", cursor.ClearEntireLine(), line)
				}

				numLayers = len(lines)
			} else {
				fmt.Fprintf(w, "\n")
				lines = map[string]string{}
				numLayers = 0

				fmt.Fprintf(w, "\r%s%s \u2500\u2500  %s %s\n", cursor.ClearEntireLine(), prefix, streamMessage.String(), streamMessage.ProgressString())
			}
		}

		lineBefore = streamMessageStr
	}
	fmt.Fprintf(w, "\n")

	return nil
}
