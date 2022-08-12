package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/qasim-sajid/clockify-api/conf"
	"github.com/qasim-sajid/clockify-api/handler"
	"github.com/qasim-sajid/clockify-api/models"
)

// IsUserAuthorized authorizes user account
func IsUserAuthorized(endpoint func(c *gin.Context, h *handler.Handler, origin *models.User), h *handler.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := verifyUserToken(c.GetHeader("Authorization"))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		} else {
			user, err := parseClaims(token)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			} else {
				if token.Valid {
					if c.Param("user_id") != "" && c.Param("user_id") != user.ID {
						c.JSON(http.StatusUnauthorized, gin.H{"error": "Not allowed to make this change"})
					} else {
						endpoint(c, h, user)
					}
				} else {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "Not Authorized"})
				}
			}
		}
	}
}

func verifyUserToken(authToken string) (*jwt.Token, error) {
	if authToken == "" {
		return nil, fmt.Errorf("Credentials missing!")
	}

	authToken = strings.TrimLeft(authToken, "Bearer")
	authToken = authToken[1:]

	token, err := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return http.StatusUnauthorized, errors.New("Not Authorized")
		}
		return []byte(conf.Signing_Key), nil
	})

	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Not Authorized")
	}

	return token, nil
}

func parseClaims(token *jwt.Token) (*models.User, error) {
	user := &models.User{}
	claims := token.Claims.(jwt.MapClaims)
	user.ID = fmt.Sprintf("%v", claims["user_id"])
	user.Username = fmt.Sprintf("%v", claims["username"])
	user.Email = fmt.Sprintf("%v", claims["email"])
	user.Name = fmt.Sprintf("%v", claims["name"])

	if claims["exp"] == nil {
		return nil, fmt.Errorf("Token expiry not found!")
	}

	expTime, ok := claims["exp"].(float64)
	if !ok {
		return nil, errors.New("Bad expiry!")
	}

	if time.Now().Unix() > int64(expTime) {
		return nil, errors.New("Token expired!")
	}

	return user, nil
}
