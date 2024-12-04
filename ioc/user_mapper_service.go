package ioc

import (
	"valsea_coding_challenge/service"
	"valsea_coding_challenge/service/impl"
)

func ProvideUserMapperService() service.UserMapperService {
	return impl.NewUserMapperServiceImpl()
}
