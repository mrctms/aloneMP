package server

import (
	"aloneMPd/media"
	"encoding/json"
	"net"
	"util"
)

type TcpServer struct {
	pause         chan bool
	mute          chan bool
	volumeUp      chan bool
	volumeDown    chan bool
	shutDown      chan bool
	selectedTrack chan string
	playerArgs    chan util.PlayerArgs
	playerInfo    *media.PlayerInformer
	listener      net.Listener
	conn          net.Conn
}

func NewTcpServer() *TcpServer {
	ts := new(TcpServer)
	ts.pause = make(chan bool)
	ts.mute = make(chan bool)
	ts.volumeUp = make(chan bool)
	ts.volumeDown = make(chan bool)
	ts.selectedTrack = make(chan string)
	ts.shutDown = make(chan bool)
	ts.playerArgs = make(chan util.PlayerArgs)
	return ts
}

func (t *TcpServer) SetPlayerInfo(info *media.PlayerInformer) {
	t.playerInfo = info
}

func (t *TcpServer) Listen(address string) error {
	err := t.startListen(address)
	if err != nil {
		return err
	}
	defer t.conn.Close()

	for {
		decoder := json.NewDecoder(t.conn)

		var msg util.ServerMessage

		decoder.Decode(&msg)

		switch msg.Command {
		case util.PLAY_COMMAND:
			t.selectedTrack <- msg.Track
		case util.NEXT_COMMAND:
		case util.PAUSE_COMMAND:
			t.pause <- true
		case util.MUTE_COMMAND:
			t.mute <- true
		case util.VOLUME_UP_COMMAND:
			t.volumeUp <- true
		case util.VOLUME_DOWN_COMMAND:
			t.volumeDown <- true
		case util.INIT_COMMAND:
			args := util.PlayerArgs{Source: msg.Source, OutputDevice: msg.OutputDevice}
			t.playerArgs <- args
		case util.SHUTDOWN_COMMAND:
			t.shutDown <- true
			err := t.startListen(address)
			if err != nil {
				return err
			}
		case "status":
			status := new(util.StatusResponse)
			if t.playerInfo != nil {
				status.TrackInfo = t.playerInfo.TrackInfo()
				status.TrackProgress = t.playerInfo.TrackProgress()
				status.Percentage = t.playerInfo.Percentage()
				status.TrackLength = t.playerInfo.TrackLength()
				status.TrackLengthFormatted = t.playerInfo.TrackLengthFormatted()
				status.TrackProgressFormatted = t.playerInfo.TrackProgressFormatted()
				status.IsPlaying = t.playerInfo.IsPlaying()
				status.InError = t.playerInfo.InError()
			}
			response, _ := json.Marshal(status)
			t.conn.Write(response)

		case "alive-check":
			t.conn.Write([]byte("4l0n3"))
		}
	}
}
func (t *TcpServer) startListen(address string) error {
	var err error
	t.listener, err = net.Listen("tcp", address)
	if err != nil {
		return err
	}
	t.conn, err = t.listener.Accept()
	if err != nil {
		return err
	}
	t.listener.Close()
	return nil
}

func (t *TcpServer) Pause() chan bool {
	return t.pause
}

func (t *TcpServer) Mute() chan bool {
	return t.mute
}

func (t *TcpServer) VolumeUp() chan bool {
	return t.volumeUp
}

func (t *TcpServer) VolumeDown() chan bool {
	return t.volumeDown
}

func (t *TcpServer) ShutDown() chan bool {
	return t.shutDown
}

func (t *TcpServer) Play() chan string {
	return t.selectedTrack
}

func (t *TcpServer) Initialize() chan util.PlayerArgs {
	return t.playerArgs
}
func (t *TcpServer) Close() {
	t.listener.Close()
	t.conn.Close()
}
