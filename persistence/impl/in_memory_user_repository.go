package impl

import (
	"errors"
	"go.uber.org/zap"
	"math/big"
	"slices"
	"sync"
	"valsea_coding_challenge/domain/entity"
	"valsea_coding_challenge/util"
)

type InMemoryUserRepository struct {
	users   *util.TypedSyncMap[string, *entity.UserEntity]
	usersMu *util.TypedSyncMap[string, *sync.RWMutex]

	log *zap.Logger
}

func (i *InMemoryUserRepository) FindById(id string) (*entity.UserEntity, error) {
	mu, ok := i.usersMu.Load(id)
	if !ok {
		return nil, errors.New("user not found")
	}

	mu.RLock()
	defer mu.RUnlock()

	user, ok := i.users.Load(id)
	if !ok {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (i *InMemoryUserRepository) FindAll() []*entity.UserEntity {
	users := make([]*entity.UserEntity, 0, i.users.Size())
	i.usersMu.Range(func(key string, value *sync.RWMutex) bool {
		value.RLock()
		defer value.RUnlock()

		u, ok := i.users.Load(key)
		if !ok {
			i.log.Fatal("user not found in iterator")
		}

		users = append(users, u)
		return true
	})

	return users
}

func (i *InMemoryUserRepository) Save(user *entity.UserEntity) error {
	if _, ok := i.users.Load(user.Id.String()); ok {
		return errors.New("user already exists")
	}

	i.users.Store(user.Id.String(), user)
	i.usersMu.Store(user.Id.String(), &sync.RWMutex{})
	return nil
}

func (i *InMemoryUserRepository) Delete(user *entity.UserEntity) error {
	if _, ok := i.users.Load(user.Id.String()); !ok {
		return errors.New("user not found")
	}

	i.users.Delete(user.Id.String())
	i.usersMu.Delete(user.Id.String())
	return nil
}

func (i *InMemoryUserRepository) Update(user *entity.UserEntity) error {
	i.users.Store(user.Id.String(), user)
	return nil
}

func sortUserIds(users []*entity.UserEntity) {
	slices.SortFunc(users, func(a, b *entity.UserEntity) int {
		aId := big.NewInt(0).SetBytes(a.Id[:])
		bId := big.NewInt(0).SetBytes(b.Id[:])
		return aId.Cmp(bId)
	})
}

func (i *InMemoryUserRepository) getUserLocks(users []*entity.UserEntity) ([]*sync.RWMutex, error) {
	// we have to sort the users to avoid deadlocks
	sortUserIds(users)

	locks := make([]*sync.RWMutex, 0, len(users))
	for _, user := range users {
		mu, ok := i.usersMu.Load(user.Id.String())
		if !ok {
			return nil, errors.New("user not found")
		}
		locks = append(locks, mu)
	}

	return locks, nil
}

func (i *InMemoryUserRepository) Lock(users []*entity.UserEntity) error {
	locks, err := i.getUserLocks(users)
	if err != nil {
		return err
	}

	for _, lock := range locks {
		lock.Lock()
	}

	return nil
}

func (i *InMemoryUserRepository) Unlock(users []*entity.UserEntity) error {
	locks, err := i.getUserLocks(users)
	if err != nil {
		return err
	}

	for _, lock := range locks {
		lock.Unlock()
	}

	return nil
}

func NewInMemoryUserRepository(log *zap.Logger) *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users:   util.NewTypedSyncMap[string, *entity.UserEntity](),
		usersMu: util.NewTypedSyncMap[string, *sync.RWMutex](),
		log:     log,
	}
}
