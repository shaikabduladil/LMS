package main

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	// "github.com/golang-jwt/jwt"

	// "github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	return string(bytes), err
}

func CheckHashPassword(reqPassword, userPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(reqPassword))

	return err == nil
}

func GetEnv(key string) string {
	return viper.GetString(key)
}

func GenerateJWT(email, secret string) (string, error) {
	claims := jwt.MapClaims{
		"email":  email,
		"expiry": time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}
