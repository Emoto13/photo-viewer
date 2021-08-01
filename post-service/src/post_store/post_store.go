package post_store

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/Emoto13/photo-viewer-rest/post-service/src/post_store/post_data"
	"github.com/gocql/gocql"
)

type PostStore interface {
	CreatePost(username string, post *post_data.PostData) error
	SearchPosts(name string) ([]*post_data.PostData, error)
	GetUserPosts(username string) ([]*post_data.PostData, error)
}

type postStore struct {
	db      *sql.DB
	session *gocql.Session
	mu      sync.RWMutex
}

func NewPostStore(db *sql.DB, session *gocql.Session) PostStore {
	return &postStore{
		db:      db,
		session: session,
		mu:      sync.RWMutex{},
	}
}

func (store *postStore) CreatePost(username string, post *post_data.PostData) error {
	posts := []*post_data.PostData{}
	err := store.session.Query(`SELECT posts FROM posts WHERE username=?`, username).Scan(&posts)
	if err != nil && err != gocql.ErrNotFound {
		fmt.Println("couldn't get posts from cassandra:", err)
		return err
	}

	posts = append([]*post_data.PostData{post}, posts...)
	err = store.session.Query(`INSERT INTO posts(username, posts) VALUES (?, ?);`, username, posts).Exec()
	if err != nil {
		fmt.Println("couldn't write posts from cassandra:", err)
		return err
	}

	return nil
}

func (store *postStore) GetUserPosts(username string) ([]*post_data.PostData, error) {
	store.mu.RLock()
	defer store.mu.RUnlock()

	posts := []*post_data.PostData{}
	err := store.session.Query(`SELECT posts FROM posts WHERE username=?`, username).Scan(&posts)
	if err != nil && err != gocql.ErrNotFound {
		fmt.Println("couldnt retrieve user posts", err.Error())
		return nil, err
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
