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

type UserServiceImpl struct {
	userRepository    persistence.UserRepository
	userMapperService service.UserMapperService

	transactionRepository persistence.TransactionRepository

	log *zap.Logger
}

func (a *UserServiceImpl) CreateUser(request transactional.CreateUserRequest) (*dto.UserDTO, error) {
	a.log.Info("Creating account for user", zap.String("owner", request.Owner), zap.String("initialBalance", request.InitialBalance.String()))

	if request.Owner == "" {
		return nil, errors.New("name of owner is required")
	}

	if request.InitialBalance.IsNegative() {
		return nil, errors.New("initial balance must be positive")
	}

	userEntity := entity.NewUserEntity(request.Owner, request.InitialBalance)

	a.log.Debug("Creating user entity", zap.String("owner", userEntity.Name), zap.String("balance", userEntity.Balance.String()))
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
	return a.userMapperService.ToUserDto(userEntity), nil
}

func (a *UserServiceImpl) GetUserById(id string) (*dto.UserDTO, error) {
	user, err := a.userRepository.FindById(id)
	if err != nil {
		return nil, err
	}

	return a.userMapperService.ToUserDto(user), nil
}

func (a *UserServiceImpl) GetAllUsers() ([]dto.UserDTO, error) {
	users := a.userRepository.FindAll()
	userDtos := make([]dto.UserDTO, 0, len(users))
	for _, user := range users {
		userDtos = append(userDtos, *a.userMapperService.ToUserDto(user))
	}
	return userDtos, nil
}

func NewUserServiceImpl(userRepository persistence.UserRepository, transactionRepository persistence.TransactionRepository, userMapperService service.UserMapperService, log *zap.Logger) *UserServiceImpl {
	return &UserServiceImpl{
		userRepository:        userRepository,
		userMapperService:     userMapperService,
		transactionRepository: transactionRepository,
		log:                   log,
	}
}
