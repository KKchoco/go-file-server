package main

import (
	"errors"
	"fmt"

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
	HTTPS   bool
}

type FilesConfig struct {
	FilesPath          string
	MaxUploadSize      int64
	KeyLength          int
	Password           string
	ObfuscateFileNames bool
	AllowedFileTypes   []string
	AdminPassword      string
}

type OtherConfig struct {
	PrettyOutput bool
}

func GetConfig() (Config, error) {
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
		return config, errors.New("error unmarshalling config file")
	}

	return config, nil
}
