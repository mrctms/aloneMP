package senders

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
	"util"
)

type TcpSender struct {
	logger *util.Logger
	conn   net.Conn
}

func NewTcpSender(address string) (*TcpSender, error) {
	var err error

	tcpSender := new(TcpSender)
	tcpSender.conn, err = net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	err = tcpSender.conn.(*net.TCPConn).SetKeepAlive(true)
	if err != nil {
		return nil, err
	}
	tcpSender.logger, err = util.NewLogger("tcp-sender.log")
	if err != nil {
		return nil, err
	}
	return tcpSender, nil
}

func (t *TcpSender) NextTrack() {
	msg := &util.ServerMessage{Command: util.NEXT_COMMAND}
	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		t.logger.Write(fmt.Sprintf("failed to send next track message %v", err))
		return
	}
	t.conn.Write(jsonMsg)
}

func (t *TcpSender) PreviousTrack() {
	msg := &util.ServerMessage{Command: util.PREVIOUS_COMMAND}
	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		t.logger.Write(fmt.Sprintf("failed to send previous track message %v", err))
		return
	}
	t.conn.Write(jsonMsg)
}

func (t *TcpSender) Mute() {
	msg := &util.ServerMessage{Command: util.MUTE_COMMAND}
	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		t.logger.Write(fmt.Sprintf("failed to send mute message %v", err))
		return
	}
	t.conn.Write(jsonMsg)
}

func (t *TcpSender) Pause() {
	msg := &util.ServerMessage{Command: util.PAUSE_COMMAND}
	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		t.logger.Write(fmt.Sprintf("failed to send pause message %v", err))
		return
	}
	t.conn.Write(jsonMsg)
}

func (t *TcpSender) VolumeUp() {
	msg := &util.ServerMessage{Command: util.VOLUME_UP_COMMAND}
	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		t.logger.Write(fmt.Sprintf("failed to send volume up message %v", err))
		return
	}
	t.conn.Write(jsonMsg)
}

func (t *TcpSender) VolumeDown() {
	msg := &util.ServerMessage{Command: util.VOLUME_DOWN_COMMAND}
	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		t.logger.Write(fmt.Sprintf("failed to send volume down message %v", err))
		return
	}
	t.conn.Write(jsonMsg)
}

func (t *TcpSender) Play(track string) {
	msg := &util.ServerMessage{Command: util.PLAY_COMMAND, Track: track}
	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		t.logger.Write(fmt.Sprintf("failed to send play track message %v", err))
		return
	}
	t.conn.Write(jsonMsg)
}

func (t *TcpSender) TrackInfo() *util.StatusResponse {
	msg := &util.ServerMessage{Command: "status"}
	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		t.logger.Write(fmt.Sprintf("failed to send status message %v", err))
		return nil
	}
	t.conn.Write(jsonMsg)

	decoder := json.NewDecoder(t.conn)

	var res util.StatusResponse

	decoder.Decode(&res)
	return &res
}

func (t *TcpSender) ShutDown() {
	msg := &util.ServerMessage{Command: util.SHUTDOWN_COMMAND}
	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		t.logger.Write(fmt.Sprintf("failed to send shutdown message %v", err))
		return
	}
	t.conn.Write(jsonMsg)
	t.conn.Close()
}

func (t *TcpSender) Initialize(source string) {
	msg := &util.ServerMessage{Command: util.INIT_COMMAND, Source: source}
	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		t.logger.Write(fmt.Sprintf("failed to send init message %v", err))
		return
	}
	t.conn.Write(jsonMsg)
}

func (t *TcpSender) IsAlive() bool {
	t.conn.SetDeadline(time.Now().Add(time.Second * 5))
	defer func() {
		var zero time.Time
		t.conn.SetDeadline(zero)
	}()
	msg := &util.ServerMessage{Command: "alive-check"}
	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		return false
	}
	t.conn.Write(jsonMsg)
	buf := make([]byte, 5)
	_, err = t.conn.Read(buf)
	if err != nil {
		return false
	}
	return true
}
