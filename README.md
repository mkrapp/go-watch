# go-watch
A simple console-based time tracker with JSON support.

Requires `gocui`:

    go get github.com/jroimartin/gocui
    
    
Run with

    go run go-watch.go

or build and run with

    go build go-watch.go
    ./go-watch

A JSON file (`test.json`) is used to store the start of the project, the total elapsed time,
and the elasped time from the last session.

## Screenshot

![screenshot of go-watch](https://cloud.githubusercontent.com/assets/5938262/10237033/060ce9d4-68ad-11e5-89db-e4b29eaf9497.png "screenshot of go-watch")
