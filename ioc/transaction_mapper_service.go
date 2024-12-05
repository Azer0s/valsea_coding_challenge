package ioc

import (
	"valsea_coding_challenge/service"
	"valsea_coding_challenge/service/impl"
)

func ProvideTransactionMapperService() service.TransactionMapperService {
	return &impl.TransactionMapperServiceImpl{}
}
