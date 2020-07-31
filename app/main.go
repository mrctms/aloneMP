package main

import (
	"aloneMP/clients"
	"aloneMP/clients/terminal"
	"aloneMP/senders"
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"

	"github.com/marcsauter/single"
)

var dir *string
var address *string

var tui = flag.Bool("tui", true, "run tui client")
var ver = flag.Bool("version", false, "show version")

var version string

func main() {
	u, err := user.Current()
	if err != nil {
		log.Fatalln(err)
	}
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalln(err)
	}
	musicDir := filepath.Join(u.HomeDir, "Music")
	dir = flag.String("dir", musicDir, "Directory with audio files")
	address = flag.String("addr", fmt.Sprintf("%s:3777", hostname), "aloneMP daemon address")
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

		var client clients.Clienter
		sender, err := senders.NewTcpSender(*address)
		if err != nil {
			log.Fatalln(err)
		}

		if *tui {
			client = terminal.NewTerminalClient()
		}

		client.SetSender(sender)
		client.Run(*dir)
	}

}
