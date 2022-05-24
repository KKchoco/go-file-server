package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	bolt "go.etcd.io/bbolt"
)

const API_PATH = "/api"

type Request struct {
	File     multipart.FileHeader `form:"file"`
	Password string               `form:"password"`
}

type File struct {
	Name  string
	Views int
}

func CreateAPI(r *gin.Engine) {
	r.POST(API_PATH+"/upload", uploadHandler)
	r.GET(API_PATH+"/:file", fileHandler)
	r.GET(API_PATH+"/:file/stats", statsHandler)
}

func statsHandler(c *gin.Context) {
	file := c.Param("file")

	// Check if file exists
	if _, err := os.Stat("./files/" + file); errors.Is(err, fs.ErrNotExist) {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": "File not found",
		})
		return
	}

	database.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte("files"))
		if b == nil {
			return errors.New("bucket does not exist")
		}

		data := File{}

		// Umarshall the byte array encoded data into a struct
		if err := json.Unmarshal(b.Get([]byte(file)), &data); err != nil {
			fmt.Println(err)
			data = File{
				Name:  file,
				Views: 0,
			}
		}

		c.JSON(200, gin.H{
			"name":  data.Name,
			"views": data.Views,
		})

		return nil
	})
}

func fileHandler(c *gin.Context) {
	file := c.Param("file")

	// Check if file exists
	if _, err := os.Stat("./files/" + file); errors.Is(err, fs.ErrNotExist) {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": "File not found",
		})
		return
	}

	// Increment in database
	database.Update(func(tx *bolt.Tx) error {

		b, err := tx.CreateBucketIfNotExists([]byte("files"))
		if err != nil {
			fmt.Println(err)
			return errors.New("could not create bucket")
		}

		data := File{}

		// Umarshall the byte array encoded data into a struct
		if err := json.Unmarshal(b.Get([]byte(file)), &data); err != nil {
			fmt.Println(err)
			data = File{
				Name:  file,
				Views: 0,
			}
		}

		// Modify Data
		data.Views += 1

		// Remarshall
		encoded, err := json.Marshal(data)
		if err != nil {
			fmt.Println(err)
			return errors.New("could not encode file")
		}

		return b.Put([]byte(data.Name), encoded)
	})

	c.File("./files/" + file)
}

func uploadHandler(c *gin.Context) {

	// Bind the request body to the struct
	var req Request
	if err := c.Bind(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Malformed request",
		})
		return
	}

	// Check if file name is too long
	if len(req.File.Filename) > 255 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "File name too long",
		})
		return
	}

	// Validate password if exists
	if len(config.Files.Password) > 0 && req.Password != config.Files.Password {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Incorrect password",
		})
		return
	}

	// Check if Content-Length exceeds max size
	if req.File.Size > config.Files.MaxUploadSize<<20 {
		c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, gin.H{
			"error": fmt.Sprintf("Request or file too large (%v > %v)", req.File.Size, config.Files.MaxUploadSize<<20),
		})
		return
	}

	// Check if valid extension
	ext := filepath.Ext(req.File.Filename)
	if !Contains(config.Files.AllowedFileTypes, ext) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("File type not allowed (%v)", ext),
		})
		return
	}

	// Generate random file name
	filename := ""
	if config.Files.ObfuscateFileNames {
		filename = RandString(config.Files.KeyLength) + ext
	} else {
		filename = req.File.Filename
	}

	// Save file to location
	if err := c.SaveUploadedFile(&req.File, "./files/"+filename); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Could not save file",
		})
	}

	// Determine protocol based on config
	protocol := ""
	if config.Server.HTTPS {
		protocol = "https://"
	} else {
		protocol = "http://"
	}

	// Return success with information
	c.JSON(200, gin.H{
		"url":  protocol + c.Request.Host + API_PATH + "/" + filename,
		"size": req.File.Size,
	})
}
