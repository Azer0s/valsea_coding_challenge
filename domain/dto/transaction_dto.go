package dto

import (
	"valsea_coding_challenge/domain/enum"
)

type TransactionDTO struct {
	Id              string               `json:"id"`
	RelatedUserId   string               `json:"account_id"`
	Amount          float64              `json:"amount"`
	Timestamp       string               `json:"timestamp"`
	TransactionType enum.TransactionType `json:"type"`
}
