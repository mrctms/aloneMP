package media

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"time"

	"github.com/dhowden/tag"
	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/flac"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/vorbis"
	"github.com/faiface/beep/wav"
)

var initSimpleRate = beep.SampleRate(44100)

type FilePlayer struct {
	pause                 chan bool
	play                  chan bool
	mute                  chan bool
	volumeUp              chan bool
	volumeDown            chan bool
	playingError          chan error
	trackToPlay           string
	terminateCurrentTrack bool
	ctrl                  *beep.Ctrl
	volume                *effects.Volume
	format                beep.Format
	streamer              beep.StreamSeekCloser
	fileToPlay            *os.File
	informer              *filePlayerInfo
	finished              chan bool
}

func NewFilePlayer() *FilePlayer {
	return &FilePlayer{informer: new(filePlayerInfo), play: make(chan bool), pause: make(chan bool), finished: make(chan bool), mute: make(chan bool), volumeUp: make(chan bool), volumeDown: make(chan bool), playingError: make(chan error)}
}

func (f *FilePlayer) Start() {
	speaker.Init(initSimpleRate, initSimpleRate.N(time.Second/2))
	for {
		select {
		case <-f.play:
			err := f.loadStreamerAndFormat(f.trackToPlay)
			if err != nil {
				f.playingError <- fmt.Errorf("Failed to load the file %v", err)
			}
			res := beep.Resample(4, f.format.SampleRate, initSimpleRate, f.streamer)
			f.ctrl = &beep.Ctrl{Streamer: res, Paused: f.informer.paused}
			f.volume = &effects.Volume{
				Streamer: f.ctrl,
				Base:     2,
				Volume:   f.informer.volume,
				Silent:   f.informer.muted,
			}
			f.informer.volume = f.volume.Volume
			f.informer.paused = f.ctrl.Paused
			f.informer.muted = f.volume.Silent
			f.informer.track = f.trackToPlay
			speaker.Play(f.volume)
			f.controlTrack()
		}
	}
}

func (f *FilePlayer) PlayingError() chan error {
	return f.playingError
}

func (f *FilePlayer) Pause() {
	f.pause <- true
}

func (f *FilePlayer) Play() {
	f.play <- true
}

func (f *FilePlayer) Mute() {
	f.mute <- true
}

func (f *FilePlayer) VolumeUp() {
	f.volumeUp <- true
}

func (f *FilePlayer) VolumeDown() {
	f.volumeDown <- true
}

func (f *FilePlayer) SetTrackToPlay(track string) {
	f.trackToPlay = track
}

func (f *FilePlayer) Info() PlayerInformer {
	return f.informer
}
func (f *FilePlayer) Finished() chan bool {
	return f.finished
}

func (f *FilePlayer) Close() {
	f.streamer.Close()
	f.format.SampleRate = 0
	speaker.Clear()
	f.terminateCurrentTrack = true
	f.informer.playing = false
}

func (f *FilePlayer) loadStreamerAndFormat(file string) error {
	var err error

	f.fileToPlay, err = os.Open(file)
	if err != nil {
		return err
	}

	trackInfo, err := tag.ReadFrom(f.fileToPlay)
	if err != nil {
		trackInfo = nil
	}
	f.informer.info = trackInfo
	ex := filepath.Ext(f.fileToPlay.Name())
	switch ex {
	case ".mp3":
		f.streamer, f.format, err = mp3.Decode(f.fileToPlay)
	case ".wav":
		f.streamer, f.format, err = wav.Decode(f.fileToPlay)
	case ".flac":
		f.streamer, f.format, err = flac.Decode(f.fileToPlay)
	case ".ogg":
		f.streamer, f.format, err = vorbis.Decode(f.fileToPlay)
	}

	if err != nil {
		return err
	}

	return nil
}

func (f *FilePlayer) controlTrack() {
	f.informer.playing = true
	f.terminateCurrentTrack = false
	defer f.fileToPlay.Close()
	for {
		select {
		case <-f.pause:
			speaker.Lock()
			f.ctrl.Paused = !f.ctrl.Paused
			f.informer.paused = f.ctrl.Paused
			speaker.Unlock()
		case <-f.mute:
			speaker.Lock()
			f.volume.Silent = !f.volume.Silent
			f.informer.muted = f.volume.Silent
			speaker.Unlock()
		case <-f.volumeUp:
			speaker.Lock()
			f.volume.Volume += 0.5
			speaker.Unlock()
		case <-f.volumeDown:
			speaker.Lock()
			f.volume.Volume -= 0.5
			speaker.Unlock()
		case <-time.After(time.Second):
			if f.terminateCurrentTrack {
				f.informer.paused = f.ctrl.Paused
				f.informer.muted = f.volume.Silent
				f.informer.volume = f.volume.Volume
				return
			}
			if f.format.SampleRate != 0 {
				speaker.Lock()
				position := f.format.SampleRate.D(f.streamer.Position()).Round(time.Second)
				length := f.format.SampleRate.D(f.streamer.Len()).Round(time.Second)
				f.informer.dur = formatProgDur(length)
				f.informer.prog = formatProgDur(position)
				f.informer.length = int(float64(position) / float64(length) * 100)
				speaker.Unlock()
				if position == length {
					f.informer.playing = false
					f.informer.paused = f.ctrl.Paused
					f.informer.muted = f.volume.Silent
					f.informer.volume = f.volume.Volume
					f.finished <- true
					return
				}
			}
		}
	}

}

func formatProgDur(d time.Duration) string {
	// thanks to https://github.com/Depado
	h := math.Mod(d.Hours(), 24)
	m := math.Mod(d.Minutes(), 60)
	s := math.Mod(d.Seconds(), 60)
	tot := fmt.Sprintf("%02d:%02d:%02d", int(h), int(m), int(s))
	return tot
}

type filePlayerInfo struct {
	volume  float64
	playing bool
	paused  bool
	muted   bool
	info    tag.Metadata
	length  int
	dur     string
	prog    string
	track   string
}

func (f *filePlayerInfo) CurrentVolume() float64 {
	return f.volume
}

func (f *filePlayerInfo) IsPlaying() bool {
	return f.playing
}

func (f *filePlayerInfo) IsPaused() bool {
	return f.paused
}

func (f *filePlayerInfo) IsMuted() bool {
	return f.muted
}

func (f *filePlayerInfo) TrackInfo() interface{} {
	return f.info
}

func (f *filePlayerInfo) TrackLength() int {
	return f.length
}

func (f *filePlayerInfo) Duration() string {
	return f.dur
}

func (f *filePlayerInfo) Progress() string {
	return f.prog
}
func (f *filePlayerInfo) PlayingTrack() string {
	return f.track
}
