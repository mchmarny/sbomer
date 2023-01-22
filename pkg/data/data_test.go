package data

import (
	"os"
	"testing"

	"github.com/mchmarny/sbomer/pkg/sbom"
	"github.com/stretchr/testify/assert"
)

const (
	testDir = "./"
)

func deleteDB() {
	os.Remove("./data.db") // nolint: errcheck
}

func TestData(t *testing.T) {
	deleteDB()

	path := "../../data/redis:7.0.8-syft-spdx.json"
	b, err := os.ReadFile(path)
	assert.NoError(t, err)
	doc, err := sbom.ParseDoc(b)
	assert.NoError(t, err)

	err = Init(testDir)
	assert.NoError(t, err)

	err = Save(doc)
	assert.NoError(t, err)

	d, err := Get(doc.ID)
	assert.NoError(t, err)
	assert.NotNil(t, d)

	subjects, err := Query()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(subjects))

	err = Close()
	assert.NoError(t, err)

	deleteDB()
}
