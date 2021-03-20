package media

import (
	"os"
	"util"
)

type FileTrack struct {
	location string
	file     *os.File
	info     util.TrackInfo
}

func (f *FileTrack) setLocation(location string) {
	f.location = location

}
func (f *FileTrack) setTrackInfo(info util.TrackInfo) {
	f.info = info
}

func (f *FileTrack) setFile(file *os.File) {
	f.file = file
}

func (f FileTrack) Location() string {
	return f.location
}

func (f FileTrack) File() *os.File {
	return f.file
}

func (f FileTrack) TrackInfo() util.TrackInfo {
	return f.info
}
