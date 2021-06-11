package helper

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/acarl005/stripansi"
	"github.com/apenella/go-docker-builder/pkg/types"
	orderedmap "github.com/wk8/go-ordered-map"
)

// SanitizeDockerOutputForIntegrationTest only keeps those output lines which contains layer status details and redacts any hashes
func SanitizeDockerOutputForIntegrationTest(r io.Reader) string {
	var res bytes.Buffer

	scanner := bufio.NewScanner(r)
	layersMap := orderedmap.New()
	shahashregexp, _ := regexp.Compile("[0-9a-z]{12,64}")

	for scanner.Scan() {
		text := stripansi.Strip(scanner.Text())
		if text != "" {
			lineSplited := strings.Split(text, types.LayerMessagePrefix)
			if len(lineSplited) > 1 {
				status := strings.Split(lineSplited[1], ":")
				if len(status) > 1 {

					layerID := strings.TrimSpace(status[0])
					isHash := shahashregexp.MatchString(layerID)
					if isHash {
						layerID = "<HASH>"
					}
					message := strings.TrimSpace(status[1])
					isHash = shahashregexp.MatchString(message)
					if isHash {
						message = "<HASH>"
					}

					layersMap.Set(layerID, message)
				}
			}
		}
	}

	for line := layersMap.Oldest(); line != nil; line = line.Next() {
		res.WriteString(fmt.Sprintf("%s: %s\n", line.Key, line.Value))
	}

	return res.String()
}
