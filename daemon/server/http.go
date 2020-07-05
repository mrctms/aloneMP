package server

import (
	"aloneMPd/media"
	"encoding/json"
	"log"
	"net/http"
)

type command struct {
	Send  string `json:"send"`
	Track string `json:"track"`
}

type HttpServer struct {
	NextTrack     chan bool
	PreviousTrack chan bool
	PauseTrack    chan bool
	MuteTrack     chan bool
	VolumeUp      chan bool
	VolumeDown    chan bool
	Status        chan bool
	ShutDown      chan bool
	SelectedTrack chan string
	Source        chan string
	PlayerInfo    media.PlayerInformer
}

func NewHttpServer() *HttpServer {
	hs := new(HttpServer)
	hs.NextTrack = make(chan bool)
	hs.PreviousTrack = make(chan bool)
	hs.PauseTrack = make(chan bool)
	hs.MuteTrack = make(chan bool)
	hs.VolumeUp = make(chan bool)
	hs.VolumeDown = make(chan bool)
	hs.Status = make(chan bool)
	hs.SelectedTrack = make(chan string)
	hs.ShutDown = make(chan bool)
	hs.Source = make(chan string)
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
			h.NextTrack <- true
		case "previous":
			h.PreviousTrack <- true
		case "mute":
			h.MuteTrack <- true
		case "pause":
			h.PauseTrack <- true
		case "up":
			h.VolumeUp <- true
		case "down":
			h.VolumeDown <- true
		case "play":
			h.SelectedTrack <- *track
		case "shutdown":
			h.ShutDown <- true
		case "init":
			h.Source <- *source
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

	if h.PlayerInfo == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if req.Method == "GET" {
		list := h.PlayerInfo.TrackList()
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

func (h *HttpServer) status(w http.ResponseWriter, req *http.Request) {

	if h.PlayerInfo == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if req.Method == "GET" {
		w.Header().Add("Content-Type", "application/json")
		response := map[string]interface{}{
			"isPlaying": h.PlayerInfo.IsPlaying(),
			"isPaused":  h.PlayerInfo.IsPaused(),
			"isMuted":   h.PlayerInfo.IsMuted(),
			"track":     h.PlayerInfo.PlayingTrack(),
			"progress":  h.PlayerInfo.Progress(),
			"length":    h.PlayerInfo.TrackLength(),
			"duration":  h.PlayerInfo.Duration(),
			"trackInfo": h.PlayerInfo.TrackInfo(),
			"inError":   h.PlayerInfo.InError(),
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

func (h *HttpServer) ListenAndServe(address string) {
	http.HandleFunc("/command", h.command)
	http.HandleFunc("/tracks", h.trackList)
	http.HandleFunc("/status", h.status)

	log.Fatalln(http.ListenAndServe(address, nil))

}
