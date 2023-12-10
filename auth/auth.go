package auth

import (
	"Project/users"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Config struct {
	MiSecreto string `json:"mi_secreto"`
}

func Cod() {
	user := users.Users()
	claims := jwt.MapClaims{
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	miSecret := GetToken()
	tokenString, err := token.SignedString([]byte(miSecret))
	if err != nil {
		fmt.Println("error 1")
	}
	decode(tokenString)
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

// decodificar
func decode(tokenString string) {
	miSecret := GetToken()
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Método de firma inválido: %v", token.Header["alg"])
		}
		return []byte(miSecret), nil
	})

	if err != nil {
		fmt.Println("error")
	}

	if token.Valid {
		// El token es válido
		claims, ok := token.Claims.(jwt.MapClaims)
		if ok && token.Valid {
			username := claims["username"].(string)
			fmt.Println(username)
		}
	} else {
		fmt.Println("token invalido")
	}
}
