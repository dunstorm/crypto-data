package database

import (
	"encoding/json"
	"os"
)

type Config struct {
	// api key
	BinanceAPIKey string

	// api secret
	BinanceAPISecret string
}

func getConfig() (Config, error) {
	data, err := os.ReadFile("database/db.json")
	var config Config
	if err == nil {
		json.Unmarshal(data, &config)
	}
	return config, err
}

func updateConfig(data Config) {
	bytes, err := json.Marshal(data)
	if err == nil {
		os.WriteFile("database/db.json", bytes, 0644)
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
