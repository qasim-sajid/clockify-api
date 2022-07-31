package auth

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/qasim-sajid/clockify-api/conf"
	"github.com/qasim-sajid/clockify-api/handler"
	"github.com/qasim-sajid/clockify-api/models"
)

// LoginForm defines login API input
type LoginForm struct {
	Identity string `json:"identity"`
	Password string `json:"password"`
}

// LoginResponse defines response on Login API
type LoginResponse struct {
	AuthToken    string `json:"auth_token"`
	RefreshToken string `json:"refresh_token"`
	UserID       string `json:"user_id"`
	Name         string `json:"name"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	Error        string `json:"error"`
}

// LoginUser handles user login requests
func LoginUser(h *handler.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		resp, status, err := handleLogin(c, h)
		if err != nil {
			c.JSON(status, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, resp)
	}
}

func handleLogin(c *gin.Context, h *handler.Handler) (response *LoginResponse, status int, err error) {
	login := &LoginForm{}
	response = &LoginResponse{}

	err = c.ShouldBindJSON(login)
	if err != nil {
		return response, http.StatusBadRequest, fmt.Errorf("invalid request format %v", err)
	}

	if login.Identity == "" || login.Password == "" {
		return response, http.StatusBadRequest, errors.New("credentials missing")
	}

	user, err := h.DB.CheckUserLogin(login.Identity, login.Password)
	if err != nil {
		if err.Error() == "invalid credentials" {
			return response, http.StatusUnauthorized, err
		}

		return response, http.StatusInternalServerError, fmt.Errorf("check login: %v", err)
	}

	response, err = getUserToken(user, h)
	if err != nil {
		return response, http.StatusInternalServerError, fmt.Errorf("get user token: %v", err)
	}

	return response, http.StatusOK, nil
}

// modify this if want to change login response
func getUserToken(user *models.User, h *handler.Handler) (response *LoginResponse, err error) {
	response = &LoginResponse{}
	response.UserID = user.ID
	response.Name = user.Name
	response.Username = user.Username
	response.Email = user.Email

	token, err := GenerateJWT(user)
	if err != nil {
		return response, err
	}

	refToken, err := GenerateRefreshJWT(user)
	if err != nil {
		return response, err
	}

	response.AuthToken = token
	response.RefreshToken = refToken

	return response, nil
}

// RefreshUserTokenPOST handles user refresh token request
func RefreshUserTokenPOST(h *handler.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		refToken := c.GetHeader("X-Refresh-Token")
		response := &LoginResponse{}

		if refToken == "" {
			response.Error = "Token Not Found."
			c.JSON(http.StatusUnauthorized, response)
			return
		}

		token, err := jwt.Parse(refToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				c.JSON(http.StatusUnauthorized, "Not Authorized")
				return nil, errors.New("Not Authorized")
			}
			return []byte(conf.Configs.RefreshSigningKey), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, "Not Authorized")
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		userID := fmt.Sprintf("%v", claims["user_id"])
		if userID == "" {
			c.JSON(http.StatusUnauthorized, "Not Authorized - No user ID")
			return
		}

		if claims["exp"] == nil {
			response.Error = "no expiry"
			c.JSON(http.StatusBadRequest, response)
			return
		}

		expTime, ok := claims["exp"].(float64)
		if !ok {
			response.Error = "invalid expiry"
			c.JSON(http.StatusBadRequest, response)
			return
		}

		if time.Now().Unix() > int64(expTime) {
			response.Error = "token expired"
			c.JSON(http.StatusConflict, response)
			return
		}

		user, err := h.DB.GetUser(userID)
		if err != nil {
			response.Error = "user not found"
			c.JSON(http.StatusUnauthorized, response)
			return
		}

		response, err = getUserToken(user, h)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, response)
	}
}
