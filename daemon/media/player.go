package media

type Player interface {
	Start(source string)
	Clear()
	Close()
	Pause()
	Play()
	Mute()
	VolumeUp()
	VolumeDown()
	SetTrackToPlay(track string)
	Info() PlayerInformer
}
