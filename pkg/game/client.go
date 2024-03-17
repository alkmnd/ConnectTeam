package game

import (
	connectteam "ConnectTeam"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
}

type User struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Client struct {
	// The actual websocket connection.
	ID       uuid.UUID `json:"id"`
	conn     *websocket.Conn
	wsServer *WsServer
	send     chan []byte
	games    map[*Game]bool
	User     User `json:"user"`
}

func newClient(conn *websocket.Conn, wsServer *WsServer, user User) *Client {
	return &Client{
		ID:       uuid.New(),
		User:     user,
		conn:     conn,
		wsServer: wsServer,
		games:    make(map[*Game]bool),
		send:     make(chan []byte, 256),
	}
}

func (client *Client) GetName() string {
	return client.User.Name
}

func (client *Client) disconnect() {
	client.wsServer.unregister <- client
	for game := range client.games {
		game.unregister <- client
	}
}

// ServeWs handles websocket requests from clients requests.
func ServeWs(wsServer *WsServer, w http.ResponseWriter, r *http.Request) {
	name, ok := r.URL.Query()["name"]

	if !ok || len(name[0]) < 1 {
		log.Println("Url Param 'name' is missing")
		return
	}

	id_string, ok := r.URL.Query()["id"]

	if !ok || len(id_string[0]) < 1 {
		log.Println("Url Param 'id' is missing")
		return
	}

	id, err := strconv.Atoi(id_string[0])
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := newClient(conn, wsServer, User{
		Id:   id,
		Name: name[0],
	})

	go client.writePump()
	go client.readPump()

	wsServer.register <- client
}

const (
	// Max wait time when writing message to peer
	writeWait = 10 * time.Second

	// Max time till next pong from peer
	pongWait = 60 * time.Second

	// Send ping interval, must be less than pong wait time
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 10000
)

func (client *Client) readPump() {
	defer func() {
		client.disconnect()
	}()

	client.conn.SetReadLimit(maxMessageSize)
	client.conn.SetReadDeadline(time.Now().Add(pongWait))
	client.conn.SetPongHandler(func(string) error { client.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	// Start endless read loop, waiting for messages from client
	for {
		_, jsonMessage, err := client.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("unexpected close error: %v", err)
			}
			break
		}

		client.handleNewMessage(jsonMessage)
	}
}

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

func (client *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		client.conn.Close()
	}()
	for {
		select {
		case message, ok := <-client.send:
			client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The WsServer closed the channel.
				client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := client.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Attach queued chat messages to the current websocket message.
			n := len(client.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-client.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
func (client *Client) handleNewMessage(jsonMessage []byte) {

	var message Message
	if err := json.Unmarshal(jsonMessage, &message); err != nil {
		log.Printf("Error on unmarshal JSON message %s", err)
		return
	}

	// Attach the client object as the sender of the messsage.
	message.Sender = &client.User

	switch message.Action {

	case SendMessageAction:
		if game := client.wsServer.findGame(message.Target.ID); game != nil {
			game.broadcast <- &message
		}
	case JoinGameAction:
		gameId := message.Target.ID

		game := client.wsServer.findGame(gameId)
		if game == nil {
			return
		}
		log.Println("joinGameAction")
		client.handleJoinGameMessage(message)
	case StartGameAction:
		log.Println("startGameAction")
		client.handleStartGameMessage(message)
	case LeaveGameAction:
		client.handleLeaveGameMessage(message)

	case SelectTopicAction:
		log.Println("selectTopicAction")
		client.handleSelectTopicGameMessage(message)
	}
}

func (client *Client) handleStartGameMessage(message Message) {
	//  меняем статус,
	gameId := message.Target.ID

	game := client.wsServer.findGame(gameId)
	game.Status = "in_progress"

	err := client.wsServer.repos.StartGame(gameId)
	if err != nil {
		log.Println("handleStartGameMessage unknown game")
		return
	}
	if len(game.Rounds) == 0 {
		var messageError Message
		messageError.Action = Error
		messageError.Target = message.Target
		messageError.Message = "number of rounds is 0"
		client.send <- messageError.encode()
		log.Println("number of rounds is 0")
		return
	}
	questions := map[int][]connectteam.Question{}

	for i, _ := range game.Rounds {
		questions[game.Rounds[i].Id], _ = client.wsServer.repos.Question.GetAllWithLimit(game.Rounds[i].Id, len(game.clients))
		game.Rounds[i].Questions = make([]string, len(game.clients))
		for j := 0; j < len(game.clients); j++ {
			game.Rounds[i].Questions[j] = questions[game.Rounds[i].Id][j].Content
		}
	}

	bytes, err := json.Marshal(game.Rounds)
	if err != nil {
		return
	}
	message.Message = string(bytes)
	game.broadcast <- &message
}

func (client *Client) handleSelectTopicGameMessage(message Message) {
	gameId := message.Target.ID

	game := client.wsServer.findGame(gameId)

	if message.Sender.Id != game.Creator {
		log.Println("handleSelectTopicGameMessage message.Sender.Id != game.Creator")
		return
	}

	if err := json.Unmarshal([]byte(message.Message), &game.Rounds); err != nil {
		log.Println("handleSelectTopicGameMessage unmarshal error")
		return
	}
}

func (client *Client) handleJoinGameMessage(message Message) {
	//gameName := message.Message
	gameId := message.Target.ID

	game := client.wsServer.findGame(gameId)
	//if game == nil {
	//	// тут из бд создаем
	//	game = client.wsServer.createGame(gameName, gameId)
	//}

	client.games[game] = true

	game.register <- client
}

func (client *Client) handleLeaveGameMessage(message Message) {
	game := client.wsServer.findGame(message.Target.ID)
	if _, ok := client.games[game]; ok {
		delete(client.games, game)
	}

	game.unregister <- client
}
