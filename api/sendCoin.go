package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/panedbay/test-task/db"
	"github.com/panedbay/test-task/model"
)

func PostAPISendCoin(c *gin.Context) {
	authHeader := c.Request.Header["Authorization"]
	if authHeader[0] == "" || len(authHeader[0]) < 8 {
		c.JSON(http.StatusUnauthorized, model.ErrorResponse{Errors: "Неавторизован."})
		return
	}
	jwtToken := authHeader[0][7:]

	username, err := GetUserFromJWT(jwtToken)

	if err != nil {
		log.Print(err)
		c.JSON(http.StatusUnauthorized, model.ErrorResponse{Errors: "Неавторизован."})
		return
	}

	var s model.SendCoinRequest
	if e := c.ShouldBindJSON(&s); e != nil {
		log.Print(e)
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Errors: "Неверный запрос"})
		return
	}

	receiver := s.ToUser
	amount := s.Amount

	if receiver == "" || amount <= 0 {
		log.Printf("Invalid ToUser or Amount fields - ToUser:%s Amount:%d", receiver, amount)
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Errors: "Неверный запрос"})
		return
	}

	_, err = db.DB.Exec("SELECT merch.f_transfer_coins($1::TEXT, $2::TEXT, $3)", username, receiver, amount)
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Errors: "Внутренняя ошибка сервера."})
		return
	}

	c.JSON(http.StatusOK, model.SendCoinResponse{Desc: "Успешный ответ"})
}
