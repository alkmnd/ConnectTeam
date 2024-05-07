package main

import (
	connectteam "ConnectTeam"
	"ConnectTeam/pkg/handler"
	"ConnectTeam/pkg/repository"
	"ConnectTeam/pkg/repository/filestorage"
	"ConnectTeam/pkg/repository/payment_gateway"
	"ConnectTeam/pkg/repository/redis"
	"ConnectTeam/pkg/service"
	"ConnectTeam/pkg/service_handler"
	"github.com/minio/minio-go"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// @title           ConnectTeam API
// @version         1.0
// @description     API Server for ConnectTeam App

// @host      localhost:8000
// @BasePath  /

// @securityDefinitions.apiKey  apiKeyAuth
// @in header
// @name Authorization

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
		logrus.Fatalf("error %s", err.Error())
	}

	rdb, err := redis.NewRedisClient(redis.Config{
		Host:     viper.GetString("redis.host"),
		Port:     viper.GetString("redis.port"),
		Password: os.Getenv("REDIS_PASSWORD"),
	})

	if err != nil {
		logrus.Fatalf("error %s", err.Error())
	}

	accessKey := os.Getenv("ACCESS_KEY")
	secretKey := os.Getenv("SECRET_KEY")

	yooClient := payment_gateway.NewYookassaClient(payment_gateway.Config{
		ShopId: viper.GetString("yookassa.shop_id"),
		ApiKey: os.Getenv("INTEGRATION_API_KEY"),
	})
	client, err := minio.New(viper.GetString("storage.endpoint"), accessKey, secretKey, false)
	if err != nil {
		log.Fatal(err)
	}

	fileStorage := filestorage.NewFileStorage(
		client,
		viper.GetString("storage.bucket"),
		viper.GetString("storage.endpoint"),
	)

	repos := repository.NewRepository(db, rdb, yooClient)
	services := service.NewService(repos, fileStorage)
	handlers := handler.NewHandler(services)

	srv := new(connectteam.Server)
	if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
		logrus.Fatalf("error %s", err.Error())
	}

	serviceHandler := service_handler.NewHandler(services, os.Getenv("SERVICE_API_KEY"))
	if err := srv.Run(viper.GetString("service_port"), serviceHandler.InitRoutes()); err != nil {
		logrus.Fatalf("error %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
