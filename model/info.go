package model

type InfoResponse struct {
	CoinHistory CoinHist `json:"coinHistory,omitempty"`
	Coins       int      `json:"coins,omitempty"`
	Inventory   []Invent `json:"inventory,omitempty"`
}

type Invent struct {
	Quantity int    `json:"quantity,omitempty"`
	Type     string `json:"type,omitempty"`
}

type RecHist struct {
	Amount   int    `json:"amount,omitempty"`
	FromUser string `json:"fromUser,omitempty"`
}

type SenHist struct {
	Amount int    `json:"amount,omitempty"`
	ToUser string `json:"toUser,omitempty"`
}

type CoinHist struct {
	Received []RecHist `json:"received,omitempty"`
	Sent     []SenHist `json:"sent,omitempty"`
}
