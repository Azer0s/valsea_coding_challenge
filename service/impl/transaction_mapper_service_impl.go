package impl

import (
	"valsea_coding_challenge/domain/dto"
	"valsea_coding_challenge/domain/entity"
)

type TransactionMapperServiceImpl struct {
}

func (t *TransactionMapperServiceImpl) ToTransactionDto(transactionEntity *entity.TransactionEntity) *dto.TransactionDTO {
	return &dto.TransactionDTO{
		Id:              transactionEntity.Id.String(),
		RelatedUserId:   transactionEntity.RelatedUserId.String(),
		Amount:          transactionEntity.Amount.InexactFloat64(),
		Timestamp:       transactionEntity.Timestamp.String(),
		TransactionType: transactionEntity.TransactionType,
	}
}
