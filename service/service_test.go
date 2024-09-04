package service_test

import (
	"atc/service"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewService(t *testing.T) {
	// Call the constructor
	// for testing
	configFileName := "config/config.yml"
	versionFileName := "config/version.yml"
	secretsFileName := "config/secrets.yml"

	root := os.Getenv("ATC_ROOT")

	// concatenate the working directory with our relative filename
	configFileName = filepath.Join(root, configFileName)
	versionFileName = filepath.Join(root, versionFileName)
	secretsFileName = filepath.Join(root, secretsFileName)

	s := service.NewService(configFileName, versionFileName, secretsFileName)

	// Verify that the service is not nil
	assert.NotNil(t, s)

	// Verify that the service's fields are initialized correctly
	assert.NotNil(t, s.Config)
	assert.NotNil(t, s.Log)
	assert.NotNil(t, s.Backend)
	assert.NotNil(t, s.Web.Handle)
}
