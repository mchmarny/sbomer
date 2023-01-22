package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	testDir = "./"
)

func TestConfig(t *testing.T) {
	c1, err := ReadOrCreate(testDir)
	assert.NoError(t, err)
	assert.NotNil(t, c1)

	c1.LastExec = time.Now()
	c1.LastVersion = "v0.0.1-test"

	err = Save(testDir, c1)
	assert.NoError(t, err)

	c2, err := ReadOrCreate(testDir)
	assert.NoError(t, err)
	assert.NotNil(t, c2)
	assert.Equal(t, c1.LastExec.Format(time.RFC3339), c2.LastExec.Format(time.RFC3339))
	assert.Equal(t, c1.LastVersion, c2.LastVersion)
}
