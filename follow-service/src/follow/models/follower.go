package models

type Follower interface {
	GetUsername() string
}

type follower struct {
	Username string
}

func NewFollower(username string) Follower {
	return &follower{Username: username}
}

func (s *follower) GetUsername() string {
	return s.Username
}
