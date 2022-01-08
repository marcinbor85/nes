package protocol

import (
	"encoding/json"
)

type Frame struct {
	Cipherkey   string `json:"cipherkey"`
	Ciphertext	string `json:"ciphertext"`
	Signature	string `json:"signature"`
}

type Message struct {
	From		string `json:"from"`
	To		 	string `json:"to"`
	Timestamp	int64  `json:"timestamp"`
	Message		string `json:"message"`
}

func (message *Message) String() string {
	text, err := json.Marshal(message)
	if err != nil {
		return ""
	}
	return string(text)
}

func MessageFromString(text string) (*Message, error) {
	message := &Message{}
	err := json.Unmarshal([]byte(text), &message)
	if err != nil {
		return nil, err
	}
	return message, nil
}

func (frame *Frame) String() string {
	text, err := json.Marshal(frame)
	if err != nil {
		return ""
	}
	return string(text)
}

func FrameFromString(text string) (*Frame, error) {
	frame := &Frame{}
	err := json.Unmarshal([]byte(text), &frame)
	if err != nil {
		return nil, err
	}
	return frame, nil
}
