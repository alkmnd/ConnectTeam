package game

import (
	"fmt"
	"log"
)

type Game struct {
	Name    string `json:"name"`
	clients map[*Client]bool
	//state
	MaxSize    int
	Status     string
	Creator    int `json:"creator_id"`
	Rounds     []Topic
	register   chan *Client
	unregister chan *Client
	broadcast  chan *Message
	ID         int `json:"id"`
}

type Topic struct {
	Id        int      `json:"id"`
	Title     string   `json:"title,omitempty"`
	Questions []string `json:"questions,omitempty"`
}

func NewGame(name string, id int, creator int, status string) *Game {
	return &Game{
		ID:      id,
		Name:    name,
		Rounds:  make([]Topic, 0),
		Creator: creator,
		Status:  status,
		MaxSize: 3,
		// state
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *Message),
	}
}

// Action
// state
// respondent

func (game *Game) GetId() int {
	return game.ID
}

func (game *Game) GetName() string {
	return game.Name
}

func (game *Game) RunGame() {
	for {
		select {
		case client := <-game.register:
			game.registerClientInGame(client)

		case client := <-game.unregister:
			game.unregisterClientInGame(client)

		case message := <-game.broadcast:
			game.broadcastToClientsInGame(message.encode())
		}
	}
}

const welcomeMessage = "%s joined the room"

func (game *Game) notifyClientJoined(client *Client) {
	message := &Message{
		Action:  JoinGameAction,
		Target:  game,
		Message: fmt.Sprintf(welcomeMessage, client.GetName()),
		Sender:  &client.User,
	}

	game.broadcastToClientsInGame(message.encode())
}

func (game *Game) registerClientInGame(client *Client) {
	log.Println("client joined")
	if len(game.clients) < game.MaxSize {
		game.clients[client] = true
		game.notifyClientJoined(client)
	}
	log.Println("registerClientInGame max number of users in game")
	return
}

func (game *Game) unregisterClientInGame(client *Client) {
	if _, ok := game.clients[client]; ok {
		delete(game.clients, client)
	}
}

func (game *Game) broadcastToClientsInGame(message []byte) {
	for client := range game.clients {
		log.Printf("broadcast message to client %s", client.GetName())
		client.send <- message
	}
}
