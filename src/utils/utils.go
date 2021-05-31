package utils

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt"
)

func GetRsa() (rsa.PublicKey, error) {
	key := rsa.PublicKey{}
	keyURI := os.Getenv("WISH_PUBLIC_KEY_URI")
	resp, err := http.Get(keyURI)
	if err != nil {
		return key, fmt.Errorf("Unable to fetch key %v", err)
	}
	defer resp.Body.Close()
	var k Key
	err = json.NewDecoder(resp.Body).Decode(&k)
	if err != nil {
		return key, fmt.Errorf("Unable to json decode server response %v", err)
	}
	return k.Key, nil
}

// TestToken can be used to test if a token is valid and get the userID from it
func TestToken(t string) (int, error) {
	token, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		cert, err := GetRsa()
		if err != nil {
			return nil, err
		}

		return &cert, nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if _, ok := claims["exp"]; !ok {
			return 0, fmt.Errorf("Old tokens without expiration are no longer valid")
		}
		return int(claims["userID"].(float64)), nil
	}

	return 0, fmt.Errorf("Unable to parse JWT")
}
