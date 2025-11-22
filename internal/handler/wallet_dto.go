package handler

type UpdateWalletRequest struct {
	WalletID      string `json:"walletId" binding:"required"`
	OperationType string `json:"operationType" binding:"required,oneof=DEPOSIT WITHDRAW"`
	Amount        int64  `json:"amount" binding:"required"`
}

type UpdateWalletResponse struct {
	WalletID   string `json:"walletId"`
	NewBalance int64  `json:"newBalance"`
}

type GetWalletResponse struct {
	WalletID string `json:"walletId"`
	Balance  int64  `json:"balance"`
}
