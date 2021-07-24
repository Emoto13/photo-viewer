package user

const (
	retrieveUser = "SELECT username, hashed_password FROM users WHERE username=$1;"
)
