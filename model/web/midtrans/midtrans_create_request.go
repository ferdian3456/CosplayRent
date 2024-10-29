package midtrans

type MidtransRequest struct {
	OrderId     string
	UserId      string `json:"user_id" binding:"required"`
	Amount      int64  `json:"amount" binding:"required"`
	CostumeId   int    `json:"costume_id" binding:"required"`
	CostumeName string `json:"costume_name" binding:"required"`
}
