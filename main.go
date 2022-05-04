package main

import (
	"database/sql"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	log "github.com/sirupsen/logrus"
	"os"
	"weight-tracker/config"
	"weight-tracker/pkg/api"
	"weight-tracker/pkg/app"
	"weight-tracker/pkg/repository"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(os.Stderr, "this is the startup error: %s\n", err)
		os.Exit(1)
	}
}

func setupDatabase(connString string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func connectRedis(connectionString string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     connectionString,
		Password: "",
		DB:       0,
	})
	log.Info("Redis client worked")
	return client, nil
}

func run() error {
	pgHost := config.Get("DB_HOST")
	pgPort := config.Get("DB_PORT")
	pgUser := config.Get("DB_USER")
	pgPassword := config.Get("DB_PASSWORD")
	pgDB := config.Get("DB_NAME")
	pgSSL := config.Get("DB_SSL")
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", pgHost, pgPort, pgUser, pgPassword, pgDB, pgSSL)

	db, err := setupDatabase(connectionString)
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		log.Warning("no connection to the database")
	} else {
		log.Info("connected to the database")
	}
	// create redis store
	redisURL := fmt.Sprintf(config.Get("REDIS_URL"))
	redisClient, err := connectRedis(redisURL)
	if err != nil {
		log.Fatal("can  not connect to the redis store")
	}

	// create storage dependency
	storage := repository.NewStorage(db)
	redisStore := repository.NewRedisStore(redisClient)

	// create router dependency
	router := fiber.New()
	router.Use(cors.New())
	router.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))

	// create services
	userService := api.NewUserService(storage)
	weightService := api.NewWeightService(storage)
	tokenService := api.NewTokenService(redisStore)

	server := app.NewServer(router, userService, weightService, tokenService)

	err = server.Run()
	if err != nil {
		return err
	}
	return nil
}
