package dto

type CardDto struct {
	Coin     int `form:"coin"`
	Quantity int `form:"quantity"`
}
