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

package ui

import (
	"fmt"
	"log"

	"github.com/dhowden/tag"
	"github.com/gdamore/tcell"
	"gitlab.com/tslocum/cview"
)

const aloneMP = `
    _    _                  __  __ ____  
   / \  | | ___  _ __   ___|  \/  |  _ \ 
  / _ \ | |/ _ \| '_ \ / _ \ |\/| | |_) |
 / ___ \| | (_) | | | |  __/ |  | |  __/ 
/_/   \_\_|\___/|_| |_|\___|_|  |_|_| 
`

type Tui struct {
	app           *cview.Application
	grid          *cview.Grid
	tracksList    *TracksList
	cmd           *CmdView
	trackInfo     *TrackInfo
	trackProgress *ProgressView
	TrackSelected chan string
	Paused        chan bool
	Mute          chan bool
	VolumeUp      chan bool
	VolumeDown    chan bool
	Quit          chan bool
}

func NewTui() *Tui {
	banner := cview.NewTextView().SetTextAlign(cview.AlignLeft).SetText(aloneMP)
	banner.SetBorder(false)
	banner.SetTextColor(tcell.ColorDarkBlue)
	t := new(Tui)
	t.tracksList = NewTracksList()

	t.cmd = NewCmdView()
	t.trackInfo = NewTrackInfo()
	t.trackProgress = NewProgressView()
	t.setCmdText()
	t.grid = cview.NewGrid().
		SetRows(4, 4).
		SetColumns(40).
		AddItem(banner, 0, 0, 12, 3, 3, 3, false).
		AddItem(t.tracksList.List, 3, 0, 10, 5, 0, 0, true).
		AddItem(t.trackInfo.TextView, 3, 5, 10, 6, 0, 0, false).
		AddItem(t.cmd.TextView, 0, 4, 3, 7, 0, 0, false).
		AddItem(t.trackProgress, 13, 0, 1, 11, 0, 0, false)

	t.app = cview.NewApplication()
	t.app.SetRoot(t.grid, true)
	t.TrackSelected = make(chan string)
	t.Paused = make(chan bool)
	t.Mute = make(chan bool)
	t.VolumeDown = make(chan bool)
	t.VolumeUp = make(chan bool)
	t.Quit = make(chan bool)
	return t
}

func (t *Tui) Run() {
	t.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter:
			t.TrackSelected <- t.tracksList.GetSelectedItemText()
		case tcell.KeyRune:
			t.Paused <- true
			return nil
		case tcell.KeyCtrlSpace:
			t.Mute <- true
			return nil
		case tcell.KeyLeft:
			t.VolumeDown <- true
			return nil
		case tcell.KeyRight:
			t.VolumeUp <- true
			return nil
		case tcell.KeyCtrlC:
			t.Quit <- true
			return nil
		}
		return event
	})
	err := t.app.Run()
	if err != nil {
		log.Println(err)
	}
}

func (t *Tui) Stop() {
	t.app.Stop()
}

func (t *Tui) Draw() {
	t.app.Draw()
}

func (t *Tui) setCmdText() {
	t.cmd.SetText(`
		[yellow](↑)(↓) Browse Songs
		[yellow] (←)(→) Volume
		[yellow](Enter) Play Selected Song 
		[yellow](Space) Pause/Resume
		[yellow](Ctrl+Space) Mute[yellow] 
		[yellow](Ctrl+C) Quit[yellow]`)
}

func (t *Tui) SetTrackInfo(info tag.Metadata) {
	var text string
	if info == nil {
		text = "[blue]File Type:\n" +
			"[blue]Title:\n" +
			"[blue]Album:\n" +
			"[blue]Artist:\n" +
			"[blue]AlbumArtist:\n" +
			"[blue]Composer:\n" +
			"[blue]Genre:\n" +
			"[blue]Year:\n"
	} else {
		text = fmt.Sprintf("[blue]File Type:[white] \n%s\n"+
			"[blue]Title: \n%s\n"+
			"[blue]Album: \n%s\n"+
			"[blue]Artist: \n%s\n"+
			"[blue]AlbumArtist: \n%s\n"+
			"[blue]Composer: \n%s\n"+
			"[blue]Genre: \n%s\n"+
			"[blue]Year: \n%d\n", info.FileType(), info.Title(), info.Album(), info.Artist(), info.AlbumArtist(), info.Composer(), info.Genre(), info.Year())
	}
	t.trackInfo.SetInfo(text)
}

func (t *Tui) PopolateTracksList(items []string) {
	t.tracksList.AddItems(items)
}

func (t *Tui) SetProgDur(prog string, dur string, percentage int) {
	if prog == "" && dur == "" {
		t.trackProgress.SetProgressTitle("00:00:00/00:00:00")
	} else {
		t.trackProgress.SetProgressTitle(fmt.Sprintf("%s/%s", prog, dur))
	}
	t.trackProgress.UpdateProgress(percentage)
}

func (t *Tui) NextTrack() string {
	return t.tracksList.NextItem()
}
