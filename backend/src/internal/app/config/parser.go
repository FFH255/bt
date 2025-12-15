package config

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

var (
	configPathFlagName = "CONFIG_PATH"
	defaultConfigPath  = "./config.yaml"

	envPathFlagName = "ENV_PATH"
	defaultEnvPath  = ".env"
)

func MustParse() *Config {
	configPath := flag.String(configPathFlagName, defaultConfigPath, "Enter config path")
	envPath := flag.String(envPathFlagName, defaultEnvPath, "Enter env path")
	flag.Parse()
	if configPath == nil || len(*configPath) == 0 {
		panic("config path is missing")
	}
	if envPath == nil || len(*envPath) == 0 {
		panic("env path is missing")
	}

	config, err := loadConfig(*configPath, *envPath)
	if err != nil {
		panic(err)
	}

	return config
}

// expandEnvVars replaces ${VAR} with environment variables
func expandEnvVars(content []byte) []byte {
	re := regexp.MustCompile(`\${([^}]+)}`)
	return re.ReplaceAllFunc(content, func(match []byte) []byte {
		envVar := strings.Trim(string(match), "${}")
		if value, exists := os.LookupEnv(envVar); exists {
			return []byte(value)
		}
		return match // return original if env var not found
	})
}

func loadConfig(configPath, envPath string) (*Config, error) {
	if err := godotenv.Load(envPath); err != nil {
		fmt.Println("Error loading .env file")
	}

	// Read the YAML file
	content, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	// Expand environment variables
	expandedContent := expandEnvVars(content)

	// Parse YAML
	var config Config
	if err := yaml.Unmarshal(expandedContent, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
