package ui

type Interfacer interface {
	TrackList() []string
	CurrentTrack() string
	NextTrack() string
	PreviousTrack() string
}
