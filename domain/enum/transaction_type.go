package enum

import (
	"encoding/json"
	"errors"
)

type TransactionType string

const (
	TransactionTypeDeposit     TransactionType = "deposit"
	TransactionTypeWithdrawal                  = "withdrawal"
	TransactionTypeTransferIn                  = "transfer_in"
	TransactionTypeTransferOut                 = "transfer_out"
)

var (
	validTransactionTypes = []TransactionType{
		TransactionTypeDeposit,
		TransactionTypeWithdrawal,
		TransactionTypeTransferIn,
		TransactionTypeTransferOut,
	}
)

func (t *TransactionType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	for _, validTransactionType := range validTransactionTypes {
		if s == string(validTransactionType) {
			*t = validTransactionType
			return nil
		}
	}

	return errors.New("invalid transaction type")
}
