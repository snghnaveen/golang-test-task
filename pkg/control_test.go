// Intended for demo only. Actual mocking (docker & AWS) and testing requires more R&D.
package pkg

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewTask(t *testing.T) {
	task := NewTask("some-docker-image",
		"some-bash-command",
		"some-cloud-watch-group",
		"some-cloud-watch-stream",
		"some-secret-access-key",
		"some-access-key", "some-region")
	assert.Equal(t, "some-secret-access-key", task.awsSecretAccessKey, "they should be equal")
	assert.NotEqual(t, "some-secret-access-key", task.awsRegion, "they should not be equal")
}
