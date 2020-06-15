package server

import (
	"aloneMP/media"
	"aloneMP/ui"
	"encoding/json"
	"io"
	"net/http"
)

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

func (h *HttpServer) next(w http.ResponseWriter, req *http.Request) {
	h.NextTrack <- true
	w.Header().Add("Content-Type", "application/json")
	response := map[string]bool{
		"next": true,
	}

	j, err := json.Marshal(response)
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}
	w.Write(j)
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

func (h *HttpServer) previous(w http.ResponseWriter, req *http.Request) {
	h.PreviousTrack <- true
	w.Header().Add("Content-Type", "application/json")
	response := map[string]bool{
		"previous": true,
	}

	j, err := json.Marshal(response)
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}
	w.Write(j)
}

func (h *HttpServer) pause(w http.ResponseWriter, req *http.Request) {
	h.PauseTrack <- true
	w.Header().Add("Content-Type", "application/json")
	response := map[string]bool{
		"pause": true,
	}
	j, err := json.Marshal(response)
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}
	w.Write(j)
}
func (h *HttpServer) mute(w http.ResponseWriter, req *http.Request) {
	h.MuteTrack <- true
	w.Header().Add("Content-Type", "application/json")
	response := map[string]bool{
		"mute": true,
	}

	j, err := json.Marshal(response)
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}
	w.Write(j)
}
func (h *HttpServer) volume(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	query := req.URL.Query()
	queryString := query.Get("go")
	if queryString == "up" {
		h.VolumeUp <- true
	}
	if queryString == "down" {
		h.VolumeDown <- true
	}
	response := map[string]string{
		"volume": queryString,
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
	http.HandleFunc("/next", h.next)
	http.HandleFunc("/tracks", h.trackList)
	http.HandleFunc("/previous", h.previous)
	http.HandleFunc("/pause", h.pause)
	http.HandleFunc("/mute", h.mute)
	http.HandleFunc("/volume", h.volume)
	http.HandleFunc("/status", h.status)

	http.ListenAndServe(address, nil)
}
