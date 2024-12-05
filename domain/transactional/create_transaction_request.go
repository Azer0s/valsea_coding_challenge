package transactional

import (
	"encoding/json"
	"errors"
	"github.com/shopspring/decimal"
	"valsea_coding_challenge/domain/enum"
)

type createTransactionRequestSchema struct {
	Amount decimal.Decimal      `json:"amount"`
	Type   enum.TransactionType `json:"type"`
}

type CreateTransactionRequest struct {
	Amount decimal.Decimal

	// Type in this case can only be `deposit` or `withdrawal`
	Type enum.TransactionType
}

func (r *CreateTransactionRequest) UnmarshalJSON(data []byte) error {
	var schema createTransactionRequestSchema
	err := json.Unmarshal(data, &schema)
	if err != nil {
		return err
	}

	if schema.Type != enum.TransactionTypeDeposit && schema.Type != enum.TransactionTypeWithdrawal {
		return errors.New("transaction type must be deposit or withdrawal")
	}

	r.Amount = schema.Amount
	r.Type = schema.Type

	return nil
}
