package senders

import (
	"encoding/json"
	"fmt"
	"net"
	"util"
)

type TcpSender struct {
	logger  *util.Logger
	tcpConn util.TcpConn
}

func NewTcpSender(address string) (*TcpSender, error) {
	var err error

	tcpSender := new(TcpSender)
	tcpSender.tcpConn.Conn, err = net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	err = tcpSender.tcpConn.Conn.(*net.TCPConn).SetKeepAlive(true)
	if err != nil {
		return nil, err
	}
	tcpSender.logger, err = util.NewLogger("tcp-sender.log")
	if err != nil {
		return nil, err
	}
	return tcpSender, nil
}

func (t *TcpSender) Mute() {
	msg := &util.ServerMessage{Command: util.MUTE_COMMAND}
	jsonMsg, _ := json.Marshal(msg)
	_, err := t.tcpConn.Write(jsonMsg)
	if err != nil {
		t.logger.Write(fmt.Sprintf("failed to send mute message %v", err))
	}
}

func (t *TcpSender) Pause() {
	msg := &util.ServerMessage{Command: util.PAUSE_COMMAND}
	jsonMsg, _ := json.Marshal(msg)
	_, err := t.tcpConn.Write(jsonMsg)
	if err != nil {
		t.logger.Write(fmt.Sprintf("failed to send pause message %v", err))
	}
}

func (t *TcpSender) VolumeUp() {
	msg := &util.ServerMessage{Command: util.VOLUME_UP_COMMAND}
	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		t.logger.Write(fmt.Sprintf("failed to send volume up message %v", err))
		return
	}
	t.tcpConn.Write(jsonMsg)
}

func (t *TcpSender) VolumeDown() {
	msg := &util.ServerMessage{Command: util.VOLUME_DOWN_COMMAND}
	jsonMsg, _ := json.Marshal(msg)
	_, err := t.tcpConn.Write(jsonMsg)
	if err != nil {
		t.logger.Write(fmt.Sprintf("failed to send volume down message %v", err))
	}
}

func (t *TcpSender) Play(track string) {
	msg := &util.ServerMessage{Command: util.PLAY_COMMAND, Track: track}
	jsonMsg, _ := json.Marshal(msg)
	_, err := t.tcpConn.Write(jsonMsg)
	if err != nil {
		t.logger.Write(fmt.Sprintf("failed to send play track message %v", err))
	}
}

func (t *TcpSender) TrackInfo() *util.StatusResponse {
	msg := &util.ServerMessage{Command: "status"}
	jsonMsg, _ := json.Marshal(msg)
	_, err := t.tcpConn.Write(jsonMsg)
	if err != nil {
		t.logger.Write(fmt.Sprintf("failed to send status message %v", err))
		return nil
	}

	decoder := json.NewDecoder(t.tcpConn.Conn)

	var res util.StatusResponse

	err = decoder.Decode(&res)
	if err != nil {
		t.logger.Write(fmt.Sprintf("failed to decode status message %v", err))
		return nil
	}
	return &res
}

func (t *TcpSender) ShutDown() {
	msg := &util.ServerMessage{Command: util.SHUTDOWN_COMMAND}
	jsonMsg, _ := json.Marshal(msg)

	_, err := t.tcpConn.Write(jsonMsg)
	if err != nil {
		t.logger.Write(fmt.Sprintf("failed to send shutdown message %v", err))
	}
	t.tcpConn.Close()
}

func (t *TcpSender) Initialize(source string) {
	msg := &util.ServerMessage{Command: util.INIT_COMMAND, Source: source, OutputDevice: ""}
	jsonMsg, _ := json.Marshal(msg)
	_, err := t.tcpConn.Write(jsonMsg)
	if err != nil {
		t.logger.Write(fmt.Sprintf("failed to send init message %v", err))
	}
}

func (t *TcpSender) IsAlive() bool {
	msg := &util.ServerMessage{Command: "alive-check"}
	jsonMsg, _ := json.Marshal(msg)
	_, err := t.tcpConn.Write(jsonMsg)
	if err != nil {
		return false
	}
	buf := make([]byte, 5)
	_, err = t.tcpConn.Read(buf)
	if err != nil {
		return false
	}
	return true
}
