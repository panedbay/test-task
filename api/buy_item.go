package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/panedbay/test-task/db"
	"github.com/panedbay/test-task/model"
)

func GetAPIBuyItem(c *gin.Context) {
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

	item := c.Param("item")

	_, err = db.DB.Exec("SELECT merch.f_buy($1::TEXT, $2::TEXT)", item, username)
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Errors: "Внутренняя ошибка сервера."})
		return
	}

	c.JSON(http.StatusOK, model.BuyItemResponse{Desc: "Успешный ответ"})
}
