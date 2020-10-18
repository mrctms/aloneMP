package media

import (
	"os"
	"path/filepath"
	"util"
)

type PlayerInformer struct {
	volume                 float64
	playing                bool
	paused                 bool
	muted                  bool
	info                   util.TrackInfo
	percentage             int
	trackLength            int64
	trackLengthFormatted   string
	trackProgress          int64
	trackProgressFormatted string
	currentTrack           Track
	rootDir                string
	inError                bool
	trackList              *util.TrackListMessage
}

func (f *PlayerInformer) setVolume(volume float64) {
	f.volume = volume
}

func (f *PlayerInformer) setRootDir(rootDir string) {
	f.rootDir = rootDir
}

func (f *PlayerInformer) setPlaying(isPlaying bool) {
	f.playing = isPlaying
}

func (f *PlayerInformer) setPaused(isPaused bool) {
	f.paused = isPaused
}

func (f *PlayerInformer) setMuted(isMuted bool) {
	f.muted = isMuted
}

func (f *PlayerInformer) setTrackInfo(info util.TrackInfo) {
	f.info = info
}

func (f *PlayerInformer) setPercentage(percentage int) {
	f.percentage = percentage
}

func (f *PlayerInformer) setTrackLength(length int64) {
	f.trackLength = length
}

func (f *PlayerInformer) setTrackLengthFormatted(length string) {
	f.trackLengthFormatted = length
}

func (f *PlayerInformer) setTrackProgress(progress int64) {
	f.trackProgress = progress
}

func (f *PlayerInformer) setTrackProgressFormatted(progress string) {
	f.trackProgressFormatted = progress
}

func (f *PlayerInformer) setCurrentTrack(track Track) {
	f.currentTrack = track
}

func (f *PlayerInformer) setError(inError bool) {
	f.inError = inError
}

func (f PlayerInformer) CurrentVolume() float64 {
	return f.volume
}

func (f PlayerInformer) IsPlaying() bool {
	return f.playing
}

func (f PlayerInformer) IsPaused() bool {
	return f.paused
}

func (f PlayerInformer) IsMuted() bool {
	return f.muted
}

func (f PlayerInformer) TrackInfo() util.TrackInfo {
	return f.info
}

func (f PlayerInformer) Percentage() int {
	return f.percentage
}

func (f PlayerInformer) TrackLength() int64 {
	return f.trackLength
}
func (f PlayerInformer) TrackLengthFormatted() string {
	return f.trackLengthFormatted
}

func (f PlayerInformer) TrackProgressFormatted() string {
	return f.trackProgressFormatted
}

func (f PlayerInformer) TrackProgress() int64 {
	return f.trackProgress
}

func (f PlayerInformer) PlayingTrack() Track {
	return f.currentTrack
}

func (f PlayerInformer) TrackList() util.TrackListMessage {
	if f.trackList == nil {
		f.loadTrackList()
	}
	return *f.trackList
}

func (f *PlayerInformer) loadTrackList() {
	f.trackList = new(util.TrackListMessage)
	f.trackList.Root = f.rootDir
	f.trackList.Content = f.getRootInfo(f.rootDir)
}
func (f *PlayerInformer) getRootInfo(rootPath string) []util.RootInfo {
	var rootInfo []util.RootInfo
	var files = util.GetKnowFilesInfo(rootPath)
	for _, v := range files {
		var ri util.RootInfo
		path := filepath.Join(rootPath, v.Name())
		ri.Name = v.Name()
		ri.Path = ""
		p, err := os.Stat(path)
		if err != nil {
			continue
		}
		ri.IsDir = p.IsDir()
		if ri.IsDir {
			dirFiles := util.GetKnowFilesInfo(path)
			var rootInfoContent []util.RootInfo
			for _, e := range dirFiles {
				var ric util.RootInfo
				ric.Name = e.Name()
				subPath := filepath.Join(path, e.Name())
				ric.Path = subPath
				sp, err := os.Stat(subPath)
				if err != nil {
					continue
				}
				ric.IsDir = sp.IsDir()
				if ric.IsDir {
					ric.Path = ""
					ric.Content = f.getRootInfo(subPath)
				}
				rootInfoContent = append(rootInfoContent, ric)
			}
			ri.Content = rootInfoContent
		} else {
			ri.Path = path
		}
		rootInfo = append(rootInfo, ri)
	}
	return rootInfo
}
func (f PlayerInformer) InError() bool {
	return f.inError
}
