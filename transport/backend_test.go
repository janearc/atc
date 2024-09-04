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
	assert.NotNil(t, c)

	// NewTransport takes a *config and a string for the secrets file
	backend, err := transport.NewTransport(c, secretsFileName)

	assert.Nil(t, err)
	assert.NotNil(t, backend)
}

func TestLoadSecrets(t *testing.T) {
	//	func LoadSecrets(secretsFileName string) (*Secrets, error) {
	secretsFileName := "config/secrets.yml"
	root := os.Getenv("ATC_ROOT")
	secretsFileName = filepath.Join(root, secretsFileName)

	_, err := transport.LoadSecrets(secretsFileName)

	assert.Nil(t, err)
}

func TestGetAuthURL(t *testing.T) {
	//	func (t *Transport) GetAuthURL() string {
	configFileName := "config/config.yml"
	versionFileName := "config/version.yml"
	secretsFileName := "config/secrets.yml"

	root := os.Getenv("ATC_ROOT")

	// concatenate the working directory with our relative filename
	configFileName = filepath.Join(root, configFileName)
	versionFileName = filepath.Join(root, versionFileName)
	secretsFileName = filepath.Join(root, secretsFileName)

	c, err := transport.LoadConfig(configFileName, versionFileName)

	assert.Nil(t, err)

	// NewTransport takes a *config and a string for the secrets file
	backend, err := transport.NewTransport(c, secretsFileName)

	assert.Nil(t, err)
	assert.NotNil(t, backend)

	// this is a test so we're not going to worry about the error
	url := backend.GetAuthURL()

	// not super clear on how to determine whether the url is valid but i think this works for now
	assert.NotNil(t, url)
}

func Test_GetConfig(t *testing.T) {
	// func (t *Transport) GetConfig() *Config {

	configFileName := "config/config.yml"
	versionFileName := "config/version.yml"
	secretsFileName := "config/secrets.yml"

	root := os.Getenv("ATC_ROOT")

	// concatenate the working directory with our relative filename
	configFileName = filepath.Join(root, configFileName)
	versionFileName = filepath.Join(root, versionFileName)
	secretsFileName = filepath.Join(root, secretsFileName)

	c, err := transport.LoadConfig(configFileName, versionFileName)

	assert.Nil(t, err)
	assert.NotNil(t, c)

	// NewTransport takes a *config and a string for the secrets file
	backend, err := transport.NewTransport(c, secretsFileName)

	assert.Nil(t, err)
	assert.NotNil(t, backend)

	config := backend.GetConfig()

	// this is just an accessor, so we don't need to run the tests that are in the config test package
	assert.NotNil(t, config)
}

// requires mocking http & strava:
// func (t *Transport) ExchangeCodeForToken(code string) error {
// func (t *Transport) GetAccessToken() string {
// func (t *Transport) GetRefreshToken() string {
// func (t *Transport) IsTokenExpired() bool {
// func (t *Transport) RefreshAccessToken(refreshToken string) (string, error) {
// func (t *Transport) FetchActivities(token string) ([]models.StravaActivity, error) {

// probably don't need to test this but maybe it makes sense for documentation
// func (t *Transport) ExampleRequest(endpoint string) ([]byte, error) {

// requires mocking http & openai
// func (t *Transport) OpenAIRequest(prompt string) (string, error) {
