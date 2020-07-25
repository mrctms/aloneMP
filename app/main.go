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
var srv = flag.String("srv", "tcp", "aloneMP daemon server type")
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
