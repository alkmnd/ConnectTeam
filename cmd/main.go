package main

import (
	connectteam "ConnectTeam"
	"ConnectTeam/pkg/handler"
	"ConnectTeam/pkg/repository"
	"ConnectTeam/pkg/repository/notification_service"
	"ConnectTeam/pkg/repository/payment_gateway"
	"ConnectTeam/pkg/repository/redis"
	"ConnectTeam/pkg/service"
	"ConnectTeam/pkg/service_handler"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"sync"
)

func main() {
	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing config: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading .env file: %s", err.Error())
	}

	logrus.Println("DB Username:", viper.GetString("db.username"))
	logrus.Println("DB Name:", viper.GetString("db.dbname"))
	logrus.Println("DB Port:", viper.GetString("db.port"))
	logrus.Println("DB Host:", viper.GetString("db.host"))
	logrus.Println("DB Password:", os.Getenv("DB_PASSWORD"))

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	})

	if err != nil {
		logrus.Fatalf("error connecting to the database: %s", err.Error())
	}

	rdb, err := redis.NewRedisClient(redis.Config{
		Host:     viper.GetString("redis.host"),
		Port:     viper.GetString("redis.port"),
		Password: os.Getenv("REDIS_PASSWORD"),
	})

	if err != nil {
		logrus.Fatalf("error connecting to Redis: %s", err.Error())
	}

	notificationService, err := notification_service.NewNotificationService(notification_service.Config{
		Host:   viper.GetString("notification_service.host"),
		Path:   viper.GetString("notification_service.path"),
		ApiKey: os.Getenv("NOTIFICATION_SERVICE_API_KEY"),
	})

	if err != nil {
		logrus.Fatalf("error initializing notification service: %s", err.Error())
	}

	yooClient := payment_gateway.NewYookassaClient(payment_gateway.Config{
		ShopId: viper.GetString("yookassa.shop_id"),
		ApiKey: os.Getenv("YOOCASSA_API_KEY"),
	})

	repos := repository.NewRepository(db, rdb, yooClient, notificationService)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	var wg sync.WaitGroup
	wg.Add(3)

	srv1 := new(connectteam.Server)
	// Запуск первого сервера в отдельной горутине.
	go func() {
		defer wg.Done()
		if err := srv1.Run(viper.GetString("port"), handlers.InitRoutes(viper.GetString("client_origin"))); err != nil {
			logrus.Fatalf("error %s", err.Error())
		}
	}()

	// Создание экземпляра сервера для второго сервера.
	srv2 := new(connectteam.Server)
	// Создание хэндлера для второго сервера.
	serviceHandler := service_handler.NewHandler(services, os.Getenv("SERVICE_API_KEY"))
	// Запуск второго сервера в отдельной горутине.
	go func() {
		defer wg.Done()
		if err := srv2.Run(viper.GetString("service_port"), serviceHandler.InitRoutes()); err != nil {
			logrus.Fatalf("error %s", err.Error())
		}
	}()

	wg.Wait()

}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
