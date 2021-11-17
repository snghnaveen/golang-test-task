// Intended for demo only. Actual mocking (docker & AWS) and testing requires more R&D.
package pkg

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfig(t *testing.T) {
	assert.NotNil(t, defaultDockerHost, "they should not be equal")
}
