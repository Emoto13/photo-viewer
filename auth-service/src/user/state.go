package user

import (
	"database/sql"
	"fmt"
	"sync"
)

type State interface {
	RetrieveUser(username string) (*User, error)
}

type state struct {
	db *sql.DB
	mu sync.RWMutex
}

func NewState(db *sql.DB) State {
	return &state{
		db: db,
		mu: sync.RWMutex{},
	}
}

func (s *state) RetrieveUser(username string) (*User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	row := s.db.QueryRow(retrieveUser, username)

	var user User
	err := row.Scan(&user.Username, &user.HashedPassword)
	if err != nil {
		return nil, err
	}

	fmt.Println(user)
	return &user, nil
}
