package main

import (
	"aloneMP/app"
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
)

// TODO: use cobra

var version string

var dir *string
var server = flag.Bool("s", false, "Run the http server")
var address = flag.String("addr", "127.0.0.1:3777", "http server address")
var ver = flag.Bool("version", false, "version")

func main() {
	mainApp := app.NewApp()
	u, err := user.Current()
	if err != nil {
		log.Fatalln(err)
	}
	musicDir := filepath.Join(u.HomeDir, "Music")
	dir = flag.String("dir", musicDir, "Directory with audio files")
	flag.Parse()

	if *ver {
		fmt.Printf("\naloneMP version: %s\n\n", version)
	} else {
		f, err := os.Stat(*dir)
		if err != nil {
			log.Fatalln(err)
		}
		if f.IsDir() {
			if *server {
				mainApp.RunHttpServer(*address)
			}
			mainApp.Run(*dir)
		} else {
			log.Fatalf("%s is not a directory\n", *dir)
		}
	}

}
