package follow

import (
	"database/sql"
	"sync"

	"github.com/Emoto13/photo-viewer-rest/follow-service/src/follow/models"
)

type FollowStore interface {
	SaveFollow(follow *models.Follow) error
	RemoveFollow(follow *models.Follow) error
	GetFollowers(username string) ([]*models.Follower, error)
	GetFollowing(username string) ([]*models.Following, error)
	GetSuggestions(username string) ([]*models.Suggestion, error)
}

type followStore struct {
	db *sql.DB
	mu sync.RWMutex
}

func NewFollowStore(db *sql.DB) FollowStore {
	return &followStore{
		db: db,
		mu: sync.RWMutex{},
	}
}

func (store *followStore) SaveFollow(follow *models.Follow) error {
	store.mu.Lock()
	defer store.mu.Unlock()

	_, err := store.db.Exec(addFollow, follow.Username, follow.Following)
	if err != nil {
		return err
	}

	return nil
}

func (store *followStore) RemoveFollow(follow *models.Follow) error {
	store.mu.Lock()
	defer store.mu.Unlock()

	_, err := store.db.Exec(removeFollow, follow.Username, follow.Following)
	if err != nil {
		return err
	}

	return nil
}

func (store *followStore) GetFollowers(username string) ([]*models.Follower, error) {
	store.mu.RLock()
	defer store.mu.RUnlock()

	rows, err := store.db.Query(getFollowers, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	followerSlice := []*models.Follower{}
	for rows.Next() {
		follower := &models.Follower{}
		err = rows.Scan(&follower.Username)
		if err != nil {
			return nil, err
		}

		followerSlice = append(followerSlice, follower)
	}

	if err != nil {
		return nil, err
	}

	return followerSlice, nil
}

func (store *followStore) GetFollowing(username string) ([]*models.Following, error) {
	store.mu.RLock()
	defer store.mu.RUnlock()

	rows, err := store.db.Query(getFollowing, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	followingSlice := []*models.Following{}
	for rows.Next() {
		following := &models.Following{}
		err = rows.Scan(&following.Username)
		if err != nil {
			return nil, err
		}
		followingSlice = append(followingSlice, following)
	}

	if err != nil {
		return nil, err
	}

	return followingSlice, nil
}

func (store *followStore) GetSuggestions(username string) ([]*models.Suggestion, error) {
	store.mu.RLock()
	defer store.mu.RUnlock()

	rows, err := store.db.Query(getSuggestions, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	suggestions := []*models.Suggestion{}
	for rows.Next() {
		suggestion := &models.Suggestion{}
		err = rows.Scan(&suggestion.Username)
		if err != nil {
			return nil, err
		}

		suggestions = append(suggestions, suggestion)
	}

	if err != nil {
		return nil, err
	}

	return suggestions, nil
}
