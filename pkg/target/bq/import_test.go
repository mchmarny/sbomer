package bq

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTargetParsing(t *testing.T) {
	t.Parallel()

	targets := []string{
		"bq://project.dataset",
		"bq://project.dataset?location=us",
		"bq://project.dataset?location=us&foo=bar",
	}
	for _, target := range targets {
		r, err := parseTarget(target)
		assert.NoError(t, err)
		assert.NotNil(t, r)
		assert.Equal(t, "project", r.projectID)
		assert.Equal(t, "dataset", r.datasetID)
		assert.Equal(t, "us", r.location)
	}
}
