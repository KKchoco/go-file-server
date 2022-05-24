package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"

	"github.com/gin-gonic/gin"
	"go.etcd.io/bbolt"
)

var config Config
var database *bbolt.DB

func main() {

	// Run preflight checks
	if err := preflight(); err != nil {
		fmt.Println(err)
		return
	}

	// Handle config
	c, err := GetConfig()
	if err != nil {
		fmt.Println(err)
		return
	}
	config = c

	// Handle database
	db, err := bbolt.Open("files.db", 0600, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()
	database = db

	// Create router
	r := gin.Default()
	r.SetTrustedProxies(nil)

	// Create routes
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"github": "https://github.com/lorencerri/sharex-server-golang",
		})
	})
	CreateAPI(r)

	r.Run(getAddr())

}

func preflight() error {
	// Check if files folder exists, if not, create it
	if _, err := os.Stat("./files"); errors.Is(err, fs.ErrNotExist) {
		fmt.Println("./files folder does not exist, creating...")
		if err := os.Mkdir("./files", 0755); err != nil {
			return errors.New("error creating ./files folder")
		}
	}
	return nil
}

func getAddr() string {
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
	return addr
}
