package main

import (
	"bytes"
	"io"
	"testing"

	helper "github.com/apenella/go-docker-builder/internal/helpers"
	"github.com/stretchr/testify/assert"
)

// TestBuildAndPush is considered and integration test and uses build-and-push to proceed with the test. It requires a test integration environment to run and must be run by `make test`
func TestBuildAndPushJoinContext(t *testing.T) {

	var buff bytes.Buffer

	buildAndPushJoinContext(io.Writer(&buff))

	expected := `1.15-alpine: Pulling from golang
<HASH>: Pushed
Digest: sha256
Status: Downloaded newer image for base-registry.go-docker-builder.test
sha256: <HASH>
latest: digest
`
	actual := helper.SanitizeDockerOutputForIntegrationTest(&buff)

	assert.Equal(t, expected, actual)
}
