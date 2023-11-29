package main

import (
	"bytes"
	"io"
	"testing"

	helper "github.com/apenella/go-docker-builder/internal/helpers"
	"github.com/stretchr/testify/assert"
)

// TestBuildAndPush is considered and integration test and uses build-and-push to proceed with the test. It requires a test integration environment to run and must be run by `make test`
func TestBuildGitContextAuth(t *testing.T) {

	var buff bytes.Buffer

	buildGitContextAuth(io.Writer(&buff))

	expected := `latest: digest
<HASH>: Layer already exists
Digest: sha256
Status: Downloaded newer image for alpine
sha256: <HASH>
custom: digest
`

	actual := helper.SanitizeDockerOutputForIntegrationTest(&buff)

	assert.Equal(t, expected, actual)
}
