package utils

import (
	"backend_test/db"
	"fmt"
	"github.com/golang-jwt/jwt"
	"net/http"
	"strings"
	"time"
)

var tokenDuration = time.Hour * time.Duration(24)

func GetNewToken(username string, key string) map[string]string {
	var validUntil = time.Now().Add(tokenDuration)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username":   username,
		"validUntil": validUntil.Format(time.RFC3339),
	})

	tokenString, _ := token.SignedString([]byte(key))

	return map[string]string{
		"username":   username,
		"validUntil": validUntil.Format(time.RFC3339),
		"token":      tokenString,
	}
}

func ValidateToken(tokenString string, key string) bool {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(key), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims["username"], claims["validUntil"])
		return true
	} else {
		fmt.Println(err)

		return false
	}
}

func GetTokenPayload(tokenString string) jwt.MapClaims {
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return nil, nil
	})
	claims := token.Claims.(jwt.MapClaims)
	return claims
}

func HandleBearerToken(w http.ResponseWriter, r *http.Request) bool {

	reqToken := r.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer ")

	if len(splitToken) < 2 {
		ReturnUnauthorized(w, "Invalid token")
		return false
	}

	fmt.Sprintf("reqToken %v", reqToken)
	reqToken = splitToken[1]

	var claims = GetTokenPayload(reqToken)
	var username = claims["username"].(string)

	rows, err := db.Instance.Query(`SELECT "_user"."id", "_session"."signing_key", "_session"."created_timestamp" FROM "_user", "_session" where "_session"."user_id" = "_user"."id" and lower("_user"."username") = '` + username + `'`)
	db.CheckError(err)

	defer rows.Close()
	var userId int
	var signingKey string
	var createdTimestamp time.Time

	for rows.Next() {
		err = rows.Scan(&userId, &signingKey, &createdTimestamp)
		db.CheckError(err)
	}

	var validUntil = createdTimestamp.Add(tokenDuration)

	if userId == 0 {
		ReturnFailure(w, "User not found")
		return false
	}

	if validUntil.Before(time.Now()) {
		ReturnUnauthorized(w, "Expired token")
		return false
	}

	var valid = ValidateToken(reqToken, signingKey)

	if !valid {
		ReturnUnauthorized(w, "Invalid token")
		return false
	}

	return valid
}

func GetUsernameFromToken(r *http.Request) *string {

	reqToken := r.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer ")

	if len(splitToken) < 2 {
		return nil
	}

	reqToken = splitToken[1]

	var claims = GetTokenPayload(reqToken)
	var username = claims["username"].(string)
	return &username
}
