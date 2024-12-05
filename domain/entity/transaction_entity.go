package entity

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"time"
	"valsea_coding_challenge/domain/enum"
)

type TransactionEntity struct {
	Id uuid.UUID

	// RelatedUserId is the account that is affected by this transaction
	// For example, if the transaction is `transfer_in` then the related account is the account which sends the money
	// If the transaction is `transfer_out` then the related account is the account which receives the money
	// If the transaction is `deposit` or `withdraw` then the related account is the account itself
	RelatedUserId uuid.UUID
	Amount        decimal.Decimal

	// Timestamp unix timestamp
	Timestamp time.Time

	TransactionType enum.TransactionType
}

func NewTransactionEntity(relatedUserId uuid.UUID, amount decimal.Decimal, transactionType enum.TransactionType) *TransactionEntity {
	return &TransactionEntity{
		Id:              uuid.New(),
		RelatedUserId:   relatedUserId,
		Amount:          amount,
		Timestamp:       time.Now(),
		TransactionType: transactionType,
	}
}
