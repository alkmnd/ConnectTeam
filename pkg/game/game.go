package game

import (
	"log"
)

type Game struct {
	Name       string `json:"name,omitempty"`
	clients    map[*Client]bool
	MaxSize    int     `json:"max_size,omitempty"`
	Status     string  `json:"status,omitempty"`
	Creator    int     `json:"creator_id,omitempty"`
	Rounds     []Topic `json:"rounds,omitempty"`
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
	log.Println("RunGame()")
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

type UserList struct {
	Users []User `json:"users"`
}

func (game *Game) listUsersInGame(client *Client) {
	for existingClient := range game.clients {
		message := &Message{
			Action: UserJoinedAction,
			Sender: &existingClient.User,
		}
		game.broadcastToClientsInGame(message.encode())
	}
}

func (game *Game) notifyClientJoined(client *Client) {
	//var users UserList
	//for i, _ := range game.clients {
	//	log.Println("notifyClientJoined")
	//	users.Users = append(users.Users, i.User)
	//}
	//bytes, _ := json.Marshal(users)

	message := &Message{
		Action:  JoinGameAction,
		Target:  game,
		Message: "",
		Sender:  &client.User,
	}

	game.broadcastToClientsInGame(message.encode())
}

func (game *Game) registerClientInGame(client *Client) {
	log.Println("client joined")
	//if game.Status == "in-progress" {
	//	log.Println("registerClientInGame max number of users in game")
	//	return
	//}
	if len(game.clients) < game.MaxSize {
		game.notifyClientJoined(client)
		game.clients[client] = true
		game.listUsersInGame(client)
		return
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
