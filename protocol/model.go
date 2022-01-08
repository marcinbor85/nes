package protocol

import (
	"encoding/json"
)

type Frame struct {
	Message string `json:"message"`
	Hash	string `json:"hash"`
}

type Message struct {
	From		string `json:"from"`
	To		 	string `json:"to"`
	Timestamp	int64  `json:"timestamp"`
	Message		string `json:"message"`
}

func (msg *Message) String() string {
	msgJson, err := json.Marshal(msg)
	if err != nil {
		return ""
	}
	return string(msgJson)
}

func MessageFromString(str string) (*Message, error) {
	msg := &Message{}
	err := json.Unmarshal([]byte(str), &msg)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func (frm *Frame) String() string {
	frmJson, err := json.Marshal(frm)
	if err != nil {
		return ""
	}
	return string(frmJson)
}
