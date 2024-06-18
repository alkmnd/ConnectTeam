package notification_service

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
)

type Config struct {
	Host   string
	Path   string
	ApiKey string
}

type NotificationService struct {
	Conn *websocket.Conn
}

func NewNotificationService(cfg Config) (*NotificationService, error) {
	u := url.URL{Scheme: "ws", Host: cfg.Host, Path: cfg.Path, RawQuery: fmt.Sprintf("apiKey=%s", cfg.ApiKey)}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)

	if err != nil {
		log.Println("Dial:", err)
		return nil, err
	}
	return &NotificationService{
		Conn: conn,
	}, nil
}

//func (n *NotificationService) UpdateConnection() {
//	n.Conn
//}
