package models

type Follow struct {
	Username  string
	Following string
}

func NewFollow(username string, following string) *Follow {
	return &Follow{Username: username, Following: following}
}
