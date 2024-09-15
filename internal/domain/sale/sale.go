package sale

import (
	"andressa-lanches/internal/domain/addition"
	"time"

	"github.com/google/uuid"
)

type Sale struct {
	ID                uuid.UUID  `json:"id"`
	Date              time.Time  `json:"date"`
	TotalAmount       float64    `json:"total_amount"`
	Discount          float64    `json:"discount,omitempty"`
	AdditionalCharges float64    `json:"additional_charges,omitempty"`
	Items             []SaleItem `json:"items"`
}

type SaleItem struct {
	SaleID     uuid.UUID           `json:"sale_id"`
	ItemID     int                 `json:"item_id"`
	ProductID  uuid.UUID           `json:"product_id"`
	Quantity   int                 `json:"quantity"`
	UnitPrice  float64             `json:"unit_price"`
	TotalPrice float64             `json:"total_price"`
	Additions  []addition.Addition `json:"additions,omitempty"`
}
