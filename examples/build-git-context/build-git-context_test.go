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

	buildGitContext(io.Writer(&buff))

	expected := `3.9: Pulling from library/alpine
<HASH>: Layer already exists
Digest: sha256
Status: Downloaded newer image for alpine
sha256: <HASH>
a-tag: digest
b-tag: digest
z-tag: digest
latest: digest
`
	actual := helper.SanitizeDockerOutputForIntegrationTest(&buff)

	assert.Equal(t, expected, actual)
}
