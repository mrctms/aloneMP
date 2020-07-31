package server

import (
	"aloneMPd/media"
	"encoding/json"
	"net"
	"util"
)

type TcpServer struct {
	nextTrack     chan bool
	previousTrack chan bool
	pause         chan bool
	mute          chan bool
	volumeUp      chan bool
	volumeDown    chan bool
	shutDown      chan bool
	selectedTrack chan string
	source        chan string
	pInfo         media.PlayerInformer
	listener      net.Listener
	conn          net.Conn
}

func NewTcpServer() *TcpServer {
	ts := new(TcpServer)
	ts.nextTrack = make(chan bool)
	ts.previousTrack = make(chan bool)
	ts.pause = make(chan bool)
	ts.mute = make(chan bool)
	ts.volumeUp = make(chan bool)
	ts.volumeDown = make(chan bool)
	ts.selectedTrack = make(chan string)
	ts.shutDown = make(chan bool)
	ts.source = make(chan string)
	return ts
}

func (t *TcpServer) SetPlayerInfo(info media.PlayerInformer) {
	t.pInfo = info
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
			t.nextTrack <- true
		case util.PREVIOUS_COMMAND:
			t.previousTrack <- true
		case util.PAUSE_COMMAND:
			t.pause <- true
		case util.MUTE_COMMAND:
			t.mute <- true
		case util.VOLUME_UP_COMMAND:
			t.volumeUp <- true
		case util.VOLUME_DOWN_COMMAND:
			t.volumeDown <- true
		case util.INIT_COMMAND:
			t.source <- msg.Source
		case util.SHUTDOWN_COMMAND:
			t.shutDown <- true
			err := t.startListen(address)
			if err != nil {
				return err
			}
		case "status":
			status := &util.StatusResponse{
				TrackInfo: t.pInfo.TrackInfo(),
				Progress:  t.pInfo.Progress(),
				Length:    t.pInfo.TrackLength(),
				Duration:  t.pInfo.Duration(),
				IsPlaying: t.pInfo.IsPlaying(),
				InError:   t.pInfo.InError(),
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

func (t *TcpServer) NextTrack() chan bool {
	return t.nextTrack
}

func (t *TcpServer) PreviousTrack() chan bool {
	return t.previousTrack
}

func (t *TcpServer) PauseTrack() chan bool {
	return t.pause
}

func (t *TcpServer) MuteTrack() chan bool {
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

func (t *TcpServer) SelectedTrack() chan string {
	return t.selectedTrack
}

func (t *TcpServer) Source() chan string {
	return t.source
}

func (t *TcpServer) PlayerInfo() media.PlayerInformer {
	return t.pInfo
}
