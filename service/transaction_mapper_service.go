package service

import (
	"valsea_coding_challenge/domain/dto"
	"valsea_coding_challenge/domain/entity"
)

type TransactionMapperService interface {
	ToTransactionDto(transactionEntity *entity.TransactionEntity) (*dto.TransactionDTO, error)
}
