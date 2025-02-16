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

	var a model.SendCoinRequest
	if err := c.ShouldBindJSON(&a); err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	receiver := a.ToUser
	amount := a.Amount

	_, err = db.DB.Exec("SELECT merch.f_transfer_coins($1::TEXT, $2::TEXT, $3)", username, receiver, amount)
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Errors: "Внутренняя ошибка сервера."})
		return
	}

	c.JSON(http.StatusOK, model.SendCoinResponse{Desc: "Успешный ответ"})
}
