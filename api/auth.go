package api

import (
	"log"
	"net/http"
	"time"

	"errors"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/panedbay/test-task/db"
	"github.com/panedbay/test-task/model"
)

func GetUserFromJWT(token string) (string, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})

	if err != nil {
		return "", err
	}

	username, ok := claims["username"]
	if !ok {
		return "", errors.New("failed to find username")
	}

	return username.(string), nil
}

func PostAPIAuth(c *gin.Context) {
	var a model.AuthRequest
	if err := c.ShouldBindJSON(&a); err != nil {
		log.Print(err)
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Errors: "Неверный запрос"})
		return
	}
	username := a.Username
	password := a.Password
	if username == "" || password == "" {
		log.Printf("Invalid username or password fields - username:%s password:%s", username, password)
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Errors: "Неверный запрос"})
		return
	}

	var result bool
	err := db.DB.QueryRow("SELECT merch.f_employee_exists($1::TEXT)", username).Scan(&result)
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Errors: "Внутренняя ошибка сервера."})
		return
	}

	if result {
		var res bool
		e := db.DB.QueryRow("SELECT merch.f_check_employee_credentials($1::TEXT, $2::TEXT)", username, password).Scan(&res)
		if e != nil || !res {
			log.Print(e)
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{Errors: "Неавторизован."})
			return
		}
		token, e := IssueJWT(username)
		if e != nil {
			log.Print(e)
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Errors: "Неавторизован."})
			return
		}
		c.JSON(http.StatusOK, model.AuthResponse{Token: token})
		return
	}
	err = addUser(a.Username, a.Password)
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Errors: "Внутренняя ошибка сервера."})
		return
	}

	token, err := IssueJWT(username)
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Errors: "Внутренняя ошибка сервера."})
		return
	}
	c.JSON(http.StatusOK, model.AuthResponse{Token: token})
}

func IssueJWT(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
		"iat":      time.Now().Unix(),
	})

	return token.SignedString([]byte("secret"))
}

func addUser(username, password string) error {
	_, e := db.DB.Exec("SELECT merch.f_add_employee($1::TEXT, $2::TEXT)", username, password)
	if e != nil {
		return e
	}
	return nil
}
