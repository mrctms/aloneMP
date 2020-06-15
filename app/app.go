package app

import (
	"aloneMP/media"
	"aloneMP/server"
	"aloneMP/ui/tui"
	"fmt"
	"time"

	"github.com/dhowden/tag"
)

type App struct {
	player     media.Player
	tui        *tui.MainTui
	httpServer *server.HttpServer
}

func NewApp() *App {
	app := new(App)
	app.player = media.NewFilePlayer()
	app.tui = tui.NewMainTui()
	app.httpServer = server.NewHttpServer()
	app.httpServer.PlayerInfo = app.player.Info()
	app.httpServer.InterfaceInfo = app.tui
	return app
}

func (a *App) Run(rootPath string) {
	ticker := time.NewTicker(time.Second).C

	a.tui.PopolateTracksList(rootPath)

	go a.player.Start()
	go a.tui.Run()

	defer func() {
		if a.player.Info().IsPlaying() {
			a.player.Close()
		}
		a.tui.Stop()
	}()

	for {
		select {
		case track := <-a.tui.TrackSelected:
			if a.player.Info().IsPlaying() {
				a.player.Close()
			}
			a.player.SetTrackToPlay(track)
			a.player.Play()
		case <-a.tui.Paused:
			a.player.Pause()
		case <-a.tui.Mute:
			a.player.Mute()
		case <-a.tui.VolumeDown:
			a.player.VolumeDown()
		case <-a.tui.VolumeUp:
			a.player.VolumeUp()
		case <-ticker:
			info, ok := a.player.Info().TrackInfo().(tag.Metadata)
			if ok {
				a.tui.SetTrackInfo(info)
			}
			a.tui.SetProgDur(a.player.Info().Progress(), a.player.Info().Duration(), a.player.Info().TrackLength())
			a.tui.Draw()
		case <-a.player.Finished():
			nextTrack := a.tui.NextTrack()
			if nextTrack != "" {
				a.player.SetTrackToPlay(nextTrack)
				a.player.Play()
			}
		case <-a.tui.Quit:
			return
		case <-a.httpServer.NextTrack:
			if a.player.Info().IsPlaying() {
				a.player.Close()
			}
			nextTrack := a.tui.NextTrack()
			if nextTrack != "" {
				a.player.SetTrackToPlay(nextTrack)
				a.player.Play()
			}
		case <-a.httpServer.PreviousTrack:
			if a.player.Info().IsPlaying() {
				a.player.Close()
			}
			previousTrack := a.tui.PreviousTrack()
			if previousTrack != "" {
				a.player.SetTrackToPlay(previousTrack)
				a.player.Play()
			}
		case <-a.httpServer.MuteTrack:
			a.player.Mute()
		case <-a.httpServer.PauseTrack:
			a.player.Pause()
		case <-a.httpServer.VolumeUp:
			a.player.VolumeUp()
		case <-a.httpServer.VolumeDown:
			a.player.VolumeDown()
		case err := <-a.player.PlayingError():
			fmt.Println(err)
			return
		}
	}
}

func (a *App) RunHttpServer(address string) {
	go a.httpServer.ListenAndServe(address)
}
