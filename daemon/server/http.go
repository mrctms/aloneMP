package server

import (
	"aloneMPd/media"
	"encoding/json"
	"net/http"
)

type HttpServer struct {
	nextTrack     chan bool
	previousTrack chan bool
	pauseTrack    chan bool
	muteTrack     chan bool
	volumeUp      chan bool
	volumeDown    chan bool
	status        chan bool
	shutDown      chan bool
	selectedTrack chan string
	source        chan string
	playerInfo    media.PlayerInformer
}

func (h *HttpServer) Listen(address string) error {
	http.HandleFunc("/command", h.command)
	http.HandleFunc("/tracks", h.trackList)
	http.HandleFunc("/status", h.statusHandler)

	return http.ListenAndServe(address, nil)
}

func (h *HttpServer) NextTrack() chan bool {
	return h.nextTrack
}

func (h *HttpServer) PreviousTrack() chan bool {
	return h.previousTrack
}

func (h *HttpServer) PauseTrack() chan bool {
	return h.pauseTrack
}

func (h *HttpServer) MuteTrack() chan bool {
	return h.muteTrack
}

func (h *HttpServer) VolumeUp() chan bool {
	return h.volumeUp
}

func (h *HttpServer) VolumeDown() chan bool {
	return h.volumeDown
}

func (h *HttpServer) ShutDown() chan bool {
	return h.shutDown
}

func (h *HttpServer) SelectedTrack() chan string {
	return h.selectedTrack
}

func (h *HttpServer) Source() chan string {
	return h.source
}

func (h *HttpServer) SetPlayerInfo(info media.PlayerInformer) {
	h.playerInfo = info
}

func (h *HttpServer) PlayerInfo() media.PlayerInformer {
	return h.playerInfo
}

func NewHttpServer() *HttpServer {
	hs := new(HttpServer)
	hs.nextTrack = make(chan bool)
	hs.previousTrack = make(chan bool)
	hs.pauseTrack = make(chan bool)
	hs.muteTrack = make(chan bool)
	hs.volumeUp = make(chan bool)
	hs.volumeDown = make(chan bool)
	hs.status = make(chan bool)
	hs.selectedTrack = make(chan string)
	hs.shutDown = make(chan bool)
	hs.source = make(chan string)
	return hs
}

func (h *HttpServer) command(w http.ResponseWriter, req *http.Request) {

	var response map[string]string
	var commandToSend *string
	var track *string
	var source *string

	if req.Method == "GET" {
		w.Header().Add("Content-Type", "application/json")
		query := req.URL.Query()
		queryString := query.Get("send")
		commandToSend = &queryString
		trackToPlay := query.Get("track")
		track = &trackToPlay
		s := query.Get("source")
		source = &s
		response = make(map[string]string)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
	if response != nil && commandToSend != nil {
		response["send"] = *commandToSend
		switch *commandToSend {
		case "next":
			h.nextTrack <- true
		case "previous":
			h.previousTrack <- true
		case "mute":
			h.muteTrack <- true
		case "pause":
			h.pauseTrack <- true
		case "up":
			h.volumeUp <- true
		case "down":
			h.volumeDown <- true
		case "play":
			h.selectedTrack <- *track
		case "shutdown":
			h.shutDown <- true
		case "init":
			h.source <- *source
		default:
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		j, err := json.Marshal(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(j)
	}

}

func (h *HttpServer) trackList(w http.ResponseWriter, req *http.Request) {

	if h.playerInfo == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if req.Method == "GET" {
		list := h.playerInfo.TrackList()
		w.Header().Add("Content-Type", "application/json")
		response := map[string][]string{
			"trackList": list,
		}

		j, err := json.Marshal(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(j)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}

func (h *HttpServer) statusHandler(w http.ResponseWriter, req *http.Request) {

	if h.playerInfo == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if req.Method == "GET" {
		w.Header().Add("Content-Type", "application/json")
		response := map[string]interface{}{
			"isPlaying": h.playerInfo.IsPlaying(),
			"isPaused":  h.playerInfo.IsPaused(),
			"isMuted":   h.playerInfo.IsMuted(),
			"track":     h.playerInfo.PlayingTrack(),
			"progress":  h.playerInfo.Progress(),
			"length":    h.playerInfo.TrackLength(),
			"duration":  h.playerInfo.Duration(),
			"trackInfo": h.playerInfo.TrackInfo(),
			"inError":   h.playerInfo.InError(),
		}

		j, err := json.Marshal(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(j)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}
