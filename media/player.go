/*
Copyright (C) MarckTomack <marcktomack@tutanota.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.
*/

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

type Player struct {
	PaRes                 chan bool // pause/resume
	Play                  chan bool
	Mute                  chan bool
	VolumeUp              chan bool
	VolumeDown            chan bool
	PlayingError          chan error
	TrackToPlay           string
	terminateCurrentTrack bool
	IsPlaying             bool
	TrackInfo             tag.Metadata
	TrackLength           int
	Duration              string
	Progress              string
	Finished              chan bool
	ctrl                  *beep.Ctrl
	volume                *effects.Volume
	format                beep.Format
	streamer              beep.StreamSeekCloser
	currentVolume         float64
	isPaused              bool
	isSilence             bool
	fileToPlay            *os.File
}

func NewPlayer() *Player {
	return &Player{Play: make(chan bool), PaRes: make(chan bool), Mute: make(chan bool), VolumeUp: make(chan bool), VolumeDown: make(chan bool), PlayingError: make(chan error), Finished: make(chan bool)}
}

func (p *Player) loadStreamerAndFormat(file string) error {
	var err error

	p.fileToPlay, err = os.Open(file)
	if err != nil {
		return err
	}

	p.TrackInfo, err = tag.ReadFrom(p.fileToPlay)
	if err != nil {
		p.TrackInfo = nil
	}

	ex := filepath.Ext(p.fileToPlay.Name())
	switch ex {
	case ".mp3":
		p.streamer, p.format, err = mp3.Decode(p.fileToPlay)
	case ".wav":
		p.streamer, p.format, err = wav.Decode(p.fileToPlay)
	case ".flac":
		p.streamer, p.format, err = flac.Decode(p.fileToPlay)
	case ".ogg":
		p.streamer, p.format, err = vorbis.Decode(p.fileToPlay)
	}

	if err != nil {
		return err
	}

	return nil
}

func (p *Player) StartPlayer() {
	speaker.Init(initSimpleRate, initSimpleRate.N(time.Second/2))
	for {
		select {
		case <-p.Play:
			err := p.loadStreamerAndFormat(p.TrackToPlay)
			if err != nil {
				p.PlayingError <- fmt.Errorf("Failed to load the file %v", err)
			}
			res := beep.Resample(4, p.format.SampleRate, initSimpleRate, p.streamer)
			p.ctrl = &beep.Ctrl{Streamer: res, Paused: p.isPaused}
			p.volume = &effects.Volume{
				Streamer: p.ctrl,
				Base:     2,
				Volume:   p.currentVolume,
				Silent:   p.isSilence,
			}
			speaker.Play(p.volume)
			p.controlTrack()
		}
	}

}

func (p *Player) controlTrack() {
	p.IsPlaying = true
	p.terminateCurrentTrack = false
	defer p.fileToPlay.Close()
	for {
		select {
		case <-p.PaRes:
			speaker.Lock()
			p.ctrl.Paused = !p.ctrl.Paused
			speaker.Unlock()
		case <-p.Mute:
			speaker.Lock()
			p.volume.Silent = !p.volume.Silent
			speaker.Unlock()
		case <-p.VolumeUp:
			speaker.Lock()
			p.volume.Volume += 0.5
			speaker.Unlock()
		case <-p.VolumeDown:
			speaker.Lock()
			p.volume.Volume -= 0.5
			speaker.Unlock()
		case <-time.After(time.Second):
			if p.terminateCurrentTrack {
				p.isPaused = p.ctrl.Paused
				p.currentVolume = p.volume.Volume
				p.isSilence = p.volume.Silent
				return
			}
			if p.format.SampleRate != 0 {
				speaker.Lock()
				position := p.format.SampleRate.D(p.streamer.Position()).Round(time.Second)
				length := p.format.SampleRate.D(p.streamer.Len()).Round(time.Second)
				p.Duration = formatProgDur(length)
				p.Progress = formatProgDur(position)
				p.TrackLength = int(float64(position) / float64(length) * 100)
				speaker.Unlock()
				if position == length {
					p.IsPlaying = false
					p.isPaused = p.ctrl.Paused
					p.currentVolume = p.volume.Volume
					p.isSilence = p.volume.Silent
					p.Finished <- true
					return
				}
			}
		}
	}

}

func (p *Player) Close() {
	p.streamer.Close()
	p.format.SampleRate = 0
	speaker.Clear()
	p.terminateCurrentTrack = true
	p.IsPlaying = false
}

func formatProgDur(d time.Duration) string {
	// thanks to https://github.com/Depado
	h := math.Mod(d.Hours(), 24)
	m := math.Mod(d.Minutes(), 60)
	s := math.Mod(d.Seconds(), 60)
	tot := fmt.Sprintf("%02d:%02d:%02d", int(h), int(m), int(s))
	return tot
}
