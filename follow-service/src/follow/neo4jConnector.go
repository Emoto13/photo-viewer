package follow

import (
	"fmt"

	"github.com/Emoto13/photo-viewer-rest/follow-service/src/follow/models"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type Neo4jConnector interface {
	CreateUser(username string) func(tx neo4j.Transaction) (interface{}, error)
	SaveFollow(followerUsername, followingUsername string) func(tx neo4j.Transaction) (interface{}, error)
	RemoveFollow(followerUsername, followingUsername string) func(tx neo4j.Transaction) (interface{}, error)

	GetFollowers(username string) func(tx neo4j.Transaction) (interface{}, error)
	GetFollowings(username string) func(tx neo4j.Transaction) (interface{}, error)
	GetSuggestions(username string) func(tx neo4j.Transaction) (interface{}, error)
}

type neo4jConnector struct{}

func NewNeo4jConnector() Neo4jConnector {
	return &neo4jConnector{}
}

func (n *neo4jConnector) CreateUser(username string) func(tx neo4j.Transaction) (interface{}, error) {
	return func(tx neo4j.Transaction) (interface{}, error) {
		_, err := tx.Run("MERGE (n:User { username: $username }) RETURN n.username", map[string]interface{}{
			"username": username,
		})
		if err != nil {
			return nil, err
		}

		return nil, nil
	}
}

func (n *neo4jConnector) SaveFollow(followerUsername, followingUsername string) func(tx neo4j.Transaction) (interface{}, error) {
	return func(tx neo4j.Transaction) (interface{}, error) {
		_, err := tx.Run(`MATCH (follower:User), (following:User)
						  WHERE follower.username = $followerUsername AND following.username = $followingUsername
						  CREATE (follower)-[r:FOLLOWS]->(following)
						  RETURN type(r);`,
			map[string]interface{}{
				"followerUsername":  followerUsername,
				"followingUsername": followingUsername,
			})
		if err != nil {
			return nil, err
		}

		return nil, nil
	}
}

func (n *neo4jConnector) RemoveFollow(followerUsername, followingUsername string) func(tx neo4j.Transaction) (interface{}, error) {
	return func(tx neo4j.Transaction) (interface{}, error) {
		_, err := tx.Run(`MATCH (follower:User)-[r:FOLLOWS]->(following:User) 
						  WHERE follower.username = $followerUsername AND following.username = $followingUsername
						  DELETE r;`,
			map[string]interface{}{
				"followerUsername":  followerUsername,
				"followingUsername": followingUsername,
			})
		if err != nil {
			return nil, err
		}

		return nil, nil
	}
}

func (n *neo4jConnector) GetFollowers(username string) func(tx neo4j.Transaction) (interface{}, error) {
	return func(tx neo4j.Transaction) (interface{}, error) {
		records, err := tx.Run("MATCH (n:User { username: $username })-->(user) RETURN user.username", map[string]interface{}{
			"username": username,
		})

		if err != nil {
			fmt.Println("couldn't retrieve followers", err)
			return nil, err
		}

		slice, err := records.Collect()
		if err != nil {
			fmt.Println("error getting followers: ", err)
			return nil, err
		}

		followers := []models.Follower{}
		for _, entry := range slice {
			fmt.Println("entry: ", entry.Values)
			follower := models.NewFollower(entry.Values[0].(string))
			followers = append(followers, follower)
		}

		return followers, nil
	}
}

func (n *neo4jConnector) GetFollowings(username string) func(tx neo4j.Transaction) (interface{}, error) {
	return func(tx neo4j.Transaction) (interface{}, error) {
		records, err := tx.Run("MATCH (n:User { username: $username })<--(user) RETURN user.username", map[string]interface{}{
			"username": username,
		})
		if err != nil {
			fmt.Println("couldn't retrieve followings", err)
			return nil, err
		}

		slice, err := records.Collect()
		if err != nil {
			fmt.Println("error getting followings: ", err)
			return nil, err
		}

		followings := []models.Following{}
		for _, entry := range slice {
			fmt.Println(entry.Values[0])

			following := models.NewFollowing(entry.Values[0].(string))
			followings = append(followings, following)
		}

		return followings, nil
	}
}

func (n *neo4jConnector) GetSuggestions(username string) func(tx neo4j.Transaction) (interface{}, error) {
	return func(tx neo4j.Transaction) (interface{}, error) {
		records, err := tx.Run(`MATCH (n:User)-[r:FOLLOWS*2]->(m:User)
								WHERE n.username = $username 
								RETURN distinct m AS Suggestion
								UNION
								MATCH (n:User), (m:User)
								WHERE NOT (n)-[:FOLLOWS]->(m) AND NOT n.username = m.username AND n.username = $username 
								RETURN m AS Suggestion
								LIMIT 10;`,
			map[string]interface{}{
				"username": username,
			})
		if err != nil {
			fmt.Println("couldn't retrieve suggestions", err)
			return nil, err
		}

		slice, err := records.Collect()
		if err != nil {
			fmt.Println("error getting suggestions: ", err)
			return nil, err
		}

		suggestions := []models.Suggestion{}
		for _, entry := range slice {
			suggestion := models.NewSuggestion(entry.Values[0].(string))
			suggestions = append(suggestions, suggestion)
		}

		return suggestions, nil
	}
}
