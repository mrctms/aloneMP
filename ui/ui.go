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

	"github.com/dhowden/tag"
	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

const aloneMP = `
    _    _                  __  __ ____  
   / \  | | ___  _ __   ___|  \/  |  _ \ 
  / _ \ | |/ _ \| '_ \ / _ \ |\/| | |_) |
 / ___ \| | (_) | | | |  __/ |  | |  __/ 
/_/   \_\_|\___/|_| |_|\___|_|  |_|_| 
`

type Ui struct {
	SongsList    *widgets.List
	SongProgress *widgets.Gauge
	SongInfo     *widgets.Paragraph
	Commands     *widgets.Paragraph
	banner       *widgets.Paragraph
}

func NewUi() *Ui {
	ui := new(Ui)
	w, h := termui.TerminalDimensions()
	ui.renderSongsList(w, h)
	ui.renderSongInfo(w, h)
	ui.renderBanner(w, h)
	ui.renderCommands(w, h)
	ui.renderProgress(w, h)
	return ui
}

func (u *Ui) SetProgDur(prog string, dur string, percentage int) {
	if prog == "" && dur == "" {
		u.SongProgress.Label = "00:00:00 / 00:00:00"
	} else {
		u.SongProgress.Label = fmt.Sprintf("%s / %s", prog, dur)
	}
	u.SongProgress.Percent = percentage
}

func (u *Ui) renderSongsList(w, h int) {
	l := widgets.NewList()
	l.Title = "Songs List"
	l.SetRect(0, 10, w/2, h-5)
	l.SelectedRowStyle = termui.NewStyle(termui.ColorWhite, termui.ColorBlue)
	u.SongsList = l
	termui.Render(u.SongsList)
}

func (u *Ui) renderSongInfo(w, h int) {
	i := widgets.NewParagraph()
	i.Title = "Info"
	i.SetRect(w, 10, w/2, h-5)
	u.SongInfo = i
	termui.Render(u.SongInfo)
}

func (u *Ui) UpdateInfo(info tag.Metadata) {
	var text string
	if info == nil {
		text = "[File Type:](fg:blue,mod:bold)\n" +
			"[Title:](fg:blue,mod:bold)\n" +
			"[Album:](fg:blue,mod:bold)\n" +
			"[Artist:](fg:blue,mod:bold)\n" +
			"[AlbumArtist:](fg:blue,mod:bold)\n" +
			"[Composer:](fg:blue,mod:bold)\n" +
			"[Genre:](fg:blue,mod:bold)\n" +
			"[Year:](fg:blue,mod:bold)\n"
	} else {
		text = fmt.Sprintf("[File Type:](fg:blue,mod:bold) %s\n"+
			"[Title:](fg:blue,mod:bold) %s\n"+
			"[Album:](fg:blue,mod:bold) %s\n"+
			"[Artist:](fg:blue,mod:bold) %s\n"+
			"[AlbumArtist:](fg:blue,mod:bold) %s\n"+
			"[Composer:](fg:blue,mod:bold) %s\n"+
			"[Genre:](fg:blue,mod:bold) %s\n"+
			"[Year:](fg:blue,mod:bold) %d\n", info.FileType(), info.Title(), info.Album(), info.Artist(), info.AlbumArtist(), info.Composer(), info.Genre(), info.Year())
	}
	u.SongInfo.Text = text
}

func (u *Ui) renderProgress(w, h int) {
	d := widgets.NewGauge()
	d.LabelStyle = termui.NewStyle(termui.ColorWhite, termui.ColorClear)
	d.Label = "00:00:00 / 00:00:00"
	d.SetRect(0, h-5, w, h-2)
	d.BarColor = termui.ColorGreen
	u.SongProgress = d
	termui.Render(u.SongProgress)

}

func (u *Ui) renderBanner(w, h int) {
	b := widgets.NewParagraph()
	b.Border = false
	b.Text = aloneMP
	b.TextStyle = termui.NewStyle(termui.ColorBlue, termui.ColorClear, termui.ModifierBold)
	b.SetRect(-1, -1, w-35, 8)
	u.banner = b
	termui.Render(u.banner)
}

func (u *Ui) renderCommands(w, h int) {
	c := widgets.NewParagraph()
	c.Border = false
	c.Text = "[[↑]](fg:yellow,mod:bold)[[↓]](fg:yellow,mod:bold) Browse Songs\n" +
		"[[Enter]](fg:yellow,mod:bold) Play Selected Song\n" +
		"[[Space]](fg:yellow,mod:bold) Pause/Resume\n" +
		"[[m]](fg:yellow,mod:bold) Mute\n" +
		"[[q]](fg:yellow,mod:bold) Quit"
	c.SetRect(w-30, 1, w, 9)
	u.Commands = c
	termui.Render(u.Commands)

}

func (u *Ui) RedrawAll() {
	termui.Render(u.SongProgress)
	termui.Render(u.SongsList)
	termui.Render(u.SongInfo)
	termui.Render(u.Commands)
}
