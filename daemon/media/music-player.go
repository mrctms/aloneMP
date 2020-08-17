package media

import (
	"fmt"
	"os"
	"util"
)

type MusicPlayer struct {
	player       Player
	pause        chan bool
	trackToPlay  chan string
	mute         chan bool
	volumeUp     chan bool
	volumeDown   chan bool
	close        chan bool
	playerInfo   *PlayerInformer
	fatalError   chan error
	loadedTracks *[]Track
}

func NewMusicPlayer() *MusicPlayer {
	mp := new(MusicPlayer)
	mp.trackToPlay = make(chan string)
	mp.pause = make(chan bool)
	mp.mute = make(chan bool)
	mp.volumeUp = make(chan bool)
	mp.volumeDown = make(chan bool)
	mp.close = make(chan bool)
	mp.fatalError = make(chan error)
	return mp
}

func (m *MusicPlayer) Initialize(args util.PlayerArgs) error {
	stat, err := os.Stat(args.Source)
	if err != nil {
		return err
	}
	if stat.IsDir() {
		m.player, err = NewFilePlayer()
		if err != nil {
			return err
		}
		err = m.player.Init(nil)
		if err != nil {
			return err
		}
		m.playerInfo = m.player.Info()
		m.loadedTracks = new([]Track)
		files := util.GetKnowFiles(args.Source)
		for _, v := range files {
			track := new(FileTrack)
			track.setLocation(v)
			*m.loadedTracks = append(*m.loadedTracks, track)
		}
	} else {
		return fmt.Errorf("%v is not a directory", args.Source)
	}
	return nil
}

func (m *MusicPlayer) Start() error {
	defer func() {
		if r := recover(); r != nil {
			m.fatalError <- fmt.Errorf("Error while playing: %v", r)
		}
	}()
	for {
		select {
		case track := <-m.trackToPlay:
			for _, v := range *m.loadedTracks {
				if track == v.Location() {
					m.player.Stop()
					m.player.Play(v)
					break
				}
			}
		case <-m.pause:
			m.player.Pause()
		case <-m.mute:
			m.player.Mute()
		case <-m.volumeUp:
			m.player.VolumeUp()
		case <-m.volumeDown:
			m.player.VolumeDown()
		case <-m.close:
			m.player.Close()
		}
	}
}

func (m *MusicPlayer) Play(track string) {
	m.trackToPlay <- track
}

func (m *MusicPlayer) Mute() {
	m.mute <- true
}

func (m *MusicPlayer) Pause() {
	m.pause <- true
}

func (m *MusicPlayer) VolumeUp() {
	m.volumeUp <- true
}

func (m *MusicPlayer) VolumeDown() {
	m.volumeDown <- true
}

func (m *MusicPlayer) Close() {
	m.close <- true
}

func (m MusicPlayer) FatalError() chan error {
	return m.fatalError
}

func (m MusicPlayer) PlayerInfo() *PlayerInformer {
	return m.playerInfo
}
