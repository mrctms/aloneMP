package senders

type Sender interface {
	NextTrack()
	PreviousTrack()
	Mute()
	Pause()
	VolumeUp()
	VolumeDown()
	Play(track string)
	TrackInfo() interface{}
	ShutDown()
	Initialize(source string)
	IsAlive() bool
}
