package transport_test

import (
	"atc/transport"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	root := os.Getenv("ATC_ROOT")

	configFileName := filepath.Join(root, "config/config.yml")
	versionFileName := filepath.Join(root, "config/version.yml")

	// fire up that constructor and read some yaml
	config, err := transport.LoadConfig(configFileName, versionFileName)

	assert.Nil(t, err)
	assert.NotNil(t, config)
}
