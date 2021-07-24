package store

import (
	"database/sql"
	"sync"
)

type UserStore interface {
	SaveUser(user *User) error
}

type userStore struct {
	mutex sync.RWMutex
	db    *sql.DB
}

func NewUserStore(db *sql.DB) UserStore {
	return &userStore{
		db: db,
	}
}

func (store *userStore) SaveUser(user *User) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	_, err := store.db.Exec(createUser, user.Username, user.HashedPassword, user.Role, user.Email)
	if err != nil {
		return err
	}

	return nil
}
