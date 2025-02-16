package model

type SendCoinRequest struct {
	Amount int `json:"amount"`

	ToUser string `json:"toUser"`
}

type SendCoinResponse struct {
	Desc string `json:"description"`
}
