package model

type Cart struct {
	Id     int  `json:"id"`
	Good   Good `json:"good"`
	Amount int  `json:"amount"`
}
