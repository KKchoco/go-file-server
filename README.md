# go-file-server

A simple file server written in Go, designed for ShareX

## Features

-   Quick & easy setup
-   Easy to use dashboard
-   Download ShareX config button
-   Password protected uploading
-   Optional filetype restrictions

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
|   POST | /api/upload                      |
|    GET | /api/{fileName}                  |
|    GET | /api/{fileName}/stats            |
|    GET | /api/{fileName}/delete/{editKey} |
|    GET | /api/files/{adminPassword}       |
