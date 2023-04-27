package utils

import (
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestSplitContainerName(t *testing.T) {
	containerName := "peer5.sorg.nindindo.com"
	name, domain := SplitContainerName(containerName)
	assert.Equal(t, name, "peer5")
	assert.Equal(t, domain, "sorg.nindindo.com")
}
