package media

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
	"util"

	"github.com/dhowden/tag"
	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/flac"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/vorbis"
	"github.com/faiface/beep/wav"
)

const initSimpleRate = beep.SampleRate(44100)

type FilePlayer struct {
	ctrl        *beep.Ctrl
	volume      *effects.Volume
	format      beep.Format
	streamer    beep.StreamSeekCloser
	logger      *util.Logger
	playerInfo  *PlayerInformer
	infoTimer   *time.Timer
	inError     bool
	trackInPlay *FileTrack
}

func NewFilePlayer() (*FilePlayer, error) {
	var err error
	fp := new(FilePlayer)
	fp.ctrl = new(beep.Ctrl)
	fp.volume = new(effects.Volume)
	fp.volume.Base = 2
	fp.volume.Streamer = fp.ctrl
	fp.playerInfo = new(PlayerInformer)
	fp.logger, err = util.NewLogger("file-player.log")
	if err != nil {
		return nil, err
	}
	return fp, nil

}

func (f *FilePlayer) Init(args interface{}) error {
	go f.runInfoTimer()
	return speaker.Init(initSimpleRate, initSimpleRate.N(time.Second/2))
}

func (f *FilePlayer) Play(track Track) {
	if f.infoTimer != nil {
		defer f.infoTimer.Reset(time.Millisecond)
		f.infoTimer.Stop()
	}

	t, ok := track.(*FileTrack)
	if !ok {
		f.inError = true
		return
	}

	f.trackInPlay = t
	err := f.loadFile(t.Location())
	if err != nil {
		f.inError = true
		return
	}
	res := beep.Resample(4, f.format.SampleRate, initSimpleRate, f.streamer)
	f.ctrl.Streamer = res
	speaker.Play(f.volume)
	f.inError = false

}

func (f *FilePlayer) runInfoTimer() {
	f.infoTimer = time.NewTimer(time.Millisecond)
	defer func() {
		if r := recover(); r != nil {
			f.logger.Write(fmt.Sprintf("error while read track info: %v", r))
			go f.runInfoTimer()
		}
	}()
	for {
		select {
		case <-f.infoTimer.C:
			f.updateInfo()
			f.infoTimer.Reset(time.Millisecond)
		}
	}
}

func (f *FilePlayer) Pause() {
	speaker.Lock()
	f.ctrl.Paused = !f.ctrl.Paused
	speaker.Unlock()
}

func (f *FilePlayer) Mute() {
	speaker.Lock()
	f.volume.Silent = !f.volume.Silent
	speaker.Unlock()
}

func (f *FilePlayer) VolumeUp() {
	speaker.Lock()
	if !(f.volume.Volume == 0) {
		f.volume.Volume += 0.5
	}
	speaker.Unlock()
}

func (f *FilePlayer) VolumeDown() {
	speaker.Lock()
	f.volume.Volume -= 0.5
	speaker.Unlock()
}

func (f *FilePlayer) updateInfo() {
	f.playerInfo.setError(f.inError)
	if f.trackInPlay != nil {
		position := f.format.SampleRate.D(f.streamer.Position()).Round(time.Second)
		length := f.format.SampleRate.D(f.streamer.Len()).Round(time.Second)
		f.playerInfo.setPaused(f.ctrl.Paused)
		f.playerInfo.setMuted(f.volume.Silent)
		f.playerInfo.setPercentage(int(float64(position) / float64(length) * 100))
		f.playerInfo.setTrackProgress(int64(position))
		f.playerInfo.setTrackLength(int64(length))
		f.playerInfo.setTrackInfo(f.trackInPlay.TrackInfo())
		f.playerInfo.setCurrentTrack(f.trackInPlay)
		f.playerInfo.setPlaying(position != length)
	}
}

func (f FilePlayer) Info() *PlayerInformer {
	return f.playerInfo
}

func (f *FilePlayer) loadFile(filePath string) error {

	fileToPlay, err := os.Open(filePath)
	if err != nil {
		return err
	}
	f.trackInPlay.setFile(fileToPlay)

	metaData, _ := tag.ReadFrom(fileToPlay)
	var trackInfo util.TrackInfo
	if metaData != nil {
		trackInfo.Title = metaData.Title()
		trackInfo.Artist = metaData.Artist()
		trackInfo.Album = metaData.Album()
		trackInfo.Genre = metaData.Genre()
		trackInfo.Year = metaData.Year()
	}
	f.trackInPlay.setTrackInfo(trackInfo)

	ex := filepath.Ext(fileToPlay.Name())
	switch ex {
	case ".mp3":
		f.streamer, f.format, err = mp3.Decode(fileToPlay)
	case ".wav":
		f.streamer, f.format, err = wav.Decode(fileToPlay)
	case ".flac":
		f.streamer, f.format, err = flac.Decode(fileToPlay)
	case ".ogg":
		f.streamer, f.format, err = vorbis.Decode(fileToPlay)
	}

	if err != nil {
		return err
	}

	return nil
}

func (f *FilePlayer) Stop() {
	if f.trackInPlay != nil {
		f.trackInPlay.File().Close()
		f.trackInPlay = nil
	}
	speaker.Clear()

}

func (f *FilePlayer) Close() {
	if f.trackInPlay != nil {
		f.trackInPlay.File().Close()
		f.trackInPlay = nil
	}
	if f.infoTimer != nil {
		f.infoTimer.Stop()
	}
	speaker.Close()
}
