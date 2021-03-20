package senders

import "util"

type Sender interface {
	Mute()
	Pause()
	VolumeUp()
	VolumeDown()
	Play(track string)
	TrackInfo() *util.StatusResponse
	TrackList() *util.TrackListMessage
	ShutDown()
	Initialize(source string)
	IsAlive() bool
}
