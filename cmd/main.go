package main

import (
	connectteam "ConnectTeam"
	"ConnectTeam/pkg/game"
	"ConnectTeam/pkg/handler"
	"ConnectTeam/pkg/repository"
	"ConnectTeam/pkg/repository/filestorage"
	"ConnectTeam/pkg/service"
	"github.com/minio/minio-go"
	"log"
	"net/http"

	"os"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	if err := initConfig(); err != nil {
		logrus.Fatalf("error")
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error")
	}

	logrus.Println(viper.GetString("db.username"))
	logrus.Println(viper.GetString("db.name"))
	logrus.Println(viper.GetString("db.port"))
	logrus.Println(viper.GetString("db.host"))
	logrus.Println(os.Getenv("DB_PASSWORD"))
	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	})

	if err != nil {
		logrus.Println("1")
		logrus.Fatalf("error %s", err.Error())
	}

	accessKey := os.Getenv("ACCESS_KEY")
	secretKey := os.Getenv("SECRET_KEY")

	client, err := minio.New(viper.GetString("storage.endpoint"), accessKey, secretKey, false)
	if err != nil {
		log.Fatal(err)
	}

	fileStorage := filestorage.NewFileStorage(
		client,
		viper.GetString("storage.bucket"),
		viper.GetString("storage.endpoint"),
	)

	repos := repository.NewRepository(db)
	services := service.NewService(repos, fileStorage)
	handlers := handler.NewHandler(services)

	wsServer := game.NewWebsocketServer(repos, services)

	go wsServer.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		game.ServeWs(wsServer, w, r)
	})
	go http.ListenAndServe(":8080", nil)

	srv := new(connectteam.Server)
	if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
		logrus.Fatalf("error %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
