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

	buildGitContextAuth(io.Writer(&buff))

	expected := `latest: Pulling from library/alpine
<HASH>: Pull complete
Digest: sha256
Status: Downloaded newer image for alpine
sha256: <HASH>
`
	actual := helper.SanitizeDockerOutputForIntegrationTest(&buff)

	assert.Equal(t, expected, actual)
}
