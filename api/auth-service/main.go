package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service struct {
	Db     *mongo.Client
	App    *fiber.App
	Config *Config
}

func (svc *Service) init() {
	app := svc.App.Group(svc.Config.BasePath)
	app.Post("/register", svc.Register)
	app.Post("/login", svc.Login)
	app.Get("/profile", JWTMiddleware(GetEnv("jwtSecret")), svc.Profile)
}

func main() {
	app := fiber.New()

	// Load config
	config, err := LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// MongoDB connection (add MongoDB setup here)
	db, err := connectMongo(config.MongoUri)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	service := &Service{
		Db:     db,
		App:    app,
		Config: config,
	}

	service.init()

	// Start the server
	log.Fatal(app.Listen(":" + config.Port))
}
