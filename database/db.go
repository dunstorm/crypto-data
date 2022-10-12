package database

import (
	"encoding/json"
	"os"
	"path"
)

type Config struct {
	// api key
	BinanceAPIKey string

	// api secret
	BinanceAPISecret string
}

func getConfigPath() string {
	// ~/.crypto-data/db.json
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return path.Join(homeDir, ".crypto-data", "db.json")
}

func getConfig() (Config, error) {
	configPath := getConfigPath()
	// if file does not exist create it
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// create full path
		os.MkdirAll(path.Dir(configPath), 0755)
		_, err := os.Create(configPath)
		if err != nil {
			return Config{}, err
		}
	}
	data, err := os.ReadFile(configPath)
	var config Config
	if err == nil {
		json.Unmarshal(data, &config)
	}
	return config, err
}

func updateConfig(data Config) {
	bytes, err := json.Marshal(data)
	configPath := getConfigPath()
	if err == nil {
		os.WriteFile(configPath, bytes, 0644)
	} else {
		panic(err)
	}
}

func FindOrCreateConfig() Config {
	config, err := getConfig()
	if err != nil {
		config = Config{
			BinanceAPIKey:    "",
			BinanceAPISecret: "",
		}
		updateConfig(config)
	}
	return config
}

func UpdateConfig(data Config) {
	updateConfig(data)
}
