package api

type Product struct {
	Id      int    `json:"Id"`
	Manager string `json:"Manager"`
	Title   string `json:"Title"`
	Price   uint   `json:"Price"`
}

type Basket struct {
	Positions []BasketPosition `json:"Items"`
}

type BasketPosition struct {
	Id       int     `json:"Id"`
	Payer    User    `json:"Payer"`
	Owner    User    `json:"Owner"`
	Product  Product `json:"Product"`
	Paid     bool    `json:"Paid"`
	Price    uint    `json:"PriceDiscount"`
	Discount uint    `json:"Discount"`
	Coupon   string  `json:"CouponCode"`
}
