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

	Build struct {
		BuildDate string `yaml:"build_date"`
		Build     string `yaml:"build"`
	} `yaml:"version"`
}

// LoadConfig reads the config.yml file and returns a Config struct.
func LoadConfig(configFileName string, versionFileName string) (*Config, error) {
	if configFileName == "" {
		// this is an absolute path but it's inside the docker container
		// we assume that if we have not been called with a filename, that
		// we aren't running in a container or we're running locally or
		// something.
		configFileName = "/app/config/config.yml"
	}

	file, err := os.Open(configFileName)
	if err != nil {
		logrus.WithError(err).Fatalf("Failed to open config file %s", configFileName)
		return nil, err
	}
	defer file.Close()

	// Create a new Config struct
	var config Config

	// Decode the config file
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		logrus.WithError(err).Fatal("Failed to decode config file")
		return nil, err
	}

	logrus.Info("Successfully loaded configuration")

	// same as above, but for the version file
	if versionFileName == "" {
		versionFileName = "/app/config/version.yml"
	}
	vf, err := os.Open(versionFileName)
	if err != nil {
		logrus.WithError(err).Fatalf("Failed to open version file %s", versionFileName)
		return nil, err
	}
	defer vf.Close()

	// Decode the version file
	vfDecoder := yaml.NewDecoder(vf)
	if err := vfDecoder.Decode(&config); err != nil {
		logrus.WithError(err).Fatalf("Failed to decode version file %s", versionFileName)
		return nil, err
	}

	return &config, nil
}
