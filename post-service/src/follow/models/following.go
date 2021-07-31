package models

type Following interface {
	GetUsername() string
}

type following struct {
	Username string
}

func NewFollowing(username string) Following {
	return &following{Username: username}
}

func (s *following) GetUsername() string {
	return s.Username
}
