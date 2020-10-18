package server

import (
	"aloneMPd/media"
	"encoding/json"
	"io"
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
	err           chan error
	playerArgs    chan util.PlayerArgs
	playerInfo    *media.PlayerInformer
	listener      net.Listener
	tcpConn       util.TcpConn
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

func (t *TcpServer) Listen(address string) {
	err := t.startListen(address)
	if err != nil {
		t.err <- err
	}
	defer t.tcpConn.Close()

	for {
		decoder := json.NewDecoder(t.tcpConn.Conn)

		var msg util.ServerMessage

		err := decoder.Decode(&msg)
		if err != nil {
			if err == io.EOF { // connection is closed. The method to check if the connection from the client is closed should be better
				t.shutDown <- true
				err := t.startListen(address)
				if err != nil {
					t.err <- err
				}
			} else {
				t.err <- err
			}
		}

		switch msg.Command {
		case util.PLAY_COMMAND:
			t.selectedTrack <- msg.Track
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
				t.err <- err
			}
		case util.STATUS:
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
			t.tcpConn.Write(response)
		case util.TRACK_LIST:
			response, _ := json.Marshal(t.playerInfo.TrackList())
			t.tcpConn.Write(response)
		case util.ALIVE_CHECK:
			t.tcpConn.Write([]byte("4l0n3"))
		}
	}
}
func (t *TcpServer) startListen(address string) error {
	var err error
	t.listener, err = net.Listen("tcp", address)
	if err != nil {
		return err
	}
	t.tcpConn.Conn, err = t.listener.Accept()
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

func (t *TcpServer) FatalError() chan error {
	return t.err
}

func (t *TcpServer) Initialize() chan util.PlayerArgs {
	return t.playerArgs
}
func (t *TcpServer) Close() {
	t.listener.Close()
	t.tcpConn.Close()
}
