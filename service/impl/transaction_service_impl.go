package impl

import (
	"errors"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"valsea_coding_challenge/domain/dto"
	"valsea_coding_challenge/domain/entity"
	"valsea_coding_challenge/domain/enum"
	"valsea_coding_challenge/persistence"
	"valsea_coding_challenge/service"
)

type TransactionServiceImpl struct {
	transactionRepository    persistence.TransactionRepository
	transactionMapperService service.TransactionMapperService
	userRepository           persistence.UserRepository
	log                      *zap.Logger
}

func (t *TransactionServiceImpl) Transfer(fromUserId string, toUserId string, amount decimal.Decimal) ([]dto.TransactionDTO, error) {
	t.log.Debug("Transferring amount", zap.String("fromUserId", fromUserId), zap.String("toUserId", toUserId), zap.String("amount", amount.String()))

	if amount.IsZero() {
		t.log.Error("Transfer amount must be non-zero", zap.String("amount", amount.String()))
		return nil, errors.New("transfer amount must be non-zero")
	}

	if fromUserId == toUserId {
		t.log.Error("Cannot transfer to self", zap.String("fromUserId", fromUserId), zap.String("toUserId", toUserId))
		return nil, errors.New("cannot transfer to self")
	}

	if amount.IsNegative() {
		t.log.Error("Transfer amount must be positive", zap.String("amount", amount.String()))
		return nil, errors.New("transfer amount must be positive")
	}

	fromUser, err := t.userRepository.FindById(fromUserId)
	if err != nil {
		t.log.Error("Failed to find from user", zap.Error(err))
		return nil, err
	}

	toUser, err := t.userRepository.FindById(toUserId)
	if err != nil {
		t.log.Error("Failed to find to user", zap.Error(err))
		return nil, err
	}

	err = t.userRepository.Lock([]*entity.UserEntity{fromUser, toUser})
	if err != nil {
		t.log.Error("Failed to lock users", zap.Error(err))
		return nil, err
	}

	defer func() {
		err = t.userRepository.Unlock([]*entity.UserEntity{fromUser, toUser})
		if err != nil {
			t.log.Error("Failed to unlock users", zap.Error(err))
		}
	}()

	if fromUser.Balance.LessThan(amount) {
		t.log.Error("Insufficient funds", zap.String("balance", fromUser.Balance.String()), zap.String("amount", amount.String()))
		return nil, errors.New("insufficient funds")
	}

	fromTransaction := entity.NewTransactionEntity(toUser.Id, amount.Neg(), enum.TransactionTypeTransferOut)
	toTransaction := entity.NewTransactionEntity(fromUser.Id, amount, enum.TransactionTypeTransferIn)

	err = t.transactionRepository.SaveForUser(fromUserId, fromTransaction)
	if err != nil {
		t.log.Error("Failed to save transaction", zap.Error(err))
		return nil, err
	}

	err = t.transactionRepository.SaveForUser(toUserId, toTransaction)
	if err != nil {
		t.log.Error("Failed to save transaction", zap.Error(err))
		t.log.Debug("Rolling back transaction", zap.String("transactionId", fromTransaction.Id.String()))

		err = t.transactionRepository.DeleteTransaction(fromUserId, fromTransaction.Id.String())

		return nil, err
	}

	fromUser.Balance = fromUser.Balance.Sub(amount)
	toUser.Balance = toUser.Balance.Add(amount)

	err = t.userRepository.Update(fromUser)
	if err != nil {
		t.log.Error("Failed to update user", zap.Error(err))
		t.log.Debug("Rolling back transaction", zap.String("transactionId", fromTransaction.Id.String()))

		err = t.transactionRepository.DeleteTransaction(fromUserId, fromTransaction.Id.String())
		if err != nil {
			t.log.Error("Failed to delete transaction", zap.Error(err))
		}

		err = t.transactionRepository.DeleteTransaction(toUserId, toTransaction.Id.String())
		if err != nil {
			t.log.Error("Failed to delete transaction", zap.Error(err))
		}

		return nil, err
	}

	err = t.userRepository.Update(toUser)
	if err != nil {
		t.log.Error("Failed to update user", zap.Error(err))
		t.log.Debug("Rolling back transaction", zap.String("transactionId", fromTransaction.Id.String()))

		fromUser.Balance = fromUser.Balance.Add(amount)
		err = t.userRepository.Update(fromUser)
		if err != nil {
			t.log.Error("Failed to update user", zap.Error(err))
		}

		err = t.transactionRepository.DeleteTransaction(fromUserId, fromTransaction.Id.String())
		if err != nil {
			t.log.Error("Failed to delete transaction", zap.Error(err))
		}

		err = t.transactionRepository.DeleteTransaction(toUserId, toTransaction.Id.String())
		if err != nil {
			t.log.Error("Failed to delete transaction", zap.Error(err))
		}

		return nil, err
	}

	return []dto.TransactionDTO{*t.transactionMapperService.ToTransactionDto(fromTransaction), *t.transactionMapperService.ToTransactionDto(toTransaction)}, nil
}

