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
	"flag"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"aloneMP/media"
	"aloneMP/ui"

	"github.com/gizak/termui/v3"
)

var (
	dir           = flag.Bool("dir", false, "Directory with audio files")
	knowExtension = [4]string{".mp3", ".wav", ".flac", ".ogg"}
)

func getFiles(dir string) []string {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatalln(err)
		}
		if contains(knowExtension, filepath.Ext(path)) {
			file := strings.TrimLeft(strings.Replace(path, dir, "", -1), "/\\")
			files = append(files, file)
		}
		return nil
	})
	if err != nil {
		log.Fatalln(err)
	}
	return files
}

func contains(s [4]string, e string) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}

func Run(files []string) {

	if err := termui.Init(); err != nil {
		log.Fatalf("Failed to initialize TUI: %v\n", err)
	}
	defer termui.Close()

	tui := ui.NewTui()
	player := media.NewPlayer()

	for _, v := range files {
		tui.SongsList.Rows = append(tui.SongsList.Rows, v)
	}

	ticker := time.NewTicker(time.Second).C

	go player.StartPlayer()

	for {
		select {
		case e := <-termui.PollEvents():
			switch e.ID {
			case "q":
				return
			case "<Up>":
				tui.SongsList.ScrollUp()
			case "<Down>":
				tui.SongsList.ScrollDown()
			case "<Enter>":
				if player.IsPlaying {
					player.Close()
				}
				player.SongToPlay = tui.SongsList.Rows[tui.SongsList.SelectedRow]
				player.Play <- true
			case "<Space>":
				player.PaRes <- true
			case "m":
				player.Mute <- true
			case "<Left>":
				player.VolumeDown <- true
			case "<Right>":
				player.VolumeUp <- true
			case "<Resize>":
				payload := e.Payload.(termui.Resize)
				tui.SongsList.SetRect(0, 10, payload.Width/2, payload.Height-5)
				tui.SongInfo.SetRect(payload.Width, 10, payload.Width/2, payload.Height-5)
				tui.Commands.SetRect(payload.Width-30, 1, payload.Width, 9)
				tui.Banner.SetRect(-1, -1, payload.Width-35, 8)
				tui.SongProgress.SetRect(0, payload.Height-5, payload.Width, payload.Height-2)
				termui.Clear()
				tui.RedrawAll()
			}
		case <-ticker:
			tui.SetProgDur(player.Progress, player.Duration, player.SongLength)
			tui.RedrawAll()
			tui.UpdateInfo(player.SongInfo)
		case <-player.Finished:
			tui.SongsList.SelectedRow++
			if tui.SongsList.SelectedRow >= len(tui.SongsList.Rows) {
				tui.SongsList.SelectedRow = 0
			}
			player.SongToPlay = tui.SongsList.Rows[tui.SongsList.SelectedRow]
			player.Play <- true
		case err := <-player.PlayingError:
			termui.Close()
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
			files := getFiles(musicDir)
			if len(files) != 0 {
				os.Chdir(musicDir)
				Run(files)
			} else {
				log.Fatalln("No audio file found")
			}
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
			files := getFiles(userDir)
			if len(files) != 0 {
				os.Chdir(userDir)
				Run(files)
			} else {
				log.Fatalln("No audio file found")
			}
		} else {
			log.Fatalf("%s is not a directory\n", userDir)
		}

	}
}
