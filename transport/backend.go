package transport

import (
	"atc/models"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

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
	refreshToken string
	expiresAt    time.Time
	openAIKey    string
	config       *Config
}

// LoadSecrets reads the secrets.yml file and returns a Secrets struct.
func LoadSecrets(secretsFileName string) (*Secrets, error) {
	// ordinarily we're running in a container somewhere out in the cosmos but for
	// testing and running locally we want to be able to pass in a specific filename
	if secretsFileName == "" {
		secretsFileName = "/app/config/secrets.yml"
	}
	file, err := os.Open(secretsFileName)
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
func NewTransport(config *Config, secretsFile string) (*Transport, error) {
	secrets, err := LoadSecrets(secretsFile)
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
		config:       config,
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

// ExchangeCodeForToken exchanges the authorization code for an access token and stores the refresh token and expiration time.
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
		logrus.WithError(err).Error("Failed to post form to exchange code for token")
		return err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		logrus.WithError(err).Error("Failed to decode response body")
		return err
	}

	if token, ok := result["access_token"].(string); ok {
		t.accessToken = token
	} else {
		logrus.Error("Failed to retrieve access token from response body")
		resultJSON, err := json.Marshal(result)
		if err != nil {
			logrus.WithError(err).Error("Failed to marshal response body")
			return fmt.Errorf("failed to retrieve access token")
		}
		logrus.WithField("result", string(resultJSON)).Error("unexpected response body")
		return fmt.Errorf("failed to retrieve access token")
	}

	if refreshToken, ok := result["refresh_token"].(string); ok {
		t.refreshToken = refreshToken
	} else {
		logrus.Error("Failed to retrieve refresh token from response body")
		return fmt.Errorf("failed to retrieve refresh token")
	}

	if expiresIn, ok := result["expires_in"].(float64); ok {
		logrus.Infof("Token expiry in %f seconds", expiresIn)
		t.expiresAt = time.Now().Add(time.Duration(expiresIn) * time.Second)
	}

	return nil
}

// GetAccessToken returns the access token.
func (t *Transport) GetAccessToken() string {
	return t.accessToken
}

// GetRefreshToken returns the stored refresh token.
func (t *Transport) GetRefreshToken() string {
	return t.refreshToken
}

// GetConfig returns the internal config used by the backend
func (t *Transport) GetConfig() *Config {
	return t.config
}

// IsTokenExpired checks if the current access token is expired.
func (t *Transport) IsTokenExpired() bool {
	return time.Now().After(t.expiresAt)
}

// RefreshAccessToken uses the refresh token to obtain a new access token.
func (t *Transport) RefreshAccessToken(refreshToken string) (string, error) {
	reqURL := fmt.Sprintf("%s/oauth/token", t.url)
	data := url.Values{
		"client_id":     {t.clientID},
		"client_secret": {t.clientSecret},
		"grant_type":    {"refresh_token"},
		"refresh_token": {refreshToken},
	}

	resp, err := t.httpClient.PostForm(reqURL, data)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if newAccessToken, ok := result["access_token"].(string); ok {
		t.accessToken = newAccessToken
		if expiresIn, ok := result["expires_in"].(float64); ok {
			t.expiresAt = time.Now().Add(time.Duration(expiresIn) * time.Second)
		}
		return newAccessToken, nil
	}

	return "", fmt.Errorf("failed to refresh access token")
}

// ExampleRequest makes an authenticated request to Strava API. This method is not
// actually used by the backend, but it's preserved for documentation's sake. please
// don't remove this.
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

// FetchActivities retrieves activities from Strava API that are of type Swim, Bike, or Run and occurred in the last six weeks.
func (t *Transport) FetchActivities(token string) ([]models.StravaActivity, error) {
	sixWeeksAgo := time.Now().AddDate(0, 0, -42).Unix()

	sports := []string{"Swim", "Ride", "Run"}

	var allActivities []models.StravaActivity

	// TODO: i feel like these endpoints should be explicitly documented somewhere in code
	//       so that maintaining them or changing them (should strava change their backend
	//       for example) is both easy to do, and easy to audit ("where am i using endpoint xyz?")

	// TODO: i also feel like this is a janky way to create urls for endpoint access.
	//       there's probably a more elegant way to do this but let's do that in the future.
	url := fmt.Sprintf("%s/api/v3/athlete/activities?access_token=%s&after=%d&per_page=200", t.url, token, sixWeeksAgo)
	resp, err := t.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch activities from Strava: %w", err)
	}
	defer resp.Body.Close()

	// Temporary structure to hold the raw JSON data
	var tempActivities []struct {
		Id                 int64     `json:"id"`
		Name               string    `json:"name"`
		Distance           float64   `json:"distance"`
		MovingTime         int       `json:"moving_time"`
		ElapsedTime        int       `json:"elapsed_time"`
		TotalElevationGain float64   `json:"total_elevation_gain"`
		Type               string    `json:"type"`
		StartDate          time.Time `json:"start_date"`
		Calories           int       `json:"calories"`
		AverageHeartRate   float64   `json:"average_heartrate"`
		MaxHeartRate       float64   `json:"max_heartrate"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tempActivities); err != nil {
		return nil, fmt.Errorf("failed to decode activities: %w", err)
	}

	// map the decoded json data to StravaActivity objects using the constructor
	for _, ta := range tempActivities {
		// this is just a really ugly grep
		for _, sport := range sports {
			if ta.Type == sport {
				// returns a pointer
				activity := models.NewStravaActivity(
					ta.Id,
					ta.Name,
					ta.Distance,
					ta.MovingTime,
					ta.ElapsedTime,
					ta.TotalElevationGain,
					ta.Type,
					ta.StartDate,
					ta.Calories,
					ta.AverageHeartRate,
					ta.MaxHeartRate,
				)
				allActivities = append(allActivities, activity)
			}
		}
	}

	return allActivities, nil
}

// OpenAIRequest sends a request to OpenAI's API and logs the full response using logrus.
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
