package store

const (
	createUser = "INSERT INTO users(username, hashed_password, role, email) VALUES ($1, $2, $3, $4)"
)
