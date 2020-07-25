package senders

import "util"

type Sender interface {
	NextTrack()
	PreviousTrack()
	Mute()
	Pause()
	VolumeUp()
	VolumeDown()
	Play(track string)
	TrackInfo() *util.StatusResponse
	ShutDown()
	Initialize(source string)
	IsAlive() bool
}
