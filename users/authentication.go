package users

import (
	"Project/auth"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func AuthenticationUser(username string, password string) (bool, *User, string, error) {
	if username == "" && password == "" {
		fmt.Println("Esta vacio")
	} else {
		user, err := getUserFromDB(username)
		if err != nil {
			return false, nil, "", err
		}
		if user == nil {
			return false, nil, "", err
		}
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
			return false, nil, "", err
		}
		token := auth.Cod(user.Username)
		return true, user, token, nil
	}
	return false, nil, "", nil

}

func DecodeToken(token string) (User, error) {
	tokenS := strings.TrimPrefix(token, "Bearer ")
	username, err := auth.Decode(tokenS)
	if err != nil {
		fmt.Println("erros 1")
		return User{}, err
	}
	user, _, err := GetUser(username)
	if err != nil {
		fmt.Println("erros")
		return User{}, err
	}
	return user, nil
}
