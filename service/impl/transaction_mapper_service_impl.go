package impl

import (
	"valsea_coding_challenge/domain/dto"
	"valsea_coding_challenge/domain/entity"
)

type TransactionMapperServiceImpl struct {
}

func (t *TransactionMapperServiceImpl) ToTransactionDto(transactionEntity *entity.TransactionEntity) (*dto.TransactionDTO, error) {
	return &dto.TransactionDTO{
		Id:              transactionEntity.Id.String(),
		RelatedAccount:  transactionEntity.RelatedAccount.String(),
		Amount:          transactionEntity.Amount.InexactFloat64(),
		Timestamp:       transactionEntity.Timestamp.String(),
		TransactionType: transactionEntity.TransactionType,
	}, nil
}
