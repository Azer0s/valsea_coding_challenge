package transactional

import "github.com/shopspring/decimal"

type CreateTransferRequest struct {
	FromUserId string          `json:"from_account_id"`
	ToUserId   string          `json:"to_account_id"`
	Amount     decimal.Decimal `json:"amount"`
}
