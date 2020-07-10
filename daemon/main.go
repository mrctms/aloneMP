package main

import (
	"aloneMPd/media"
	"aloneMPd/server"
	"flag"
	"fmt"
	"log"

	"github.com/marcsauter/single"
)

var address = flag.String("addr", "127.0.0.1:3777", "address")
var version string
var ver = flag.Bool("version", false, "show version")

func main() {
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
		httpServer := server.NewHttpServer()

		go httpServer.ListenAndServe(*address)

		for {
			select {
			case dir := <-httpServer.Source:
				var err error
				player, err = media.NewFilePlayer()
				if err != nil {
					log.Fatalln(err)
				}
				httpServer.PlayerInfo = player.Info()

				defer func() {
					if player.Info().IsPlaying() {
						player.Close()
					}
				}()
				go player.Start(dir)
			case track := <-httpServer.SelectedTrack:
				if player != nil {
					if player.Info().IsPlaying() {
						player.Clear()
					}
					player.SetTrackToPlay(track)
					player.Play()
				}
			case <-httpServer.NextTrack:
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
			case <-httpServer.PreviousTrack:
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
			case <-httpServer.MuteTrack:
				if player != nil {
					player.Mute()
				}
			case <-httpServer.PauseTrack:
				if player != nil {
					player.Pause()
				}
			case <-httpServer.VolumeUp:
				if player != nil {
					player.VolumeUp()
				}
			case <-httpServer.VolumeDown:
				if player != nil {
					player.VolumeDown()
				}
			case <-httpServer.ShutDown:
				if player != nil {
					if player.Info().IsPlaying() {
						player.Close()
					}
				}
			}
		}
	}
}
