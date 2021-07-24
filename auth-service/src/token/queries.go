package token

const (
	getUsernameFromToken = "SELECT username FROM tokens WHERE token = $1;"
	createToken          = `INSERT INTO tokens(token, username) VALUES ($1, $2) 
							ON CONFLICT (username) 
							DO UPDATE SET token = $3;`
	removeToken = "DELETE FROM tokens WHERE username = $1 AND token = $2;"
)
