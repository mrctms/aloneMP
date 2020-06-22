package server

import (
	"aloneMP/media"
	"aloneMP/ui"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type command struct {
	Send string `json:"send"`
}

type HttpServer struct {
	NextTrack     chan bool
	PreviousTrack chan bool
	PauseTrack    chan bool
	MuteTrack     chan bool
	VolumeUp      chan bool
	VolumeDown    chan bool
	Status        chan bool
	PlayerInfo    media.PlayerInformer
	InterfaceInfo ui.Interfacer
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
	return hs
}
func (h *HttpServer) command(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	var response map[string]string
	var commandToSend *string

	if req.Method == "GET" {
		query := req.URL.Query()
		queryString := query.Get("send")
		commandToSend = &queryString
		response = make(map[string]string)
	} else if req.Method == "POST" {
		var cmd command
		response = make(map[string]string)
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}
		err = json.Unmarshal(body, &cmd)
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}
		commandToSend = &cmd.Send
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
		case "up":
			h.VolumeUp <- true
		case "down":
			h.VolumeDown <- true
		default:
			response["error"] = fmt.Sprintf("Unknow command %s", *commandToSend)
		}

		j, err := json.Marshal(response)
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}
		w.Write(j)
	}
}

func (h *HttpServer) trackList(w http.ResponseWriter, req *http.Request) {
	list := h.InterfaceInfo.TrackList()
	w.Header().Add("Content-Type", "application/json")
	response := map[string][]string{
		"trackList": list,
	}

	j, err := json.Marshal(response)
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}
	w.Write(j)

}

func (h *HttpServer) status(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	response := map[string]interface{}{
		"isPlaying": h.PlayerInfo.IsPlaying(),
		"isPaused":  h.PlayerInfo.IsPaused(),
		"isMuted":   h.PlayerInfo.IsMuted(),
		"track":     h.PlayerInfo.PlayingTrack(),
		"progress":  h.PlayerInfo.Progress(),
	}

	j, err := json.Marshal(response)
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}
	w.Write(j)
}

func (h *HttpServer) ListenAndServe(address string) {
	http.HandleFunc("/command", h.command)
	http.HandleFunc("/tracks", h.trackList)
	http.HandleFunc("/status", h.status)

	http.ListenAndServe(address, nil)
}
