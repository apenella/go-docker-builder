package main

import (
	"bytes"
	"io"
	"strings"
	"testing"

	helper "github.com/apenella/go-docker-builder/internal/helpers"
	"github.com/stretchr/testify/assert"
)

func TestBuildPathContext(t *testing.T) {

	var buff bytes.Buffer

	err := buildPathContext(io.Writer(&buff))
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