func (t *TransactionServiceImpl) lockSingleUser(userId string, callback func(userEntity *entity.UserEntity) error) error {
	user, err := t.userRepository.FindById(userId)
	if err != nil {
		return err
	}

	err = t.userRepository.Lock([]*entity.UserEntity{user})
	if err != nil {
		return err
	}

	defer func() {
		err = t.userRepository.Unlock([]*entity.UserEntity{user})
		if err != nil {
			t.log.Error("Failed to unlock user", zap.Error(err))
		}
	}()

	return callback(user)
}

func (t *TransactionServiceImpl) Deposit(userId string, amount decimal.Decimal) (*dto.TransactionDTO, error) {
	t.log.Debug("Depositing amount", zap.String("userId", userId), zap.String("amount", amount.String()))

	if amount.IsNegative() {
		t.log.Error("Deposit amount must be positive", zap.String("amount", amount.String()))
		return nil, errors.New("deposit amount must be positive")
	}

	if amount.IsZero() {
		t.log.Error("Deposit amount must be non-zero", zap.String("amount", amount.String()))
		return nil, errors.New("deposit amount must be non-zero")
	}

	var transaction *entity.TransactionEntity
	err := t.lockSingleUser(userId, func(user *entity.UserEntity) error {
		transaction = entity.NewTransactionEntity(user.Id, amount, enum.TransactionTypeDeposit)
		err := t.transactionRepository.SaveForUser(userId, transaction)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		t.log.Error("Failed to deposit amount", zap.Error(err))
		return nil, err
	}

	return t.transactionMapperService.ToTransactionDto(transaction), nil
}

func (t *TransactionServiceImpl) Withdraw(userId string, amount decimal.Decimal) (*dto.TransactionDTO, error) {
	t.log.Debug("Withdrawing amount", zap.String("userId", userId), zap.String("amount", amount.String()))

	if amount.IsPositive() {
		t.log.Error("Withdrawal amount must be negative", zap.String("amount", amount.String()))
		return nil, errors.New("withdrawal amount must be negative")
	}

	if amount.IsZero() {
		t.log.Error("Withdrawal amount must be non-zero", zap.String("amount", amount.String()))
		return nil, errors.New("withdrawal amount must be non-zero")
	}

	var transaction *entity.TransactionEntity
	err := t.lockSingleUser(userId, func(user *entity.UserEntity) error {
		if user.Balance.LessThan(amount.Neg()) {
			t.log.Error("Insufficient funds", zap.String("balance", user.Balance.String()), zap.String("amount", amount.String()))
			return errors.New("insufficient funds")
		}

		transaction = entity.NewTransactionEntity(user.Id, amount, enum.TransactionTypeWithdrawal)
		err := t.transactionRepository.SaveForUser(userId, transaction)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		t.log.Error("Failed to withdraw amount", zap.Error(err))
		return nil, err
	}

	return t.transactionMapperService.ToTransactionDto(transaction), nil
}

func (t *TransactionServiceImpl) GetTransactions(userId string) ([]dto.TransactionDTO, error) {
	transactions, err := t.transactionRepository.GetAllTransactionsForUser(userId)
	if err != nil {
		return nil, err
	}

	var transactionDTOs []dto.TransactionDTO
	for _, transaction := range transactions {
		transactionDTOs = append(transactionDTOs, *t.transactionMapperService.ToTransactionDto(transaction))
	}

	return transactionDTOs, nil
}

func NewTransactionServiceImpl(log *zap.Logger, transactionRepository persistence.TransactionRepository, transactionMapperService service.TransactionMapperService, userRepository persistence.UserRepository) *TransactionServiceImpl {
	return &TransactionServiceImpl{
		transactionRepository:    transactionRepository,
		transactionMapperService: transactionMapperService,
		log:                      log,
		userRepository:           userRepository,
	}
}
