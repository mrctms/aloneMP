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
var srv = flag.String("srv", "tcp", "aloneMPd server type")
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

		var client clients.Clienter
		var sender senders.Sender

		if *srv == "tcp" {
			sender, err = senders.NewTcpSender(*address)
			if err != nil {
				return
			}
		} else if *srv == "http" {
			sender, err = senders.NewHttpSender(*address)
			if err != nil {
				return
			}
		}

		if *tui {
			client = terminal.NewTerminalClient()
		}

		client.SetSender(sender)
		client.Run(*dir)
	}

}
