package file

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileParsing(t *testing.T) {
	ctx := context.TODO()
	locations := []string{
		"../../data/swaggerapi_petstore_1.0.6-syft-spdx.json",
		"https://raw.githubusercontent.com/chainguard-dev/bom-shelter/main/in-the-wild/spdx/powershell-2.2.6.spdx.json",
	}
	for _, location := range locations {
		content, err := GetContent(ctx, location)
		assert.NoError(t, err)
		assert.NotEmpty(t, content)
	}
}
