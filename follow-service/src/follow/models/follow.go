package models

type Follow interface {
	GetUsername() string
	GetFollowing() string
}

type follow struct {
	Username  string
	Following string
}

func NewFollow(username string, following string) Follow {
	return &follow{Username: username, Following: following}
}

func (f *follow) GetUsername() string {
	return f.Username
}

func (f *follow) GetFollowing() string {
	return f.Following
}
