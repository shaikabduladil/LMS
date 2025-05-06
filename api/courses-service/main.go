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
	app.Post("/createcourse", svc.CreateCourse)
	app.Get("/getcourses", svc.GetCourses)
	app.Get("/course/:id/getcourse", svc.GetOneCourse)
	app.Put("/course/:id/update", svc.UpdateCourse)
	app.Delete("/course/:id/delete", svc.DeleteCourse)
}
func main() {
	app := fiber.New()
	config, err := LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config", err)
	}

	db, err := connectMongo(config.MongoUri)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB: ", err)
	}
	//intialize service,it initializes at original address
	service := &Service{
		Db:     db,
		App:    app,
		Config: config,
	}

	service.init()

	log.Fatal(app.Listen(":" + config.Port))
}
