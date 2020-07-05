package media

type PlayerInformer interface {
	CurrentVolume() float64
	IsPlaying() bool
	IsPaused() bool
	IsMuted() bool
	TrackInfo() *TrackInfo
	PlayingTrack() string
	TrackList() []string
	NextTrack() string
	PreviousTrack() string
	TrackLength() int
	Duration() string
	Progress() string
	InError() bool
}

type TrackInfo struct {
	Title  string `json:"title"`
	Artist string `json:"artist"`
	Album  string `json:"album"`
}
