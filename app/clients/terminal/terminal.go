package terminal

import (
	"aloneMP/senders"
	"aloneMP/ui/tui"
	"fmt"
	"time"
)

type TerminalClient struct {
	mainTui *tui.MainTui
	sender  senders.Sender
}

func NewTerminalClient() *TerminalClient {
	terminalClient := new(TerminalClient)
	mainTui := tui.NewMainTui()
	terminalClient.mainTui = mainTui
	return terminalClient
}

func (t *TerminalClient) SetSender(sender senders.Sender) {
	t.sender = sender
}

func (t *TerminalClient) Run(rootDir string) {
	ticker := time.NewTicker(time.Second).C
	defer func() {
		t.mainTui.Stop()
		t.sender.ShutDown()
	}()
	go t.mainTui.Run()
	t.mainTui.PopolateTracksList(rootDir)
	t.sender.Initialize(rootDir)

	for {
		select {
		case track := <-t.mainTui.TrackSelected:
			t.sender.Play(track)
		case <-t.mainTui.Mute:
			t.sender.Mute()
		case <-t.mainTui.Paused:
			t.sender.Pause()
		case <-t.mainTui.VolumeUp:
			t.sender.VolumeUp()
		case <-t.mainTui.VolumeDown:
			t.sender.VolumeDown()
		case <-ticker:
			if !t.sender.IsAlive() {
				t.mainTui.Stop()
				fmt.Println("aloneMPd is not alive")
				return
			}
			info, ok := t.sender.TrackInfo().(senders.StatusResponse)
			if ok {
				if info.InError {
					t.mainTui.NextTrack()
					t.sender.Play(t.mainTui.CurrentTrack())
				} else if !info.IsPlaying {
					if (info.Duration == info.Progress) && (info.Length != 0) {
						t.mainTui.NextTrack()
						t.sender.Play(t.mainTui.CurrentTrack())
					}
				}
				t.mainTui.SetProgDur(info.Progress, info.Duration, info.Length)
				t.mainTui.SetTrackInfo(info.TrackInfo.Title, info.TrackInfo.Artist, info.TrackInfo.Album)
			}
			t.mainTui.Draw()
		case <-t.mainTui.Quit:
			return
		}
	}
}
