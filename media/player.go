package media

type Player interface {
	Start()
	Close()
	PlayingError() chan error
	Pause()
	Play()
	Mute()
	VolumeUp()
	VolumeDown()
	SetTrackToPlay(track string)
	Info() PlayerInformer
	Finished() chan bool
}
