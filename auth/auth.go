package auth

import (
	"encoding/json"
	"errors"
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

func Decode(tokenString string) (string, error) {
	miSecret := GetToken()
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de firma inválido: %v", token.Header["alg"])
		}
		return []byte(miSecret), nil
	})
	if err != nil {
		return "error tk", err
	}
	if !token.Valid {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("error al obtener los claims del token")
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return "", fmt.Errorf("reclamación 'exp' no válida en el token")
	}

	expirationTime := time.Unix(int64(exp), 0)
	if time.Now().After(expirationTime) {
		return "", errors.New("token expirado")
	}

	username, ok := claims["username"].(string)
	if !ok {
		return "", fmt.Errorf("username no encontrado en los claims")
	}
	return username, nil
}
