package transport_test

import (
	"atc/transport"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	// fire up that constructor and read some yaml
	config, err := transport.LoadConfig()
	assert.Nil(t, err)
	assert.NotNil(t, config)
}
