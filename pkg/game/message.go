package game

import (
	"encoding/json"
	"log"
)

const JoinGameAction = "join-game"
const LeaveGameAction = "leave-game"
const SendMessageAction = "send-message"
const SelectTopicAction = "select-topic"
const StartGameAction = "start-game"
const Error = "error"
const UserJoinedAction = "join-success"

const UserLeftAction = "user-left"

type Message struct {
	Action  string          `json:"action"`
	Payload json.RawMessage `json:"payload,omitempty"`
	Target  *Game           `json:"target"`
	Sender  *User           `json:"sender"`
}

type MessageReceive struct {
	Action  string          `json:"action"`
	Message json.RawMessage `json:"message,omitempty"`
	Target  *Game           `json:"target"`
	Sender  string          `json:"sender"`
}

func (message *Message) encode() []byte {
	json, err := json.Marshal(message)
	if err != nil {
		println("123")
		log.Println(err)
	}

	return json
}
