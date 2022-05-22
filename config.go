package main

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Config struct {
	Server ServerConfig
	Files  FilesConfig
	Other  OtherConfig
}

type ServerConfig struct {
	Address string
	Port    string
	Domain  string
}

type FilesConfig struct {
	FilesPath        string
	MaxUploadSize    int64
	KeyLength        int
	AllowedFileTypes []string
}

type OtherConfig struct {
	PrettyOutput bool
}

func GetConfig() Config {
	var config Config

	// Set config.yml information
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yml")

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}

	// Provide default values
	viper.SetDefault("other.prettyOutput", false)
	viper.SetDefault("files.keyLength", 7)
	viper.SetDefault("files.maxUploadSize", 20)

	// Assert config file to config variable
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal().Err(err).Msg("Unable to convert config.yml into struct")
	}

	return config
}

func isValidConfig(config Config) bool {
	errCount := 0

	if len(config.Server.Address) == 0 {
		log.Error().Msg("server.address is required")
		errCount += 1
	}
	if len(config.Server.Port) == 0 {
		log.Error().Msg("server.port is required")
		errCount += 1
	}
	if len(config.Files.FilesPath) == 0 {
		log.Error().Msg("files.filesPath is required")
		errCount += 1
	}

	if errCount > 0 {
		log.Fatal().Msgf("%v errors occured while parsing config.yml", errCount)
		return false
	} else {
		return true
	}
}

func DisplayConfig(config Config) {
	log.Info().Msg("Loaded config!")
	log.Info().Msgf("%+v", config)
}
