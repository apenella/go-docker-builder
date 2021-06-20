package main

import (
	"bytes"
	"io"
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

	expected := `3.13: Pulling from alpine
<HASH>: Pull complete
Digest: sha256
Status: Downloaded newer image for base-registry.go-docker-builder.test
sha256: <HASH>
`
	actual := helper.SanitizeDockerOutputForIntegrationTest(&buff)

	assert.Equal(t, expected, actual)
}
