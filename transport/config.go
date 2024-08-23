package transport

import (
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	Server struct {
		Port        int    `yaml:"port"`
		RedirectURI string `yaml:"redirect_uri"`
	} `yaml:"server"`

	Strava struct {
		Url string `yaml:"url"`
	} `yaml:"strava"`
}

func LoadConfig() (*Config, error) {
	// NOTE: this is a weird path because of docker
	file, err := os.Open("/app/config/config.yml")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config := &Config{}
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(config); err != nil {
		return nil, err
	}

	return config, nil
}
