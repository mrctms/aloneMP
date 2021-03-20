package cli

import (
	"aloneMP/senders"
	"fmt"
	"time"
	"util"
)

type TerminalClient struct {
	tui    *Tui
	sender senders.Sender
}

func NewTerminalClient() *TerminalClient {
	terminalClient := new(TerminalClient)
	tui := NewTui()
	terminalClient.tui = tui
	return terminalClient
}

func (t *TerminalClient) SetSender(sender senders.Sender) {
	t.sender = sender
}

func (t *TerminalClient) Run(rootDir string) {

	ticker := time.NewTicker(time.Second).C
	aliveTicker := time.NewTicker(time.Second * 10).C
	defer func() {
		t.tui.Stop()
		t.sender.ShutDown()
	}()
	go t.tui.Run()

	t.sender.Initialize(rootDir)

	var tl *util.TrackListMessage
	c := 0
	for tl == nil {
		if c == 5 {
			t.tui.Stop()
			fmt.Println("timeout reached")
			return
		}
		tl = t.sender.TrackList()
		time.Sleep(1 * time.Second)
		c++
	}

	t.tui.PopolateTracksList(*tl)
	for {
		select {
		case track := <-t.tui.TrackSelected:
			t.sender.Play(track)
		case <-t.tui.Mute:
			t.sender.Mute()
		case <-t.tui.Paused:
			t.sender.Pause()
		case <-t.tui.VolumeUp:
			t.sender.VolumeUp()
		case <-t.tui.VolumeDown:
			t.sender.VolumeDown()
		case <-aliveTicker:
			if !t.sender.IsAlive() {
				t.tui.Stop()
				fmt.Println("aloneMPd is not alive")
				return
			}
		case <-ticker:
			info := t.sender.TrackInfo()
			if info != nil {
				if info.InError {
					t.tui.NextTrack()
					t.sender.Play(t.tui.CurrentTrack())
				} else if !info.IsPlaying {
					if (info.TrackLength == info.TrackProgress) && (info.Percentage != 0) {
						t.tui.NextTrack()
						t.sender.Play(t.tui.CurrentTrack())
					}
				}
				t.tui.SetProgDur(util.FormatProgDur(time.Duration(info.TrackProgress)), util.FormatProgDur(time.Duration(info.TrackLength)), info.Percentage)
				t.tui.SetTrackInfo(info.TrackInfo)
				t.tui.Draw()
			}
		case <-t.tui.Quit:
			return
		}
	}
}
