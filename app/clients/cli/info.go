package cli

import (
	"gitlab.com/tslocum/cview"
)

type TrackInfo struct {
	*cview.TextView
}

func NewTrackInfo() *TrackInfo {
	trackInfo := new(TrackInfo)
	trackInfo.TextView = cview.NewTextView().SetTextAlign(cview.AlignCenter)
	//trackInfo.TextView.SetTitle("Track Info")
	trackInfo.TextView.SetBorder(false)
	trackInfo.TextView.SetDynamicColors(true)
	return trackInfo
}

func (t *TrackInfo) SetInfo(info string) {
	t.TextView.SetText(info)
}
