package transport

import (
	"atc/models"
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
	clientID      string
	clientSecret  string
	redirectURI   string
	url           string
	httpClient    *http.Client
	accessToken   string
	refreshToken  string
	expiresAt     time.Time
	openAIKey     string
	config        *Config
	authenticated bool
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
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			logrus.WithError(err).Error("failed to close secrets file")
			return
		}
	}(file)

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
		clientID:      secrets.Strava.ClientID,
		clientSecret:  secrets.Strava.ClientSecret,
		redirectURI:   config.Server.RedirectURI,
		url:           config.Strava.Url,
		httpClient:    &http.Client{},
		openAIKey:     secrets.OpenAI.APIKey,
		config:        config,
		authenticated: false,
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

// GetAccessToken returns the access token.
func (t *Transport) GetAccessToken() string {
	if t.accessToken == "" {
		logrus.Info("GetAccessToken called with undefined token")
	}
	return t.accessToken
}

// SetAccessToken writes the access token.
func (t *Transport) SetAccessToken(token string) {
	if t.accessToken != "" {
		logrus.Info("SetAccessToken() overwriting existing token")
	} else {
		logrus.Info("SetAccessToken() writing new token")
	}
	t.accessToken = token
}

// GetRefreshToken returns the stored refresh token.
func (t *Transport) GetRefreshToken() string {
	if t.refreshToken == "" {
		logrus.Info("GetRefreshToken called with undefined token")
	}
	return t.refreshToken
}

// SetRefreshToken writes the refresh token
func (t *Transport) SetRefreshToken(token string) {
	if t.refreshToken != "" {
		logrus.Info("SetRefreshToke() overwriting existing token")
	} else {
		logrus.Info("SetRefreshToken() writing new token")
	}
	t.refreshToken = token
}

// GetConfig returns the internal config used by the backend
func (t *Transport) GetConfig() *Config {
	return t.config
}

// IsTokenExpired checks if the current access token is expired.
func (t *Transport) IsTokenExpired() bool {
	return time.Now().After(t.expiresAt)
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
		t.AuthBad()
		return err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		logrus.WithError(err).Error("Failed to decode response body")
		t.AuthBad()
		return err
	}

	if token, ok := result["access_token"].(string); ok {
		t.SetAccessToken(token)
		t.AuthGood()
	} else {
		logrus.Error("Failed to retrieve access token from response body")
		resultJSON, err := json.Marshal(result)
		if err != nil {
			t.AuthBad()
			logrus.WithError(err).Error("Failed to marshal response body")
			return fmt.Errorf("failed to retrieve access token")
		}
		logrus.WithField("result", string(resultJSON)).Error("unexpected response body")
		t.AuthBad()
		return fmt.Errorf("failed to retrieve access token")
	}

	if refreshToken, ok := result["refresh_token"].(string); ok {
		t.refreshToken = refreshToken
	} else {
		logrus.Error("Failed to retrieve refresh token from response body")
		t.AuthBad()
		return fmt.Errorf("failed to retrieve refresh token")
	}

	if expiresIn, ok := result["expires_in"].(float64); ok {
		logrus.Infof("Token expiry in %f seconds", expiresIn)
		t.expiresAt = time.Now().Add(time.Duration(expiresIn) * time.Second)
	} else {
		logrus.Info("Strange or missing expiry data in response")
	}

	return nil
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

	// NOTE: This doesn't actually set cookies because we're not talking to the client
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

