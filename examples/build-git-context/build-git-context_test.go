package main

import (
	"bytes"
	"io"
	"testing"

	helper "github.com/apenella/go-docker-builder/internal/helpers"
	"github.com/stretchr/testify/assert"
)

// TestBuildAndPush is considered and integration test and uses build-and-push to proceed with the test. It requires a test integration environment to run and must be run by `make test`
func TestBuildGitContext(t *testing.T) {

	var buff bytes.Buffer

	err := buildGitContext(io.Writer(&buff))
	if err != nil {
		t.Error(err.Error())
	}

	expected := `sha256: <HASH>
<HASH>: Layer already exists
a-tag: digest
b-tag: digest
z-tag: digest
latest: digest
`

	actual := helper.SanitizeDockerOutputForIntegrationTest(&buff)

	assert.Equal(t, expected, actual)
}
