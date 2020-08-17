package util

type TrackInfo struct {
	Title  string `json:"title"`
	Artist string `json:"artist"`
	Album  string `json:"album"`
	Genre  string `json:"genre"`
	Year   int    `json:"year"`
}

type StatusResponse struct {
	TrackInfo              TrackInfo
	TrackProgress          int64  `json:"trackProgress"`
	Percentage             int    `json:"percentage"`
	TrackLength            int64  `json:"trackLength"`
	TrackProgressFormatted string `json:"trackProgressFormatted"`
	TrackLengthFormatted   string `json:"trackLengthFormatted"`
	IsPlaying              bool   `json:"isPlaying"`
	InError                bool   `json:"inError"`
}

type ServerMessage struct {
	Command      string `json:"command"`
	Track        string `json:"track"`
	Source       string `json:"source"`
	OutputDevice string `json:"device"`
}

type FilePathInfo struct {
	FilePath string `json:"filePath"`
	IsDir    bool   `json:"isDir"`
}

type PlayerArgs struct {
	Source       string
	OutputDevice string
}
