package auth

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/qasim-sajid/clockify-api/conf"
	"github.com/qasim-sajid/clockify-api/models"
)

// GenerateJWT generates JWT with payload of user info passed
func GenerateJWT(user *models.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["user_id"] = user.ID
	claims["username"] = user.Username
	claims["email"] = user.Email
	claims["name"] = user.Name
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString([]byte(conf.Configs.SigningKey))
	if err != nil {
		return "", fmt.Errorf("Unable to generate token: %s", err.Error())
	}

	return tokenString, nil
}

// GenerateRefreshJWT generates refresh token for user
func GenerateRefreshJWT(user *models.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["user_id"] = user.ID
	claims["exp"] = time.Now().Add(time.Hour * 6).Unix()
	tokenString, err := token.SignedString([]byte(conf.Configs.RefreshSigningKey))

	if err != nil {
		return "", fmt.Errorf("Unable to generate token: %s", err.Error())
	}

	return tokenString, nil
}
