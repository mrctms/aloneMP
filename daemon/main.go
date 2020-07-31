package main

import (
	"aloneMPd/media"
	"aloneMPd/server"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/marcsauter/single"
)

var address *string
var version string

var ver = flag.Bool("version", false, "show version")

func main() {
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalln(err)
	}
	address = flag.String("addr", fmt.Sprintf("%s:3777", hostname), "address")

	flag.Parse()
	if *ver {
		fmt.Printf("\naloneMPd version: %s\n\n", version)
	} else {
		s := single.New("aloneMPd")

		if err := s.CheckLock(); err != nil && err == single.ErrAlreadyRunning {
			log.Fatalln("another instance is running")
		} else {
			defer s.TryUnlock()
		}
		var player *media.FilePlayer

		defer func() {
			if player != nil {
				if player.Info().IsPlaying() {
					player.Close()
				}
			}

		}()
		listener := server.NewTcpServer()

		go listener.Listen(*address)

		for {
			select {
			case dir := <-listener.Source():
				var err error
				player, err = media.NewFilePlayer()
				if err != nil {
					return
				}
				listener.SetPlayerInfo(player.Info())
				go player.Start(dir)
			case track := <-listener.SelectedTrack():
				if player != nil {
					if player.Info().IsPlaying() {
						player.Clear()
					}
					player.SetTrackToPlay(track)
					player.Play()
				}
			case <-listener.NextTrack():
				if player != nil {
					if player.Info().IsPlaying() {
						player.Clear()
					}
					nextTrack := player.Info().NextTrack()
					if nextTrack != "" {
						player.SetTrackToPlay(nextTrack)
						player.Play()
					}
				}
			case <-listener.PreviousTrack():
				if player != nil {
					if player.Info().IsPlaying() {
						player.Clear()
					}
					previousTrack := player.Info().PreviousTrack()
					if previousTrack != "" {
						player.SetTrackToPlay(previousTrack)
						player.Play()
					}
				}
			case <-listener.MuteTrack():
				if player != nil {
					player.Mute()
				}
			case <-listener.PauseTrack():
				if player != nil {
					player.Pause()
				}
			case <-listener.VolumeUp():
				if player != nil {
					player.VolumeUp()
				}
			case <-listener.VolumeDown():
				if player != nil {
					player.VolumeDown()
				}
			case <-listener.ShutDown():
				if player != nil {
					if player.Info().IsPlaying() {
						player.Close()
					}
				}
			}
		}
	}
}
