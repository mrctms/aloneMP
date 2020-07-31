package util

type TrackInfo struct {
	Title       string `json:"title"`
	Artist      string `json:"artist"`
	Album       string `json:"album"`
	AlbumArtist string `json:"albumArtist"`
	Composer    string `json:"composer"`
	Genre       string `json:"genre"`
	Year        int    `json:"year"`
}

type StatusResponse struct {
	TrackInfo *TrackInfo
	Progress  string `json:"progress"`
	Length    int    `json:"length"`
	Duration  string `json:"duration"`
	IsPlaying bool   `json:"isPlaying"`
	InError   bool   `json:"inError"`
}

type ServerMessage struct {
	Command string `json:"command"`
	Track   string `json:"track"`
	Source  string `json:"source"`
}

type FilePathInfo struct {
	FilePath string `json:"filePath"`
	IsDir    bool   `json:"isDir"`
}
