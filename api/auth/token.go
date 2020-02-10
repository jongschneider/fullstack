package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var API_SECRET = os.Getenv("API_SECRET")

func CreateToken(user_id uint32) (string, error) {
	claims := jwt.MapClaims{
		"authorized": true,
		"user_id":    user_id,
		"exp":        time.Now().Add(time.Hour * 1).Unix(), // Token expires after 1 hour
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// return token.SignedString([]byte(os.Getenv("API_SECRET")))
	return token.SignedString([]byte(API_SECRET))
}

func TokenValid(r *http.Request) error {
	tokenString := ExtractToken(r)
	jwtParseFn := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(API_SECRET), nil
	}

	token, err := jwt.Parse(tokenString, jwtParseFn)
	if err != nil {
		return err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil
	}

	Pretty(claims)
	return nil
}

func Pretty(data interface{}) {
	b, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(string(b))
}

func ExtractToken(r *http.Request) string {
	keys := r.URL.Query()
	token := keys.Get("token")
	if token != "" {
		return token
	}

	bearerToken := r.Header.Get("Authorization")
	kv := strings.Split(bearerToken, " ")
	if len(kv) != 2 {
		return ""
	}
	return kv[1]

}

func ExtractTokenID(r *http.Request) (uint32, error) {
	tokenString := ExtractToken(r)
	jwtParseFn := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(API_SECRET), nil
	}

	token, err := jwt.Parse(tokenString, jwtParseFn)
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, nil
	}

	uid, err := strconv.ParseUint(fmt.Sprintf("%.0f", claims["user_id"]), 10, 32)
	if err != nil {
		return 0, err
	}
	return uint32(uid), nil
}
