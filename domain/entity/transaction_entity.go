package entity

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"time"
	"valsea_coding_challenge/domain/enum"
)

type TransactionEntity struct {
	Id uuid.UUID

	// RelatedAccount is the account that is affected by this transaction
	// For example, if the transaction is `transfer_in` then the related account is the account which sends the money
	// If the transaction is `transfer_out` then the related account is the account which receives the money
	// If the transaction is `deposit` or `withdraw` then the related account is the account itself
	RelatedAccount uuid.UUID
	Amount         decimal.Decimal

	// Timestamp unix timestamp
	Timestamp time.Time

	TransactionType enum.TransactionType
}
