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

const API_PATH = "/api"

type Request struct {
	File     multipart.FileHeader `form:"file"`
	Password string               `form:"password"`
}

func CreateAPI(r *gin.Engine) {
	r.POST(API_PATH+"/upload", uploadHandler)
	r.GET(API_PATH+"/:file", fileHandler)
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
		"url":  protocol + c.Request.Host + "/" + filename,
		"size": req.File.Size,
	})
}
