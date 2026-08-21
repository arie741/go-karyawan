package auth

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Handler struct {
	repo *Repository
}

func NewHandler(repo *Repository) *Handler {
	return &Handler{repo}
}

func (h *Handler) LoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{})
}

func (h *Handler) Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	t := jwt.New(jwt.GetSigningMethod("HS256"))

	t.Claims = &UserClaims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		},
		"level1",
		UserInfo{username},
	}

	secret := os.Getenv("SECRET")

	authToken, err := t.SignedString([]byte(secret))
	if err != nil {
		log.Fatal(err)
	}

	var authorizedUser string
	users, userErr := h.repo.GetAllUsers()
	if userErr != nil {
		panic(userErr)
	}

	for _, user := range users {
		if user.Name == username && user.Password == password {
			authorizedUser = user.Name
		}
	}

	if authorizedUser != "" {
		c.SetCookie("token", authToken, 3600*24, "/", "", false, true)
		c.Header("HX-Redirect", "/")
		c.Status(http.StatusOK)
	} else {
		c.HTML(http.StatusUnauthorized, "login-message", gin.H{
			"message": "Wrong username or password!",
		})
	}
}

func (h *Handler) Logout(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "", false, true)
	c.Header("HX-Redirect", "/")
	c.Status(http.StatusOK)
}
