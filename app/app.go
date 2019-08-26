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

	go pr.StartPlayer(tui)

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
					pr.Next <- true
					pr.Play <- true
				} else {
					pr.Play <- true
				}
			case "<Space>":
				pr.PaRes <- true
			case "m":
				pr.Mute <- true
			case "<Resize>":
				return // TODO
			}
		case <-ticker:
			tui.SetProgDur(pr.Progress, pr.Duration, pr.SongLenght)
			tui.RedrawAll()
			tui.UpdateInfo(pr.SongInfo)
		case <-pr.PlayingError:
			termui.Close()
			log.Fatalln(pr.ErrorMsg)
		}

	}
}
