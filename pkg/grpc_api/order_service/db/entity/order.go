package entity

import "time"

type CreateOrderReq struct {
	AddressID     string  `json:"address_id" validate:"required"`
	CartID        string  `json:"cart_id" validate:"required"`
	PromoID       *string `json:"promo_id"`
	Amount        float64 `json:"amount" validate:"required"`
	PaymentMode   string  `json:"action" validate:"payment_mode"`
	UserID        string  `json:"user_id"`
	TransactionID *string `json:"transaction_id"`
}

type UserOrderRes struct {
	ID           string  `json:"id"`
	ProductID    string  `json:"product_id"`
	ProductName  string  `json:"product_name"`
	ProductImage string  `json:"product_image"`
	Size         float32 `json:"size"`
	Price        float64 `json:"price"`
	Status       string  `json:"status"`
}

type GetTrackingDetailsReq struct {
	OrderItemID string `json:"order_item_id" validate:"required"`
}

type TrackingDetailsRes struct {
	Status  int32     `json:"status"`
	Title   string    `json:"title"`
	Summary string    `json:"summary"`
	Time    time.Time `json:"time"`
}
