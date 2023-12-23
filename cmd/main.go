package main

import (
	connectteam "ConnectTeam"
	"ConnectTeam/pkg/handler"
	"ConnectTeam/pkg/repository"
	"ConnectTeam/pkg/service"
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

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