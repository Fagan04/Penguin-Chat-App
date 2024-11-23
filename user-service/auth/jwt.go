package auth

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey []byte

func init() {
	config, err := os.ReadFile("config.json")
	if err != nil {
		log.Fatalf("Failed to read config.json: %v", err)
	}

	var data map[string]string
	err = json.Unmarshal(config, &data)
	if err != nil {
		log.Fatalf("Failed to parse config.json: %v", err)
	}

	key, exists := data["jwt_secret_key"]
	if !exists || key == "" {
		log.Fatal("jwt_secret_key is not set in config.json")
	}

	jwtKey = []byte(key)
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func GenerateJWT(username string, userID int) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour).Unix()
	claims := Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime,
			Id:        strconv.Itoa(userID), // Convert userID to string and assign it
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func ValidateJWT(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			switch {
			case ve.Errors&jwt.ValidationErrorMalformed != 0:
				return nil, errors.New("malformed token")
			case ve.Errors&jwt.ValidationErrorExpired != 0:
				return nil, errors.New("token expired")
			case ve.Errors&jwt.ValidationErrorNotValidYet != 0:
				return nil, errors.New("token not valid yet")
			default:
				return nil, errors.New("invalid token")
			}
		}
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
