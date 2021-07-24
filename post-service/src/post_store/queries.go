package post_store

const (
	getPostsOfFollowing = `SELECT posts.username, images.path, posts.name, posts.created_on 
						   FROM followers
						   INNER JOIN posts ON followers.following=posts.username
						   INNER JOIN images ON images.id = posts.image_id
						   WHERE followers.username = $1
						   ORDER BY posts.created_on DESC;`
	searchPosts = `SELECT posts.name, images.path, posts.username,posts.created_on 
				   FROM posts
				   INNER JOIN images ON posts.image_id=images.id
				   WHERE LOWER(posts.name) LIKE LOWER($1);`
)
