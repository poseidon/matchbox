package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequirementMatches(t *testing.T) {
	// requirements
	reqs := map[string]string{
		"region": "Central US",
		"zone":   "us-central1-a",
		"lot":    "42",
	}
	attrs := map[string]string{
		"uuid": "16e7d8a7-bfa9-428b-9117-363341bb330b",
	}
	// labels
	labels := map[string]string{
		"region": "Central US",
		"zone":   "us-central1-a",
		"lot":    "42",
	}
	query := map[string]string{
		"uuid": "16e7d8a7-bfa9-428b-9117-363341bb330b",
	}
	lacking := map[string]string{
		"region": "Central US",
	}

	cases := []struct {
		reqs     map[string]string
		labels   map[string]string
		expected bool
	}{
		{reqs, labels, true},
		{attrs, query, true},
		{reqs, lacking, false},
		// zero requirements match any label set
		{map[string]string{}, labels, true},
	}
	for _, c := range cases {
		r := RequirementSet(c.reqs)
		assert.Equal(t, c.expected, r.Matches(c.labels))
	}
}
