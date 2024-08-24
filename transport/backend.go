package transport

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
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
	OpenAI struct {
		APIKey string `yaml:"api_key"`
	} `yaml:"openai"`
}

type Transport struct {
	clientID     string
	clientSecret string
	redirectURI  string
	url          string
	httpClient   *http.Client
	accessToken  string
	openAIKey    string
}

// LoadSecrets reads the secrets.yml file and returns a Secrets struct.
func LoadSecrets() (*Secrets, error) {
	file, err := os.Open("/app/config/secrets.yml")
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
		openAIKey:    secrets.OpenAI.APIKey,
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

func (t *Transport) OpenAIRequest(prompt string) (string, error) {
	requestBody, err := json.Marshal(map[string]interface{}{
		"model": "gpt-3.5-turbo",
		"messages": []map[string]string{
			{"role": "system", "content": "You are a helpful assistant."},
			{"role": "user", "content": prompt},
		},
	})
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+t.openAIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := t.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	// Log the full response body using logrus
	logrus.WithField("response_body", string(body)).Info("Full OpenAI Response")

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if choices, ok := response["choices"].([]interface{}); ok && len(choices) > 0 {
		if message, ok := choices[0].(map[string]interface{})["message"].(map[string]interface{}); ok {
			return message["content"].(string), nil
		}
	}

	return "", fmt.Errorf("no valid response from OpenAI")
}
