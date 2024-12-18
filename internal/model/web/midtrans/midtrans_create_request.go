package midtrans

type MidtransRequest struct {
	OrderId         string `json:"orderid" binding:"required"`
	CustomerName    string `json:"customer_name" binding:"required"`
	CustomerID      string `json:"customer_id" binding:"required"`
	CustomerEmail   string `json:"customer_email" binding:"required"`
	CostumeId       int    `json:"costume_id" binding:"required"`
	CostumeName     string `json:"costume_name" binding:"required"`
	CostumeCategory string `json:"costume_category" binding:"required"`
	Price           int64  `json:"costume_price" binding:"required"`
	FinalAmount     int64  `json:"final_amount" binding:"required"`
	MerchantName    string `json:"merchant_name" binding:"required"`
}
