package order

type OrderUpdateRequest struct {
	StatusOrder string `validate:"required" json:"status_order"`
	Description string `json:"description,omitempty"`
}
