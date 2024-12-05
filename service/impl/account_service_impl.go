package impl

import (
	"errors"
	"go.uber.org/zap"
	"valsea_coding_challenge/domain/dto"
	"valsea_coding_challenge/domain/entity"
	"valsea_coding_challenge/domain/transactional"
	"valsea_coding_challenge/persistence"
	"valsea_coding_challenge/service"
)

type AccountServiceImpl struct {
	userRepository    persistence.UserRepository
	userMapperService service.UserMapperService

	transactionRepository persistence.TransactionRepository

	log *zap.Logger
}

func (a *AccountServiceImpl) CreateAccount(request transactional.CreateUserRequest) (*dto.UserDTO, error) {
	if request.Owner == "" {
		return nil, errors.New("name of owner is required")
	}

	if request.InitialBalance.IsNegative() {
		return nil, errors.New("initial balance must be positive")
	}

	userEntity := entity.NewUserEntity(request.Owner, request.InitialBalance)

	a.log.Debug("Creating account", zap.String("owner", request.Owner), zap.String("initialBalance", request.InitialBalance.String()))
	err := a.userRepository.Save(userEntity)
	if err != nil {
		return nil, err
	}

	a.log.Debug("Creating transaction history for user", zap.String("userId", userEntity.Id.String()))
	err = a.transactionRepository.CreateHistory(userEntity.Id.String())
	if err != nil {
		//rollback user creation
		err = a.userRepository.Delete(userEntity)
		if err != nil {
			a.log.Error("Failed to rollback user creation", zap.Error(err))
		}

		return nil, err
	}

	a.log.Debug("Created account", zap.String("id", userEntity.Id.String()))
	return a.userMapperService.ToUserDto(userEntity)
}

func (a *AccountServiceImpl) GetAccountById(id string) (*dto.UserDTO, error) {
	user, err := a.userRepository.FindById(id)
	if err != nil {
		return nil, err
	}

	return a.userMapperService.ToUserDto(user)
}

func (a *AccountServiceImpl) GetAllAccounts() ([]dto.UserDTO, error) {
	users := a.userRepository.FindAll()
	userDtos := make([]dto.UserDTO, 0, len(users))
	for _, user := range users {
		userDto, err := a.userMapperService.ToUserDto(user)
		if err != nil {
			return nil, err
		}
		userDtos = append(userDtos, *userDto)
	}
	return userDtos, nil
}

func NewAccountServiceImpl(userRepository persistence.UserRepository, transactionRepository persistence.TransactionRepository, userMapperService service.UserMapperService, log *zap.Logger) *AccountServiceImpl {
	return &AccountServiceImpl{
		userRepository:        userRepository,
		userMapperService:     userMapperService,
		transactionRepository: transactionRepository,
		log:                   log,
	}
}
