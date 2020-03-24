package response

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/apenella/go-docker-builder/pkg/types"
)

type DefaultResponse struct {
	Prefix string
}

// Response
func (d *DefaultResponse) Write(w io.Writer, r io.ReadCloser) error {
	scanner := bufio.NewScanner(r)
	prefix := d.Prefix

	for scanner.Scan() {
		streamMessage := &types.BuildResponseBodyStreamMessage{}
		line := scanner.Bytes()
		err := json.Unmarshal(line, &streamMessage)
		if err != nil {
			return errors.New("(responser:Response) Error unmarshalling line '" + string(line) + "' " + err.Error())
		}

		fmt.Fprintf(w, "%s \u2500\u2500  %s\n", prefix, streamMessage.String())
	}

	return nil
}
