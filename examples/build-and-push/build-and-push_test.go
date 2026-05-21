package main

import (
	"bytes"
	"io"
	"strings"
	"testing"

	helper "github.com/apenella/go-docker-builder/internal/helpers"
	"github.com/stretchr/testify/assert"
)

// TestBuildAndPush is considered and integration test and uses build-and-push to proceed with the test. It requires a test integration environment to run and must be run by `make test`
func TestBuildAndPush(t *testing.T) {

	var buff bytes.Buffer

	err := buildAndPush(io.Writer(&buff))
	if err != nil {
		t.Error(err.Error())
	}

	expected := `sha256: <HASH>
<HASH>: Layer already exists
tag1: digest
latest: digest
`

	actual := helper.SanitizeDockerOutputForIntegrationTest(&buff)

	expectedLines := strings.Split(strings.TrimSpace(expected), "\n")
	actualLines := strings.Split(strings.TrimSpace(actual), "\n")

	assert.ElementsMatch(t, actualLines, expectedLines)

}
