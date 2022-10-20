package main

import (
	"errors"
	"fmt"
	"io/fs"
	"io"
	"os"

	"github.com/gin-gonic/gin"
	"go.etcd.io/bbolt"
)

//const GIN_MODE = "release"
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

	//gin.SetMode(gin.ReleaseMode)
    gin.DisableConsoleColor()
    f, _ := os.Create("gin.log") // can be /var/log/gin.log for example
    gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	r := gin.Default()

	//r.SetTrustedProxies([]string{"127.0.0.1", "10.241.66.1"})
    r.SetTrustedProxies([]string{"0.0.0.0/0"}) // If you use nginx as reverse proxy then use address of your nginx, like example above
	// Create routes
	CreateAPI(r)

	r.Run(getAddr())

}

func preflight() error {
	// Check if files folder exists, if not, create it
	if _, err := os.Stat(config.Files.FilesPath); errors.Is(err, fs.ErrNotExist) {
		fmt.Println(config.Files.FilesPath + " folder does not exist, creating...")
		if err := os.Mkdir(config.Files.FilesPath, 0755); err != nil {
			return errors.New("error creating " + config.Files.FilesPath +  " folder")
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
