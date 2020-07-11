package media

import "util"

type PlayerInformer interface {
	CurrentVolume() float64
	IsPlaying() bool
	IsPaused() bool
	IsMuted() bool
	TrackInfo() *util.TrackInfo
	PlayingTrack() string
	TrackList() []string
	NextTrack() string
	PreviousTrack() string
	TrackLength() int
	Duration() string
	Progress() string
	InError() bool
}
