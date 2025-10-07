package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

// This structure maps the YAML code to a go structure .
type Config struct {
	LoadBalancer struct {
		Service string   `yaml:"service"`
		Port    string   `yaml:"port"`
		Servers []string `yaml:"servers"`
	} `yaml:"loadbalancer"`
}

// This function reads the YAML file and sends the readed data back to the `cfg` variable
func LoadConfig(filePath string) (*Config, error) {
	configFile, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(configFile, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// This function prints the loaded configuration details to verify proper execution.
func CheckAndLogConfig(cfg *Config) {
	log.Printf("Service: %s", cfg.LoadBalancer.Service)
	log.Printf("Port: %s", cfg.LoadBalancer.Port)
	log.Printf("Servers: %v", cfg.LoadBalancer.Servers)
	}