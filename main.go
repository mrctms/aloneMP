/*
Copyright (C) MarckTomack <marcktomack@tutanota.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"aloneMP/media"
	"aloneMP/ui"
	"flag"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"time"
)

var dir = flag.Bool("dir", false, "Directory with audio files")

func Run(rootPath string) {

	ticker := time.NewTicker(time.Second).C

	player := media.NewPlayer()

	tui := ui.NewTui()

	tui.PopolateTracksList(rootPath)

	go player.StartPlayer()
	go tui.Run()
	defer func() {
		if player.IsPlaying {
			player.Close()
		}
		tui.Stop()
	}()

	for {
		select {
		case track := <-tui.TrackSelected:
			if player.IsPlaying {
				player.Close()
			}
			player.TrackToPlay = track
			player.Play <- true
		case paused := <-tui.Paused:
			player.PaRes <- paused
		case mute := <-tui.Mute:
			player.Mute <- mute
		case volumeDown := <-tui.VolumeDown:
			player.VolumeDown <- volumeDown
		case volumeUp := <-tui.VolumeUp:
			player.VolumeUp <- volumeUp
		case <-ticker:
			tui.SetTrackInfo(player.TrackInfo)
			tui.SetProgDur(player.Progress, player.Duration, player.TrackLength)
			tui.Draw()
		case <-player.Finished:
			nextTrack := tui.NextTrack()
			if nextTrack != "" {
				player.TrackToPlay = nextTrack
				player.Play <- true
			}

		case <-tui.Quit:
			return
		case err := <-player.PlayingError:
			log.Fatalln(err)
		}
	}
}

func main() {
	flag.Parse()
	if len(os.Args) == 1 {
		u, _ := user.Current()
		musicDir := filepath.Join(u.HomeDir, "Music")
		f, err := os.Stat(musicDir)
		if err != nil {
			log.Fatalln(err)
		}
		if f.IsDir() {
			Run(musicDir)
		} else {
			log.Fatalf("%s is not a directory\n", musicDir)
		}
	}
	if *dir {
		userDir := os.Args[2]
		f, err := os.Stat(userDir)
		if err != nil {
			log.Fatalln(err)
		}
		if f.IsDir() {
			Run(userDir)
		} else {
			log.Fatalf("%s is not a directory\n", userDir)
		}

	}
}
