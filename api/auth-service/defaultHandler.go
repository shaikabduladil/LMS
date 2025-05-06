package main

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
)

func (svc *Service) Register(c *fiber.Ctx) error {
	var details User

	if err := c.BodyParser(&details); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if details.Name == "" || details.Password == "" || details.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Mandatory fields should not be Empty"})
	}

	hashedPassword, err := HashPassword(details.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error While Hashing Password"})
	}

	details.Password = hashedPassword
	details.CreatedAt = time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := svc.Db.Database(GetEnv("mongo.db")).Collection("users")

	result, err := collection.InsertOne(ctx, details)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User Created Sucessfully",
		"result":  result,
	})

}

func (svc *Service) Login(c *fiber.Ctx) error {

	type LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password`
	}

	var requestDetails LoginRequest

	err := c.BodyParser(&requestDetails)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	coll := svc.Db.Database(GetEnv("mongo.db")).Collection("users")

	err = coll.FindOne(ctx, fiber.Map{"email": requestDetails.Email}).Decode(&user)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	//compare Password

	if !CheckHashPassword(requestDetails.Password, user.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid Credentials"})
	}

	//Generate JWT Token
	token, err := GenerateJWT(user.Email, GetEnv("jwtSecret"))

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	expiresAt := time.Now().Add(72 * time.Hour)
	tokenDoc := Token{
		Email:     requestDetails.Email,
		CreatedAt: time.Now(),
		ExpiresAt: expiresAt,
		Token:     token,
	}

	TokenCollection := svc.Db.Database(GetEnv("mongo.db")).Collection("tokens")
	_, err = TokenCollection.InsertOne(ctx, tokenDoc)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error while inserting Token"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"token": token})

}

func (svc *Service) Profile(c *fiber.Ctx) error {
	email := c.Locals("userEmail").(string)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user User

	err := svc.Db.Database("lms").Collection("users").FindOne(ctx, fiber.Map{"email": email}).Decode(&user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error while Fetching User"})
	}

	//empty password while sending
	user.Password = ""
	return c.JSON(user)

}
