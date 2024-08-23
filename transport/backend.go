package transport

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"gopkg.in/yaml.v2"
)

type Secrets struct {
	Strava struct {
		ClientID     string `yaml:"client_id"`
		ClientSecret string `yaml:"client_secret"`
	} `yaml:"strava"`
}

type Transport struct {
	clientID     string
	clientSecret string
	redirectURI  string
	url          string
	httpClient   *http.Client
	accessToken  string
}

// LoadSecrets reads the secrets.yml file and returns a Secrets struct.
func LoadSecrets() (*Secrets, error) {
	file, err := os.Open("config/secrets.yml")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	secrets := &Secrets{}
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(secrets); err != nil {
		return nil, err
	}

	return secrets, nil
}

// NewTransport initializes the Transport with secrets and config.
func NewTransport(config *Config) (*Transport, error) {
	secrets, err := LoadSecrets()
	if err != nil {
		return nil, err
	}

	return &Transport{
		clientID:     secrets.Strava.ClientID,
		clientSecret: secrets.Strava.ClientSecret,
		redirectURI:  config.Server.RedirectURI,
		url:          config.Strava.Url,
		httpClient:   &http.Client{},
	}, nil
}

// GetAuthURL generates the Strava OAuth URL for authentication.
func (t *Transport) GetAuthURL() string {
	return fmt.Sprintf(
		"%s/oauth/authorize?client_id=%s&redirect_uri=%s&response_type=code&scope=read,activity:read",
		t.url,
		t.clientID,
		url.QueryEscape(t.redirectURI),
	)
}

// ExchangeCodeForToken exchanges the authorization code for an access token.
func (t *Transport) ExchangeCodeForToken(code string) error {
	reqURL := fmt.Sprintf("%s/oauth/token", t.url)
	data := url.Values{
		"client_id":     {t.clientID},
		"client_secret": {t.clientSecret},
		"code":          {code},
		"grant_type":    {"authorization_code"},
		"redirect_uri":  {t.redirectURI},
	}

	resp, err := t.httpClient.PostForm(reqURL, data)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	if token, ok := result["access_token"].(string); ok {
		t.accessToken = token
		return nil
	}

	return fmt.Errorf("failed to retrieve access token")
}

// GetAccessToken returns the access token.
func (t *Transport) GetAccessToken() string {
	return t.accessToken
}

// ExampleRequest makes an authenticated request to Strava API.
func (t *Transport) ExampleRequest(endpoint string) ([]byte, error) {
	req, err := http.NewRequest("GET", t.url+"/api/v3"+endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+t.accessToken)
	resp, err := t.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
