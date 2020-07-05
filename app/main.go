package main

import (
	"aloneMP/clients"
	"aloneMP/clients/terminal"
	"aloneMP/senders"
	"flag"
	"fmt"
	"log"
	"os/user"
	"path/filepath"

	"github.com/marcsauter/single"
)

var dir *string
var address = flag.String("addr", "127.0.0.1:3777", "aloneMP daemon address")
var tui = flag.Bool("tui", true, "run tui client")
var ver = flag.Bool("version", false, "show version")

var version string

func main() {
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
		s := single.New("aloneMP")

		if err := s.CheckLock(); err != nil && err == single.ErrAlreadyRunning {
			log.Fatalln("another instance is running")
		} else {
			defer s.TryUnlock()
		}

		httpSender, err := senders.NewHttpSender(*address)
		if err != nil {
			log.Fatalln(err)
		}
		var client clients.Clienter
		if *tui {
			client = terminal.NewTerminalClient()
		}

		client.SetSender(httpSender)
		client.Run(*dir)
	}

}
