package models

type Suggestion interface {
	GetUsername() string
}

type suggestion struct {
	Username string
}

func NewSuggestion(username string) Suggestion {
	return &suggestion{Username: username}
}

func (s *suggestion) GetUsername() string {
	return s.Username
}
