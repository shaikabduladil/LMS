package main

import (
	"context"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Make sure the Config type is capitalized (exported)
type Config struct {
	Port      string `json:"port"`
	MongoUri  string `json:"mongoUri"`
	JWTSecret string `json:"jwtSecret"`
	BasePath  string `json:"basePath"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config") // without .json extension
	viper.SetConfigType("json")
	viper.AddConfigPath("configFiles") // folder name

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func connectMongo(uri string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return client, nil
}

func JWTMiddleware(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenStr := c.Get("Authorization")

		if tokenStr == "" {
			return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"error": "Missing Token"})
		}

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid Token"})
		}

		claims, ok := token.Claims.(jwt.MapClaims)

		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid Claims"})
		}

		email := claims["email"].(string)
		c.Locals("userEmail", email)

		return c.Next()
	}
}
