package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

type (
	Data struct {
		Type        string    `json:"type"`
		SrcComId    string    `json:"srcComId"`
		DstComId    string    `json:"dstComId"`
		Diff        int       `json:"diff"`
		SubmittedAt time.Time `json:"submittedAt"`
	}
	Message struct {
		Code string `json:"messageCode"`
		Data Data   `json:"data,omitempty"`
	}
	ReceiveMsg struct {
		Code string `json:"code"`
		Data struct {
			Point     int       `json:"point"`
			SrcComId  string    `json:"srcComId"`
			UpdatedAt time.Time `json:"updatedAt"`
		} `json:"data"`
		Message     string `json:"message"`
		MessageCode string `json:"messageCode"`
	}
	Websocker struct {
		userId      string
		host        string
		wsSessionID string
		cnn         *websocket.Conn
		name        string
		token       string
	}
)

const (
	DataSelf = "self"
	Submit   = "point.submit"
	Ping     = "conn.ping"
	Changed  = "point.changed"
)

var (
	errInitCnn = errors.New("create connection failed")
)

func NewWebsocker(host string, userId string, name string, token string) *Websocker {

	return &Websocker{
		userId: userId,
		host:   host,
		name:   name,
		token:  token,
	}
}

func (w *Websocker) ConnectUntilSuccess() string {
	for {
		wsSessionId, err := getWsSessionId(w.host, w.userId, w.token)
		if err != nil {
			time.Sleep(200 * time.Millisecond)
			continue
		}
		return wsSessionId
	}
}

func (w *Websocker) Connect() error {
	w.wsSessionID = w.ConnectUntilSuccess()

	wsUrl := fmt.Sprintf("ws/v1/players/%s/ws/%s", w.userId, w.wsSessionID)

	u := url.URL{Scheme: "wss", Host: w.host, Path: wsUrl}

	cnn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return errInitCnn
	}
	w.cnn = cnn

	return nil
}
func (w *Websocker) ReadMessage() (messageType int, p []byte, err error) {
	return w.cnn.ReadMessage()

}
func (w *Websocker) Close() error {
	err := w.cnn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		return err
	}
	if err := w.cnn.Close(); err != nil {
		return err
	}
	return nil
}

func (w *Websocker) sendMessage(msg []byte) error {
	if err := w.cnn.WriteMessage(websocket.TextMessage, msg); err != nil {
		return fmt.Errorf("send message: %w", err)
	}
	return nil
}

func (w *Websocker) submitPoint(diff int) error {
	msg := Message{
		Code: Submit,
		Data: Data{
			Type:        DataSelf,
			SrcComId:    w.userId,
			DstComId:    w.userId,
			Diff:        diff,
			SubmittedAt: time.Now(),
		},
	}
	bytes, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal message: %w", err)
	}

	if err := w.sendMessage(bytes); err != nil {
		return fmt.Errorf("submit point: %w", err)
	}

	return nil
}

func (w *Websocker) ping() error {
	msg := Message{
		Code: Ping,
	}
	bytes, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal message: %w", err)
	}

	if err := w.sendMessage(bytes); err != nil {
		return fmt.Errorf("submit point: %w", err)
	}

	return nil
}
