package senders

type Sender interface {
	NextTrack()
	PreviousTrack()
	Mute()
	Pause()
	VolumeUp()
	VolumeDown()
	Play(track interface{})
	TrackInfo() interface{}
	ShutDown()
	Initialize(source string)
	IsAlive() bool
}
