package post_store

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/Emoto13/photo-viewer-rest/post-service/src/post_store/post_data"
)

type PostStore interface {
	RetrieveFollowingPosts(username string) ([]*post_data.PostData, error)
	SearchPosts(name string) ([]*post_data.PostData, error)
}

type postStore struct {
	db *sql.DB
	mu sync.RWMutex
}

func NewPostStore(db *sql.DB) PostStore {
	return &postStore{
		db: db,
		mu: sync.RWMutex{},
	}
}

func (store *postStore) RetrieveFollowingPosts(username string) ([]*post_data.PostData, error) {
	store.mu.RLock()
	defer store.mu.RUnlock()

	rows, err := store.db.Query(getPostsOfFollowing, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := []*post_data.PostData{}
	for rows.Next() {
		post := &post_data.PostData{}

		err = rows.Scan(&post.Name, &post.Path, &post.Owner, &post.CreatedOn)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (store *postStore) SearchPosts(name string) ([]*post_data.PostData, error) {
	store.mu.RLock()
	defer store.mu.RUnlock()

	rows, err := store.db.Query(searchPosts, name+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := []*post_data.PostData{}
	for rows.Next() {
		post := &post_data.PostData{}
		err = rows.Scan(&post.Name, &post.Path, &post.Owner, &post.CreatedOn)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	fmt.Println(posts)
	return posts, nil
}
