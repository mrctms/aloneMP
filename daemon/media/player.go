package media

type Player interface {
	Start(source string)
	Close()
	Pause()
	Play()
	Mute()
	VolumeUp()
	VolumeDown()
	SetTrackToPlay(track string)
	Info() PlayerInformer
}
