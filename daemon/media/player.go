package media

type Player interface {
	Init(args interface{}) error
	Play(track Track)
	Mute()
	Pause()
	VolumeUp()
	VolumeDown()
	Stop()
	Close()
	Info() *PlayerInformer
}
