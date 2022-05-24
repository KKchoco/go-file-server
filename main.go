package main

import (
	"errors"
	"fmt"
	"io/fs"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

var config Config

type Request struct {
	File     multipart.FileHeader `form:"file"`
	Password string               `form:"password"`
}

func main() {
	// Check if files folder exists, if not, create it
	if _, err := os.Stat("./files"); errors.Is(err, fs.ErrNotExist) {
		fmt.Println("./files folder does not exist, creating...")
		if err := os.Mkdir("./files", 0755); err != nil {
			fmt.Println("Error creating ./files folder")
			return
		}
	}

	c, err := GetConfig()
	if err != nil {
		fmt.Println(err)
		return
	}
	config = c

	r := gin.Default()
	r.SetTrustedProxies(nil)

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"github": "https://github.com/lorencerri/sharex-server-golang",
		})
	})
	r.POST("/upload", uploadHandler)
	r.GET("/:file", func(c *gin.Context) {
		file := c.Param("file")

		// Check if file exists
		if _, err := os.Stat("./files/" + file); errors.Is(err, fs.ErrNotExist) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "File not found",
			})
			return
		}

		c.File("./files/" + file)
	})

	addr := ""
	if config.Server.Address != "" {
		addr = config.Server.Address
	}
	if config.Server.Port != "" {
		addr += ":" + config.Server.Port
	}
	if addr == "" {
		addr = ":8080"
	}

	r.Run(addr)

}

func uploadHandler(c *gin.Context) {

	// Bind the request body to the struct
	var req Request
	if err := c.Bind(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Malformed request",
		})
		return
	}

	// Validate password if exists
	if len(config.Files.Password) > 0 && req.Password != config.Files.Password {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Incorrect password",
		})
		return
	}

	// Check if Content-Length exceeds max size
	if req.File.Size > config.Files.MaxUploadSize<<20 {
		c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, gin.H{
			"message": fmt.Sprintf("Request or file too large (%v > %v)", req.File.Size, config.Files.MaxUploadSize<<20),
		})
		return
	}

	// Check if valid extension
	ext := filepath.Ext(req.File.Filename)
	if !Contains(config.Files.AllowedFileTypes, ext) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("File type not allowed (%v)", ext),
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
			"message": "Error saving file",
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
		"url":  protocol + c.Request.Host + "/" + filename,
		"size": req.File.Size,
	})
}
