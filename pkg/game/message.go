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
const UserJoinedAction = "user-join"
const UserLeftAction = "user-left"

type Message struct {
	Action  string `json:"action"`
	Message string `json:"message"`
	Target  *Game  `json:"target"`
	Sender  *User  `json:"sender"`
	// time
}

func (message *Message) encode() []byte {
	json, err := json.Marshal(message)
	if err != nil {
		println("meow")
		log.Println(err)
	}

	return json
}