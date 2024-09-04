package transport_test

import (
	"atc/transport"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// this is for NewTransport, which is different than NewService, which is a wrapper
// around the transport backend(s)
func TestNewTransport(t *testing.T) {
	//	func NewTransport(configFileName, versionFileName, secretsFileName string) (*Transport, error) {
	configFileName := "config/config.yml"
	versionFileName := "config/version.yml"
	secretsFileName := "config/secrets.yml"

	root := os.Getenv("ATC_ROOT")

	// concatenate the working directory with our relative filename
	configFileName = filepath.Join(root, configFileName)
	versionFileName = filepath.Join(root, versionFileName)
	secretsFileName = filepath.Join(root, secretsFileName)

	// we're testing the transport here so we're not going to go through all the config tests
	c, err := transport.LoadConfig(configFileName, versionFileName)

	assert.Nil(t, err)

	// NewTransport takes a *config and a string for the secrets file
	_, err = transport.NewTransport(c, secretsFileName)

	assert.Nil(t, err)
}

func TestLoadSecrets(t *testing.T) {
	//	func LoadSecrets(secretsFileName string) (*Secrets, error) {
	secretsFileName := "config/secrets.yml"
	root := os.Getenv("ATC_ROOT")
	secretsFileName = filepath.Join(root, secretsFileName)

	_, err := transport.LoadSecrets(secretsFileName)

	assert.Nil(t, err)
}
