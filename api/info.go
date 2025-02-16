package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/panedbay/test-task/db"
	"github.com/panedbay/test-task/model"
)

func GetAPIInfo(c *gin.Context) {

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

	var coins int
	err = db.DB.QueryRow("SELECT merch.f_get_user_coins($1::TEXT)", username).Scan(&coins)
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Errors: "Внутренняя ошибка сервера."})
		return
	}

	results, err := db.DB.Query("SELECT * FROM merch.f_get_employee_inventory($1::TEXT)", username)
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Errors: "Внутренняя ошибка сервера."})
		return
	}
	defer results.Close()

	var inventory []model.Invent
	for results.Next() {
		var inv model.Invent
		if e := results.Scan(&inv.Type, &inv.Quantity); err != nil {
			log.Print(e)
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Errors: "Внутренняя ошибка сервера."})
			return
		}
		inventory = append(inventory, inv)
	}

	results, err = db.DB.Query("SELECT * FROM merch.f_get_transfers_sender($1::TEXT)", username)
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Errors: "Внутренняя ошибка сервера."})
		return
	}
	defer results.Close()

	var SH []model.SenHist
	for results.Next() {
		var sh model.SenHist
		if e := results.Scan(&sh.ToUser, &sh.Amount); err != nil {
			log.Print(e)
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Errors: "Внутренняя ошибка сервера."})
			return
		}
		SH = append(SH, sh)
	}

	results, err = db.DB.Query("SELECT * FROM merch.f_get_transfers_receiver($1::TEXT)", username)
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Errors: "Внутренняя ошибка сервера."})
		return
	}
	defer results.Close()

	var RH []model.RecHist
	for results.Next() {
		var rh model.RecHist
		if err := results.Scan(&rh.FromUser, &rh.Amount); err != nil {
			log.Print(err)
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Errors: "Внутренняя ошибка сервера."})
			return
		}
		RH = append(RH, rh)
	}

	var CH = model.CoinHist{Received: RH, Sent: SH}
	var info = model.InfoResponse{Coins: coins, Inventory: inventory, CoinHistory: CH}

	c.JSON(http.StatusOK, info)

}