// GetAthleteProfile retrieves the athlete's profile from Strava
func (t *Transport) GetAthleteProfile() (*models.Athlete, error) {
	if t.Authenticated() == false {
		logrus.Warn("GetAthleteProfile called but not authenticated")
		return &models.Athlete{}, fmt.Errorf("not authenticated")
	}
	req, err := http.NewRequest("GET", t.url+"/api/v3/athlete", nil)
	if err != nil {
		return &models.Athlete{}, err
	}

	req.Header.Set("Authorization", "Bearer "+t.GetAccessToken())
	resp, err := t.httpClient.Do(req)
	if err != nil {
		logrus.WithError(err).Error("failed to fetch athlete profile")
		return &models.Athlete{}, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logrus.WithError(err).Error("failed to close response body")
			return
		}
	}(resp.Body)

	placeholder := struct {
		ID            int64     `json:"id"`
		Username      *string   `json:"username"`
		ResourceState int       `json:"resource_state"`
		Firstname     string    `json:"firstname"`
		Lastname      string    `json:"lastname"`
		City          string    `json:"city"`
		State         string    `json:"state"`
		Country       string    `json:"country"`
		Sex           string    `json:"sex"`
		Premium       bool      `json:"premium""`
		Summit        bool      `json:"summit"`
		CreatedAt     time.Time `json:"created_at"`
		UpdatedAt     time.Time `json:"updated_at"`
		BadgeTypeID   int       `json:"badge_type_id"`
		ProfileMedium string    `json:"profile_medium"`
		Profile       string    `json:"profile"`
		Friend        *string   `json:"friend"`
		Follower      *string   `json:"follower"`
	}{} // this is a placeholder struct to hold the decoded json data

	if err := json.NewDecoder(resp.Body).Decode(&placeholder); err != nil {
		logrus.WithError(err).Error("failed to decode athlete profile")
		return &models.Athlete{}, err
	}

	th := models.Thresholds{}
	th.Run.ThresholdHR = t.config.Athlete.Run.ThresholdHR
	th.Swim.ThresholdHR = t.config.Athlete.Swim.ThresholdHR
	th.Bike.ThresholdHR = t.config.Athlete.Bike.ThresholdHR

	athlete := models.NewAthlete(
		fmt.Sprintf("%d", placeholder.ID),
		placeholder.Firstname,
		placeholder.Lastname,
		placeholder.Sex,
		&th)

	logrus.Infof("Fetched athlete profile for %s", athlete.FullName())
	return athlete, nil
}

// FetchActivities retrieves activities from Strava API that are of type Swim, Bike, or Run and occurred in the last six weeks.
func (t *Transport) FetchActivities() ([]models.StravaActivity, error) {
	if t.Authenticated() == false {
		logrus.Warn("FetchActivities called but not authenticated")
		return []models.StravaActivity{}, fmt.Errorf("not authenticated")
	}

	token := t.GetAccessToken()

	sixWeeksAgo := time.Now().AddDate(0, 0, -42).Unix()

	sports := []string{"Swim", "Ride", "Run"}

	var allActivities []models.StravaActivity

	// TODO: i feel like these endpoints should be explicitly documented somewhere in code
	//       so that maintaining them or changing them (should strava change their backend
	//       for example) is both easy to do, and easy to audit ("where am i using endpoint xyz?")

	// TODO: i also feel like this is a janky way to create urls for endpoint access.
	//       there's probably a more elegant way to do this but let's do that in the future.

	u, err := url.Parse(fmt.Sprintf("%s/api/v3/athlete/activities", t.url))
	if err != nil {
		logrus.WithError(err).Error("failed to parse URL")
		return []models.StravaActivity{}, err
	}
	// Add query parameters
	params := url.Values{}
	params.Add("access_token", token)
	params.Add("after", fmt.Sprintf("%d", sixWeeksAgo)) // Convert int to string for query params
	params.Add("per_page", "200")

	u.RawQuery = params.Encode()

	resp, err := t.httpClient.Get(u.String())
	if err != nil {
		logrus.WithError(err).Error("failed to fetch activities from Strava")
		return allActivities, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logrus.WithError(err).Error("failed to close response body")
			return
		}
	}(resp.Body)

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
		logrus.WithError(err).Errorf("FetchActivities() failed to decode response body")
		return allActivities, err
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

// AuthGood sets the authenticated flag to true.
func (t *Transport) AuthGood() {
	if t.authenticated == true {
		logrus.Warn("AuthGood called but already authenticated")
		return
	}

	t.authenticated = true
}

// AuthBad sets the authenticated flag to false.
func (t *Transport) AuthBad() {
	// github issue #1, this needs to push out to the reauth flow

	if t.authenticated == false {
		logrus.Warn("AuthBad called but already unauthenticated")
		return
	}

	t.authenticated = false
}

// Authenticated returns the authenticated flag.
func (t *Transport) Authenticated() bool {
	return t.authenticated
}
