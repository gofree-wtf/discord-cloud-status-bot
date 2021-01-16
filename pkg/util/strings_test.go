package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSubstringAfter(t *testing.T) {
	commandPrefix := "!status"
	userMsg := commandPrefix + " health"
	command := SubstringAfter(userMsg, commandPrefix+" ")
	assert.Equal(t, "health", command)
}
