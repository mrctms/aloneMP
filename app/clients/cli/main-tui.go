package cli

import (
	"fmt"
	"util"

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
		AddItem(t.tracksList.TreeView, 3, 0, 10, 5, 0, 0, true).
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

func (t *Tui) Run() error {
	t.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter:
			trackToPlay := t.tracksList.GetSelectedTrackName()
			if trackToPlay != "" {
				t.TrackSelected <- trackToPlay
				return nil
			}
		case tcell.KeyCtrlP:
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
		return err
	}
	return nil
}

func (t *Tui) Stop() {
	t.app.Stop()
}

func (t *Tui) Draw() {
	t.app.Draw()
}

func (t *Tui) setCmdText() {
	t.cmd.SetText(`
		[yellow](↑)(↓) Browse Track
		[yellow] (←)(→) Volume
		[yellow](Enter) Play Selected Track 
		[yellow](Ctrl+P) Pause/Resume
		[yellow](Ctrl+Space) Mute[yellow] 
		[yellow](Ctrl+C) Quit[yellow]`)
}

func (t *Tui) SetTrackInfo(info util.TrackInfo) {
	var text string
	text = fmt.Sprintf("[blue]Title:[white] \n%s\n"+
		"[blue]Album:[white] \n%s\n"+
		"[blue]Artist:[white] \n%s\n"+
		"[blue]Genre:[white] \n%s\n"+
		"[blue]Year:[white] \n%d\n", info.Title, info.Album, info.Artist, info.Genre, info.Year)

	t.trackInfo.SetInfo(text)
}

func (t *Tui) PopolateTracksList(rootPath string) {
	t.tracksList.AddItems(rootPath)
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
	t.tracksList.NextTrack()
	return t.tracksList.GetSelectedTrackName()
}

func (t *Tui) PreviousTrack() string {
	t.tracksList.PreviousTrack()
	return t.tracksList.GetSelectedTrackName()
}

func (t *Tui) TrackList() []string {
	return t.tracksList.GetAllTracks()
}
func (t *Tui) CurrentTrack() string {
	return t.tracksList.GetSelectedTrackName()
}
