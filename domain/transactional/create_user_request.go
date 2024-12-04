package transactional

import (
	"github.com/shopspring/decimal"
)

type CreateUserRequest struct {
	Owner          string          `json:"owner"`
	InitialBalance decimal.Decimal `json:"initial_balance"`
}
