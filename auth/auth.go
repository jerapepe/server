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
		"exp":      time.Now().Add(time.Minute).Unix(),
		//"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	configData, err := os.ReadFile("config.json")
	if err != nil {
		fmt.Println("Error al leer el archivo de configuración:", err)
		return
	}
	var config Config
	err = json.Unmarshal(configData, &config)
	if err != nil {
		fmt.Println("Error al parsear el archivo de configuración:", err)
		return
	}

	miSecret := config.MiSecreto

	tokenString, err := token.SignedString([]byte(miSecret))
	if err != nil {
		fmt.Println("error 1")
	}
	fmt.Println("Se codifico correctamente")
	decode(tokenString)
}

// decodificar
func decode(tokenString string) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Método de firma inválido: %v", token.Header["alg"])
		}
		return []byte("mi_clave_secreta"), nil
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
