package follow

const (
	addFollow      = "INSERT INTO followers(username, following) VALUES ($1, $2);"
	removeFollow   = "DELETE FROM followers WHERE username = $1 AND following = $2;"
	getFollowers   = "SELECT username FROM followers WHERE following = $1;"
	getFollowing   = "SELECT following FROM followers WHERE username = $1;"
	getSuggestions = `SELECT DISTINCT users.username
					  FROM users
					  WHERE users.username NOT IN (SELECT followers.following FROM followers WHERE followers.username=$1 )
					  AND users.username!=$1
					  LIMIT 10;`
)
