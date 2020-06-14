package media

type PlayerInformer interface {
	CurrentVolume() float64
	IsPlaying() bool
	IsPaused() bool
	IsMuted() bool
	TrackInfo() interface{}
	PlayingTrack() string
	TrackLength() int
	Duration() string
	Progress() string
}
