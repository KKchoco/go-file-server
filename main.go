package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const filesPath = "/files"

var config Config

type ResponseObject struct {
	Name string
	Url  string
	Size uint
}

func main() {

	config = GetConfig()

	// Configure console to display color
	if config.Other.PrettyOutput {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	// Validate config
	if isValidConfig(config) {
		DisplayConfig(config)
		createHandlers(config)
	} else {
		log.Fatal().Msg("Aborting")
	}
}

func createHandlers(config Config) {

	// Create Router
	router := mux.NewRouter()

	// Handle Routes
	router.HandleFunc("/", StatsHandler)
	router.HandleFunc("/upload", UploadHandler)
	router.Handle("/{file}", http.FileServer(http.Dir("./files")))

	// Listen for requests
	address := config.Server.Address + ":" + config.Server.Port
	log.Info().Msg("Listening for connections...")
	http.ListenAndServe(address, router)
}

func StatsHandler(response http.ResponseWriter, request *http.Request) {
	log.Info().Msg("Request received to stats handler")
}

func UploadHandler(response http.ResponseWriter, request *http.Request) {
	log.Info().Msg("Beginning the upload process...")
	start := time.Now()

	// Validate the request's password
	password := request.FormValue("password")
	if len(config.Files.Password) > 0 && password != config.Files.Password {
		returnError(response, http.StatusUnauthorized, "Invalid password provided in 'password' field of Form Data")
		return
	}

	// Calculate the file size in bytes
	maxBytes := config.Files.MaxUploadSize << 20 // Megabytes -> Bytes Conversion

	// Limit the request's body size
	request.Body = http.MaxBytesReader(response, request.Body, maxBytes)

	// Get the first file from the request
	if file, fileHeader, err := request.FormFile("file"); err == nil {
		defer file.Close() // When the enclosing function ends, close the file
		log.Info().Msgf("File Size: %v", fileHeader.Size)

		// Read the bytes and store it in a variable
		if bytes, err := ioutil.ReadAll(file); err == nil {

			// Get the file type based on the file's bytes
			// This is better than getting the Content-Type from
			// the fileHeader map as it can be changed by the user
			contentType := http.DetectContentType(bytes)
			log.Info().Msgf("File Type: %v", contentType)

			// Get the string file extension based on the content type
			extensions, err := mime.ExtensionsByType(contentType)
			if err == nil {
				ext := extensions[0] // Get the first extension
				log.Info().Msgf("File Extension: %v", ext)

				// Check if the file type is included in the allowed file types
				if contains(config.Files.AllowedFileTypes, ext) {

					// Generate a random file name + information
					id, _ := gonanoid.Generate("ABCDEFGHIJKLMNPQRSTUVWXYZabcdefghjkmnpqrstuvwxyz", config.Files.KeyLength)
					name := id + ext
					path := filepath.Join(".", config.Files.FilesPath, name)

					// Create the file
					if createdFile, err := os.Create(path); err == nil {
						defer createdFile.Close() // When the enclosing function ends, close the file

						// Write the bytes to the file
						if _, err := createdFile.Write(bytes); err == nil {

							// Determine protocol based on config
							protocol := ""
							if config.Server.HTTPS {
								protocol = "https://"
							} else {
								protocol = "http://"
							}

							// Create response object to be passed
							jsonObj := ResponseObject{
								Name: name,
								Url:  protocol + request.Host + "/" + name,
								Size: uint(fileHeader.Size),
							}

							if jsonResp, err := json.Marshal(jsonObj); err == nil {
								log.Info().Msgf("Created file with url: %v", jsonObj.Url)

								// Return JSON object of information
								response.Header().Set("Content-Type", "application/json")
								response.WriteHeader(http.StatusOK)
								response.Write(jsonResp)

								duration := time.Since(start)
								log.Info().Msgf("Captured and uploaded image in %v", duration)

							} else {
								returnError(response, http.StatusInternalServerError, "Unexpected error, unable to create JSON response object")
							}
						} else {
							returnError(response, http.StatusInternalServerError, "Unexpected error, unable to write to the created file")
						}
					} else {
						returnError(response, http.StatusInternalServerError, "Unexpected error, unable to create a new file")
					}
				} else {
					returnError(response, http.StatusInternalServerError, fmt.Sprintf("The extension %v is not in the predefined list of allowed extensions", ext))
				}
			} else {
				returnError(response, http.StatusInternalServerError, fmt.Sprintf("Unable to get file extension of Content-Type %v", contentType))
			}
		} else {
			returnError(response, http.StatusInternalServerError, "Unexpected error, unable to read the file")
		}
	} else {
		returnError(response, http.StatusBadRequest, fmt.Sprintf("The file size is limited to %v Bytes, although the approx. size of the request was %v Bytes", maxBytes, request.ContentLength))
	}
}

func returnError(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	w.Write([]byte(message))
	log.Error().Msgf("HTTP Error %v: %v", code, message)
}

// https://freshman.tech/snippets/go/check-if-slice-contains-element/
func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
