package follow

import (
	"fmt"
	"sync"

	"github.com/Emoto13/photo-viewer-rest/follow-service/src/follow/models"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type FollowStore interface {
	CreateUser(username string) error
	SaveFollow(follow models.Follow) error
	RemoveFollow(follow models.Follow) error
	GetFollowers(username string) ([]*models.Follower, error)
	GetFollowing(username string) ([]*models.Following, error)
	GetSuggestions(username string) ([]*models.Suggestion, error)
}

type followStore struct {
	driver    neo4j.Driver
	connector Neo4jConnector
	mu        sync.RWMutex
}

func NewFollowStore(driver neo4j.Driver, connector Neo4jConnector) FollowStore {
	return &followStore{
		driver:    driver,
		connector: connector,
		mu:        sync.RWMutex{},
	}
}

func (store *followStore) CreateUser(username string) error {
	store.mu.Lock()
	defer store.mu.Unlock()

	session := store.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := store.connector.CreateUser(username)
	_, err := session.WriteTransaction(query)
	if err != nil {
		return err
	}

	return nil
}

func (store *followStore) SaveFollow(follow models.Follow) error {
	store.mu.Lock()
	defer store.mu.Unlock()

	session := store.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := store.connector.SaveFollow(follow.GetUsername(), follow.GetFollowing())
	_, err := session.WriteTransaction(query)
	if err != nil {
		return err
	}

	return nil
}

func (store *followStore) RemoveFollow(follow models.Follow) error {
	store.mu.Lock()
	defer store.mu.Unlock()

	session := store.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := store.connector.RemoveFollow(follow.GetUsername(), follow.GetFollowing())
	_, err := session.WriteTransaction(query)
	if err != nil {
		return err
	}

	return nil
}

func (store *followStore) GetFollowers(username string) ([]*models.Follower, error) {
	store.mu.RLock()
	defer store.mu.RUnlock()

	session := store.driver.NewSession(neo4j.SessionConfig{})
	fmt.Println("Session opening finished")

	defer session.Close()

	query := store.connector.GetFollowers(username)
	followers, err := session.ReadTransaction(query)
	fmt.Println("Transaction not finished")

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	fmt.Println("Transaction finished")
	fmt.Println(followers)
	return followers.([]*models.Follower), nil
}

func (store *followStore) GetFollowing(username string) ([]*models.Following, error) {
	store.mu.RLock()
	defer store.mu.RUnlock()

	session := store.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := store.connector.GetFollowings(username)
	followings, err := session.ReadTransaction(query)
	if err != nil {
		fmt.Println("error getting followings:", err)
		return nil, err
	}

	return followings.([]*models.Following), nil
}

func (store *followStore) GetSuggestions(username string) ([]*models.Suggestion, error) {
	store.mu.RLock()
	defer store.mu.RUnlock()

	session := store.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := store.connector.GetSuggestions(username)
	suggestions, err := session.ReadTransaction(query)
	if err != nil {
		return nil, err
	}

	return suggestions.([]*models.Suggestion), nil
}
