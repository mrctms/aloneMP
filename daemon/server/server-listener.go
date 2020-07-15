package server

import "aloneMPd/media"

type ServerListener interface {
	Listen(address string) error
	NextTrack() chan bool
	PreviousTrack() chan bool
	PauseTrack() chan bool
	MuteTrack() chan bool
	VolumeUp() chan bool
	VolumeDown() chan bool
	ShutDown() chan bool
	SelectedTrack() chan string
	Source() chan string
	SetPlayerInfo(info media.PlayerInformer)
	PlayerInfo() media.PlayerInformer
}
