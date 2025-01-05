package config

import "log"

var AppConfig *Config

type Config struct {
	DatabaseFile string
	ServerPort   string
	BaseURL      string
}

func LoadConfig() {
	AppConfig = &Config{
		DatabaseFile: "urls.db",
		ServerPort:   "8080",
		BaseURL:      "http://localhost:8080",
	}

	log.Printf("Configuration loaded: %+v\n", AppConfig)
}
