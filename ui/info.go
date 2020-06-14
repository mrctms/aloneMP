package ui

import (
	"gitlab.com/tslocum/cview"
)

type TrackInfo struct {
	*cview.TextView
}

func NewTrackInfo() *TrackInfo {
	trackInfo := new(TrackInfo)
	trackInfo.TextView = cview.NewTextView().SetTextAlign(cview.AlignLeft)
	trackInfo.TextView.SetTitle("Track Info")
	trackInfo.TextView.SetBorder(true)
	trackInfo.TextView.SetDynamicColors(true)
	return trackInfo
}

func (t *TrackInfo) SetInfo(info string) {
	t.TextView.SetText(info)
}
