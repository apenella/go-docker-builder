package main

import (
	"bytes"
	"io"
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

	expected := `3.13: Pulling from alpine
<HASH>: Layer already exists
Digest: sha256
Status: Downloaded newer image for base-registry.go-docker-builder.test
sha256: <HASH>
tag1: digest
latest: digest
`

	actual := helper.SanitizeDockerOutputForIntegrationTest(&buff)

	assert.Equal(t, expected, actual)

}
