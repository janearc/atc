package transport

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"os"
)

// Config struct to hold the configuration values from config.yml
type Config struct {
	Server struct {
		Port        int    `yaml:"port"`
		RedirectURI string `yaml:"redirect_uri"`
	} `yaml:"server"`

	Strava struct {
		Url string `yaml:"url"`
	} `yaml:"strava"`

	Athlete struct {
		Run struct {
			ThresholdHR float64 `yaml:"threshold_hr"`
		} `yaml:"run"`
		Swim struct {
			ThresholdHR float64 `yaml:"threshold_hr"`
		} `yaml:"swim"`
		Bike struct {
			ThresholdHR float64 `yaml:"threshold_hr"`
		} `yaml:"bike"`
	} `yaml:"athlete"`
}

// LoadConfig reads the config.yml file and returns a Config struct.
func LoadConfig() (*Config, error) {
	file, err := os.Open("/app/config/config.yml")
	if err != nil {
		logrus.WithError(err).Fatal("Failed to open config file")
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		logrus.WithError(err).Fatal("Failed to decode config file")
		return nil, err
	}

	logrus.Info("Successfully loaded configuration")
	return &config, nil
}
