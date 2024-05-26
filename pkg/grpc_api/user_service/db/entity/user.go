package entity

import "github.com/akmal4410/gestapo/pkg/grpc_api/product_service/db/entity"

type GetHomeRes struct {
	Discount  *DiscountRes           `json:"discount,omitempty"`
	Merchants []MerchantRes          `json:"merchants,omitempty"`
	Products  []entity.GetProductRes `json:"products,omitempty"`
}

type DiscountRes struct {
	ProductID    string  `json:"product_id"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	Percentage   float64 `json:"percentage"`
	ProductImage string  `json:"product_image"`
	CardColor    uint32  `json:"card_color"`
}

type MerchantRes struct {
	MerchantID string  `json:"merchant_id"`
	Name       string  `json:"name"`
	ImageURL   *string `json:"image_url,omitempty"`
}

type AddRemoveWishlistReq struct {
	Action    string `json:"action" validate:"wishlist_action"`
	ProductID string `json:"product_id" validate:"required"`
	UserID    string `json:"user_id" validate:"required"`
}

type AddToCartReq struct {
	ProductID string  `json:"product_id" validate:"required"`
	Size      float64 `json:"size" validate:"required"`
	Quantity  int32   `json:"quantity" validate:"required"`
	Price     float64 `json:"price" validate:"required"`
	UserID    string  `json:"user_id" validate:"required"`
}

type CartRes struct {
	CartID string  `json:"cart_id"`
	UserID string  `json:"user_id"`
	Price  float64 `json:"price"`
}

type CartItemRes struct {
	ProductID         string  `json:"product_id"`
	CartID            string  `json:"cart_id"`
	CartItemID        string  `json:"cart_item_id"`
	ImageURL          string  `json:"image_url"`
	Name              string  `json:"name"`
	Size              float64 `json:"size"`
	Price             float64 `json:"price"`
	Quantity          int32   `json:"quantity"`
	AvailableQuantity int32   `json:"available_quantity"`
	InventoryID       string  `json:"inventory_id"`
}

type CheckoutCartItemsReq struct {
	CartID string         `json:"cart_id" validate:"required"`
	Data   []*CheckoutReq `json:"data" validate:"required"`
}

type CheckoutReq struct {
	CartItemID string `json:"cart_item_id" validate:"required"`
	Quantity   int32  `json:"quantity" validate:"required"`
}

type AddAddressReq struct {
	Title       string  `json:"title" validate:"required"`
	AddressLine string  `json:"address_line" validate:"required"`
	Country     string  `json:"country" validate:"required"`
	City        string  `json:"city" validate:"required"`
	PostalCode  *int64  `json:"postal_code"`
	Landmark    *string `json:"landmark"`
	IsDefault   *bool   `json:"is_default"`
	UserID      string  `json:"user_id"`
}

type GetAddressRes struct {
	AddressID   string  `json:"address_id"`
	UserID      string  `json:"user_id"`
	Title       string  `json:"title"`
	AddressLine string  `json:"address_line"`
	Country     *string `json:"country"`
	City        *string `json:"city"`
	PostalCode  *int64  `json:"postal_code"`
	Landmark    *string `json:"landmark"`
	IsDefault   *bool   `json:"is_default"`
}

type EditAddressReq struct {
	Title       *string `json:"title"`
	AddressLine *string `json:"address_line"`
	Country     *string `json:"country"`
	City        *string `json:"city"`
	PostalCode  *int64  `json:"postal_code"`
	Landmark    *string `json:"landmark"`
	IsDefault   *bool   `json:"is_default"`
}

type GetOrdersReq struct {
	Type string `json:"type" validate:"order_type"`
}
