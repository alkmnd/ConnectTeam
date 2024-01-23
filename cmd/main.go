package main

import (
	connectteam "ConnectTeam"
	"ConnectTeam/pkg/handler"
	"ConnectTeam/pkg/repository"
	"ConnectTeam/pkg/service"
	"net/http"

	"os"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var upgrader = websocket.Upgrader{}

func main() {
	if err := initConfig(); err != nil {
		logrus.Fatalf("error")
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error")
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host: viper.GetString("db.host"), 
		Port:  viper.GetString("db.port"), 
		Username:  viper.GetString("db.username"),
		DBName:  viper.GetString("db.dbname"), 
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	})

	if err != nil {
		logrus.Fatalf("error %s", err.Error())
	}


	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)
	
	http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
        handlers.Echo(w, r)
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



