package transport_test

import (
	"atc/transport"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadSecrets(t *testing.T) {
	//	func LoadSecrets(secretsFileName string) (*Secrets, error) {
	secretsFileName := "config/secrets.yml"
	root := os.Getenv("ATC_ROOT")
	secretsFileName = filepath.Join(root, secretsFileName)

	_, err := transport.LoadSecrets(secretsFileName)

	assert.Nil(t, err)
}
