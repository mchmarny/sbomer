package sbom

import (
	"os"
	"path"
	"testing"
	"time"

	"github.com/mchmarny/sbomer/pkg/doc"
	"github.com/stretchr/testify/assert"
)

func testDoc(t *testing.T, d *doc.Document, path string) {
	assert.NotNil(t, d, "doc: %s", path)
	assert.NotNil(t, d.ID, "ID: %s", path)
	assert.NotNil(t, d.Subject, "Subject: %s", path)
	assert.NotNil(t, d.SubjectVersion, "SubjectVersion: %s", path)
	assert.NotEmpty(t, d.Format, "Format: %s", path)
	assert.NotEmpty(t, d.FormatVersion, "FormatVersion: %s", path)
	assert.NotEmpty(t, d.Provider, "Provider: %s", path)
	assert.Greater(t, d.Created, time.Now().AddDate(-1, 0, 0).Unix(), "Created: %s", path)
	assert.NotEmpty(t, d.Items, "Items: %s", path)

	for i, p := range d.Items {
		assert.NotEmpty(t, p.ID, "ID: %s[%d] - %s", path, i, p.ToString())
		assert.NotEmpty(t, p.Name, "Name: %s[%d] - %s", path, i, p.ToString())
		assert.NotEmpty(t, p.Version, "Version: %s[%d] - %s", path, i, p.ToString())
		for j, c := range p.Contexts {
			assert.NotEmpty(t, c.Type, "Type: %s[%d][%d] - %s", path, i, j, c.ToString())
			assert.NotEmpty(t, c.Key, "Key: %s[%d][%d] - %s", path, i, j, c.ToString())
			assert.NotEmpty(t, c.Value, "Value: %s[%d][%d] - %s", path, i, j, c.ToString())
		}
	}
}

func TestParsingInvalidDoc(t *testing.T) {
	var in []byte
	_, err := ParseDoc(in)
	assert.Error(t, err)
}

func TestParsing(t *testing.T) {
	files, err := os.ReadDir("../../data")
	assert.NoError(t, err)

	for _, f := range files {
		path := path.Join("../../data", f.Name())
		b, err := os.ReadFile(path)
		assert.NoError(t, err)
		doc, err := ParseDoc(b)
		assert.NoError(t, err, "failed to parse %s", path)
		testDoc(t, doc, path)
	}
}
