package main

import (
	"aloneMPd/media"
	"aloneMPd/server"
	"flag"
	"fmt"
	"log"
	"os"
	"util"

	"github.com/marcsauter/single"
)

var address *string
var version string

var ver = flag.Bool("version", false, "show version")

func main() {
	logger, err := util.NewLogger("daemon-main")
	if err != nil {
		log.Fatalln(err)
	}
	hostname, err := os.Hostname()
	if err != nil {
		logger.Write(fmt.Sprintf("failed to get hostname %v", err))
		log.Fatalln(err)
	}
	address = flag.String("addr", fmt.Sprintf("%s:3777", hostname), "address")

	flag.Parse()
	if *ver {
		fmt.Printf("\naloneMPd version: %s\n\n", version)
	} else {

		s := single.New("aloneMPd")

		if err := s.CheckLock(); err != nil && err == single.ErrAlreadyRunning {
			msg := "another instance is running"
			logger.Write(msg)
			log.Fatalln(msg)
		} else {
			defer s.TryUnlock()
		}
		player := media.NewMusicPlayer()
		listener := server.NewTcpServer()

		defer func() {
			if r := recover(); r != nil {
				msg := fmt.Sprintf("Something went wrong: %v", r)
				logger.Write(msg)
				fmt.Println(msg)
			} else {
				if player != nil {
					player.Close()
				}
				if listener != nil {
					listener.Close()
				}
			}
		}()

		go listener.Listen(*address)

		for {
			select {
			case pArgs := <-listener.Initialize():
				err := player.Initialize(pArgs)
				if err != nil {
					panic(err)
				}
				go player.Start()
				listener.SetPlayerInfo(player.PlayerInfo())
			case track := <-listener.Play():
				player.Play(track)
			case <-listener.Mute():
				player.Mute()
			case <-listener.Pause():
				player.Pause()
			case <-listener.VolumeUp():
				player.VolumeUp()
			case <-listener.VolumeDown():
				player.VolumeDown()
			case <-listener.ShutDown():
				player.Close()
			case err := <-player.FatalError():
				panic(err)
			case err := <-listener.FatalError():
				panic(err)
			}
		}
	}
}
