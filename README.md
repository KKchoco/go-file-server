# go-file-server

A simple file server written in Go, designed for ShareX

## Features

-   Quick & easy setup
-   Easy to use dashboard
-   Download ShareX config button
-   Password protected uploading
-   Optional filetype restrictions

## Main differences from original branch

-   "Superkey" cookie that saves on valid admin password and removes requierment of both passwords 
-   API_PATH set to "/", main page on "/public/" so files hosted kinda like in [chibisafe](https://github.com/chibisafe/chibisafe)
-   Upload date saves to db and shows on dashboard 
-   Logging to file
-   Working "filesPath" setting in config
-   Better filename randomization
-   "Download ShareX Config" button hidden by default, use DevTools to unhide the button
-   Huge frontend optimizations, initial page size is only 6.8kb (4.46kb gzipped), no external dependences

## Preview

[<img alt="demo_gif" src="https://fs.plexidev.org/api/pICAQZm.gif" />](https://fs.plexidev.org/api/pICAQZm.gif)
[<img alt="demo_gif" src="https://fs.plexidev.org/api/ahYHMSG.gif" />](https://fs.plexidev.org/api/ahYHMSG.gif)

## Usage

1. Install [Go](https://go.dev) ([Ubuntu](https://github.com/golang/go/wiki/Ubuntu))
2. Clone repo `git clone https://github.com/lorencerri/sharex-server-golang`
3. Install dependencies `go get`
4. Copy & modify `example.config.yml` -> `config.yml`
5. Run program `go run .`

## API

| Method | Endpoint                         |
| -----: | :------------------------------- |
|   POST | /upload                          |
|    GET | /{fileName}                      |
|    GET | /{fileName}/stats                |
|    GET | /{fileName}/delete/{editKey}     |
|    GET | /files/{adminPassword}           |
