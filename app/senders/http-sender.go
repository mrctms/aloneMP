package senders

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"util"
)

type StatusResponse struct {
	TrackInfo trackInfo
	Progress  string `json:"progress"`
	Length    int    `json:"length"`
	Duration  string `json:"duration"`
	IsPlaying bool   `json:"isPlaying"`
	InError   bool   `json:"inError"`
}

type trackInfo struct {
	Title  string `json:"title"`
	Artist string `json:"artist"`
	Album  string `json:"album"`
}

type HttpSender struct {
	baseUrl string
	logger  *util.Logger
}

func NewHttpSender(address string) (*HttpSender, error) {

	hs := new(HttpSender)
	if !(strings.HasPrefix(address, "http") || strings.HasPrefix(address, "https")) {
		address = "http://" + address
	}
	u, err := url.Parse(address)
	if err != nil {
		return nil, err
	}
	hs.baseUrl = u.String()
	hs.logger, err = util.NewLogger("http-sender.log")
	if err != nil {
		return nil, err
	}
	return hs, nil
}

func (h *HttpSender) NextTrack() {
	req := fmt.Sprintf("%s?send=next", h.baseUrl)
	res, err := http.Get(req)
	if err != nil {
		h.logger.Write(fmt.Sprintf("error on next track: %v", err))
		return
	}
	h.logger.Write(fmt.Sprintf("send next track %s", res.Status))
}
func (h *HttpSender) Play(track interface{}) {
	trackToPlay, ok := track.(string)
	if ok {
		req := fmt.Sprintf("%s/command?send=play&track=%s", h.baseUrl, url.QueryEscape(trackToPlay))
		res, err := http.Get(req)
		if err != nil {
			h.logger.Write(fmt.Sprintf("error on play track: %v", err))
			return
		}
		h.logger.Write(fmt.Sprintf("send play track %s", res.Status))
	}
}

func (h *HttpSender) PreviousTrack() {
	req := fmt.Sprintf("%s/command?send=previous", h.baseUrl)
	res, err := http.Get(req)
	if err != nil {
		h.logger.Write(fmt.Sprintf("error on previous track: %v", err))
		return
	}
	h.logger.Write(fmt.Sprintf("send previous track %s", res.Status))
}

func (h *HttpSender) Mute() {
	req := fmt.Sprintf("%s/command?send=mute", h.baseUrl)
	res, err := http.Get(req)
	if err != nil {
		h.logger.Write(fmt.Sprintf("error on mute track: %v", err))
		return
	}
	h.logger.Write(fmt.Sprintf("send mute track %s", res.Status))
}

func (h *HttpSender) Pause() {
	req := fmt.Sprintf("%s/command?send=pause", h.baseUrl)
	res, err := http.Get(req)
	if err != nil {
		h.logger.Write(fmt.Sprintf("error on pause track: %v", err))
		return
	}
	h.logger.Write(fmt.Sprintf("send pause track %s", res.Status))
}

func (h *HttpSender) VolumeUp() {
	req := fmt.Sprintf("%s/command?send=up", h.baseUrl)
	res, err := http.Get(req)
	if err != nil {
		h.logger.Write(fmt.Sprintf("error on volume up: %v", err))
		return
	}
	h.logger.Write(fmt.Sprintf("send volume up %s", res.Status))
}

func (h *HttpSender) VolumeDown() {
	req := fmt.Sprintf("%s/command?send=down", h.baseUrl)
	res, err := http.Get(req)
	if err != nil {
		h.logger.Write(fmt.Sprintf("error on volume down: %v", err))
		return
	}
	h.logger.Write(fmt.Sprintf("send volume down %s", res.Status))
}

func (h *HttpSender) TrackInfo() interface{} {
	var info StatusResponse

	req := fmt.Sprintf("%s/status", h.baseUrl)
	res, err := http.Get(req)
	if err != nil {
		h.logger.Write(fmt.Sprintf("error on track info: %v", err))
	}
	if res != nil {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			h.logger.Write(fmt.Sprintf("error when read body: %v", err))
		}
		if body != nil {
			json.Unmarshal(body, &info)
		}

	}
	return info
}
func (h *HttpSender) ShutDown() {
	req := fmt.Sprintf("%s/command?send=shutdown", h.baseUrl)
	res, err := http.Get(req)
	if err != nil {
		h.logger.Write(fmt.Sprintf("error on shutdown: %v", err))
		return
	}
	h.logger.Write(fmt.Sprintf("send shutdown %s", res.Status))
}

func (h *HttpSender) Initialize(source string) {
	req := fmt.Sprintf("%s/command?send=init&source=%s", h.baseUrl, url.QueryEscape(source))
	res, err := http.Get(req)
	if err != nil {
		h.logger.Write(fmt.Sprintf("error on init: %v", err))
		return
	}
	h.logger.Write(fmt.Sprintf("send init %s", res.Status))
}

func (h *HttpSender) IsAlive() bool {
	_, err := http.Get(h.baseUrl)
	if err != nil {
		h.logger.Write(fmt.Sprintf("daemon is not alive: %v", err))
		return false
	}
	return true
}
