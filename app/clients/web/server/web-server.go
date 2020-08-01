package server

import (
	"aloneMP/senders"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"util"
)

type WebServer struct {
	address string
	sender  senders.Sender
	baseDir string
}

func NewWebServer(address string) *WebServer {
	webSrv := new(WebServer)
	webSrv.address = address
	return webSrv
}

func (a *WebServer) SetSender(sender senders.Sender) {
	a.sender = sender
}

func (a *WebServer) command(w http.ResponseWriter, req *http.Request) {
	response := make(map[string]string)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Content-Type", "application/json")

	if req.Method == "POST" {
		decoder := json.NewDecoder(req.Body)
		var msg util.ServerMessage
		decoder.Decode(&msg)

		response["send"] = msg.Command
		switch msg.Command {
		case util.NEXT_COMMAND:
			a.sender.NextTrack()
		case util.PREVIOUS_COMMAND:
			a.sender.PreviousTrack()
		case util.MUTE_COMMAND:
			a.sender.Mute()
		case util.PAUSE_COMMAND:
			a.sender.Pause()
		case util.VOLUME_UP_COMMAND:
			a.sender.VolumeUp()
		case util.VOLUME_DOWN_COMMAND:
			a.sender.VolumeDown()
		case util.PLAY_COMMAND:
			a.sender.Play(msg.Track)
		case util.SHUTDOWN_COMMAND:
			a.sender.ShutDown()
		default:
			response["error"] = "unknow command"
		}

	} else {
		response["error"] = "no GET request"
	}

	j, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(j)

}

func (a *WebServer) trackList(w http.ResponseWriter, req *http.Request) {
	response := make(map[string]string)
	filePathInfoList := new([]util.FilePathInfo)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Content-Type", "application/json")

	if req.Method == "GET" {
		a.popolateTrackList(a.baseDir, filePathInfoList)
	} else {
		response["error"] = "no GET request"
	}
	j, err := json.Marshal(filePathInfoList)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(j)
}

func (a *WebServer) popolateTrackList(rootDir string, slice *[]util.FilePathInfo) {
	files := util.GetKnowFilesInfo(rootDir)
	for _, v := range files {
		var fpi util.FilePathInfo
		path := filepath.Join(rootDir, v.Name())
		fpi.FilePath = path
		f, err := os.Stat(path)
		if err != nil {
			continue
		}
		fpi.IsDir = f.IsDir()
		*slice = append(*slice, fpi)
		if f.IsDir() {
			a.popolateTrackList(path, slice)
		}

	}
}

func (a *WebServer) Run(source string) {
	defer a.sender.ShutDown()
	a.baseDir = source
	a.sender.Initialize(source)
	http.HandleFunc("/tracks", a.trackList)
	http.HandleFunc("/command", a.command)

	http.ListenAndServe(a.address, nil)
}
