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

package player

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
	"github.com/faiface/beep/wav"
	"github.com/marcktomack/aloneMP/ui"
)

type Player struct {
	PaRes        chan bool // pause/resume
	Play         chan bool
	Next         chan bool
	Mute         chan bool
	PlayingError chan error
	SongInfo     tag.Metadata
	SongLenght   int
	IsPlaying    bool
	Duration     string
	Progress     string
	ErrorMsg     string
	finished     bool
}

func NewPlayer() *Player {
	return new(Player)
}

func (p *Player) StartPlayer(rows *ui.Ui) {
	p.Play = make(chan bool)
	for {
		select {
		case <-p.Play:
			if p.finished {
				rows.SongsList.SelectedRow++
				if rows.SongsList.SelectedRow >= len(rows.SongsList.Rows) {
					rows.SongsList.SelectedRow = 0
				}
				p.playSong(rows.SongsList.Rows[rows.SongsList.SelectedRow])
			} else {
				p.playSong(rows.SongsList.Rows[rows.SongsList.SelectedRow])
			}
		}
	}
}

func (p *Player) playSong(file string) {

	var streamer beep.StreamSeekCloser
	var format beep.Format
	var decodeErr, tagErr error

	p.PlayingError = make(chan error)

	f, err := os.Open(file)
	if err != nil {
		p.PlayingError <- err
		p.ErrorMsg = fmt.Sprintf("%v", err)
		return
	}

	p.SongInfo, tagErr = tag.ReadFrom(f)
	if tagErr != nil {
		p.SongInfo = nil
	}

	ex := filepath.Ext(f.Name())
	switch ex {
	case ".mp3":
		streamer, format, decodeErr = mp3.Decode(f)
	case ".wav":
		streamer, format, decodeErr = wav.Decode(f)
	case ".flac":
		streamer, format, decodeErr = flac.Decode(f)
	}

	if decodeErr != nil {
		p.PlayingError <- decodeErr
		p.ErrorMsg = fmt.Sprintf("%v", decodeErr)
		return
	}

	defer streamer.Close()

	p.PaRes = make(chan bool)
	p.Next = make(chan bool)
	p.Mute = make(chan bool)

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/2))
	ctrl := &beep.Ctrl{Streamer: streamer, Paused: false}
	volume := &effects.Volume{
		Streamer: ctrl,
		Base:     2,
		Volume:   0,
		Silent:   false,
	}
	speaker.Play(volume)
	p.IsPlaying = true
	for {
		select {
		case <-p.PaRes:
			speaker.Lock()
			ctrl.Paused = !ctrl.Paused
			speaker.Unlock()
		case <-p.Next:
			p.IsPlaying = false
			return
		case <-p.Mute:
			speaker.Lock()
			volume.Silent = !volume.Silent
			speaker.Unlock()
		case <-time.After(time.Second):
			speaker.Lock()
			position := format.SampleRate.D(streamer.Position()).Round(time.Second)
			lenght := format.SampleRate.D(streamer.Len()).Round(time.Second)
			p.Duration = formatProgDur(lenght)
			p.Progress = formatProgDur(position)
			p.SongLenght = int(float64(position) / float64(lenght) * 100)
			speaker.Unlock()
			if position == lenght {
				p.finished = true
				p.IsPlaying = false
				return
			}

		}
	}
}

func formatProgDur(d time.Duration) string {
	h := math.Mod(d.Hours(), 24)
	m := math.Mod(d.Minutes(), 60)
	s := math.Mod(d.Seconds(), 60)
	tot := fmt.Sprintf("%02d:%02d:%02d", int(h), int(m), int(s))
	return tot
}
