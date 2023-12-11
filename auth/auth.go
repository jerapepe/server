package auth

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Config struct {
	MiSecreto string `json:"mi_secreto"`
}

func Cod(username string) string {
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	miSecret := GetToken()
	tokenString, err := token.SignedString([]byte(miSecret))
	if err != nil {
		fmt.Println("error 1")
	}
	return tokenString
}

func GetToken() string {
	configData, err := os.ReadFile("config.json")
	if err != nil {
		fmt.Println("Error leer archivo:", err)
	}
	var config Config
	err = json.Unmarshal(configData, &config)
	if err != nil {
		fmt.Println("Error parsear archivo:", err)
	}
	miSecret := config.MiSecreto
	return miSecret
}

func Decode(tokenString string) bool {
	miSecret := GetToken()
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Método de firma inválido: %v", token.Header["alg"])
		}
		return []byte(miSecret), nil
	})
	if err != nil {
		fmt.Println("error")
		return false
	}
	if token.Valid {
		claims, ok := token.Claims.(jwt.MapClaims)
		if ok && token.Valid {
			username := claims["username"].(string)
			fmt.Println(username)
			return true
		}
	} else {
		fmt.Println("token invalido")
		return false
	}
	return true
}
