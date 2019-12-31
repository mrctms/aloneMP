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

package app

import (
	"log"
	"time"

	"github.com/gizak/termui/v3"
	"github.com/marcktomack/aloneMP/player"
	"github.com/marcktomack/aloneMP/ui"
)

func Run(files []string) {

	if err := termui.Init(); err != nil {
		log.Fatalf("Failed to initialize TUI: %v\n", err)
	}
	defer termui.Close()

	tui := ui.NewUi()
	pr := player.NewPlayer()

	for _, v := range files {
		tui.SongsList.Rows = append(tui.SongsList.Rows, v)
	}

	ticker := time.NewTicker(time.Second).C

	go pr.StartPlayer()

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
				if pr.IsPlaying {
					pr.Close()
				}
				pr.SongToPlay = tui.SongsList.Rows[tui.SongsList.SelectedRow]
				pr.Play <- true
			case "<Space>":
				pr.PaRes <- true
			case "m":
				pr.Mute <- true
			case "<Left>":
				pr.VolumeDown <- true
			case "<Right>":
				pr.VolumeUp <- true
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
			tui.SetProgDur(pr.Progress, pr.Duration, pr.SongLength)
			tui.RedrawAll()
			tui.UpdateInfo(pr.SongInfo)
		case <-pr.Finished:
			tui.SongsList.SelectedRow++
			if tui.SongsList.SelectedRow >= len(tui.SongsList.Rows) {
				tui.SongsList.SelectedRow = 0
			}
			pr.SongToPlay = tui.SongsList.Rows[tui.SongsList.SelectedRow]
			pr.Play <- true
		case <-pr.PlayingError:
			termui.Close()
			log.Fatalln(pr.ErrorMsg)
		}

	}
}
