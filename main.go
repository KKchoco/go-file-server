package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var Address string = "127.0.0.1"
var Port string = "8080"

func main() {

	config := GetConfig()

	// Configure output to display color
	if config.Other.PrettyOutput {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	// Validate config
	if ValidateConfig(config) {
		DisplayConfig(config)
		createHandlers(config)
	} else {
		fmt.Println("Aborting.")
	}
}

func createHandlers(config Config) {

	// Create Router
	r := mux.NewRouter()

	// Handle Routes
	r.HandleFunc("/", StatsHandler)

	// Listen for requests
	address := config.Server.Address + ":" + config.Server.Port
	log.Info().Msg("Listening for connections...")
	http.ListenAndServe(address, r)
}

func StatsHandler(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("Request received to /")
}
