package model

type BuyItemRequest struct {
	ItemName string `json:"item_name"`
}

type BuyItemResponse struct {
	Desc string `json:"description"`
}
