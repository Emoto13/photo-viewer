package tests

import (
	"testing"

	"github.com/Emoto13/photo-viewer-rest/follow-service/src/follow"
	"github.com/Emoto13/photo-viewer-rest/follow-service/src/follow/models"
	"github.com/Emoto13/photo-viewer-rest/follow-service/tests/setup"
)

func assertEqualsError(a, b error, message string) {
	if a != b {
		panic(message)
	}
}

func TestSaveFollow(t *testing.T) {
	db, _ := setup.OpenPostgresDatabaseConnection()
	db.Exec(setup.CreateUsersTable)
	db.Exec(setup.CreateFollowersTable)
	db.Exec(setup.CreateTestUser)
	db.Exec(setup.CreateTestFollower)

	followStore := follow.NewFollowStore(db)

	var tests = []struct {
		name     string
		input    *models.Follow
		expected error
		message  string
	}{
		{"Test SaveFollow", models.NewFollow("TestFollower", "TestUser"), nil, "SaveFollow should work correctly"},
	}

	for _, test := range tests {
		err := followStore.SaveFollow(test.input)
		assertEqualsError(err, test.expected, test.message)
	}

	db.Exec(setup.DropUsersTable)
	db.Exec(setup.DropFollowersTable)
}

func TestRemoveFollow(t *testing.T) {
	db, _ := setup.OpenPostgresDatabaseConnection()
	db.Exec(setup.CreateUsersTable)
	db.Exec(setup.CreateFollowersTable)
	db.Exec(setup.CreateTestUser)
	db.Exec(setup.CreateTestFollower)
	db.Exec(setup.FollowTestUser)

	followStore := follow.NewFollowStore(db)

	var tests = []struct {
		name     string
		input    *models.Follow
		expected error
		message  string
	}{
		{"Test RemoveFollow", models.NewFollow("TestFollower", "TestUser"), nil, "RemoveFollow should work correctly"},
	}

	for _, test := range tests {
		err := followStore.RemoveFollow(test.input)
		assertEqualsError(err, test.expected, test.message)
	}

	db.Exec(setup.DropUsersTable)
	db.Exec(setup.DropFollowersTable)
}

func TestGetFollowers(t *testing.T) {
	db, _ := setup.OpenPostgresDatabaseConnection()
	db.Exec(setup.CreateUsersTable)
	db.Exec(setup.CreateFollowersTable)
	db.Exec(setup.CreateTestUser)
	db.Exec(setup.CreateTestFollower)
	db.Exec(setup.FollowTestUser)

	followStore := follow.NewFollowStore(db)

	var tests = []struct {
		name     string
		input    string
		expected error
		message  string
	}{
		{"Test GetFollowers", "TestUser", nil, "GetFollowers should work correctly"},
	}

	for _, test := range tests {
		_, err := followStore.GetFollowers(test.input)
		assertEqualsError(err, test.expected, test.message)
	}
	db.Exec(setup.DropUsersTable)
	db.Exec(setup.DropFollowersTable)
}

func TestGetFollowing(t *testing.T) {
	db, _ := setup.OpenPostgresDatabaseConnection()
	db.Exec(setup.CreateUsersTable)
	db.Exec(setup.CreateFollowersTable)
	db.Exec(setup.CreateTestUser)
	db.Exec(setup.CreateTestFollowing)
	db.Exec(setup.FollowTestFollowing)

	followStore := follow.NewFollowStore(db)

	var tests = []struct {
		name     string
		input    string
		expected error
		message  string
	}{
		{"Test GetFollowing", "TestFollowing", nil, "GetFollowing should work correctly"},
	}

	for _, test := range tests {
		_, err := followStore.GetFollowing(test.input)
		assertEqualsError(err, test.expected, test.message)
	}
	db.Exec(setup.DropUsersTable)
	db.Exec(setup.DropFollowersTable)
}

func TestGetSuggestions(t *testing.T) {
	db, _ := setup.OpenPostgresDatabaseConnection()
	db.Exec(setup.CreateUsersTable)
	db.Exec(setup.CreateFollowersTable)
	db.Exec(setup.CreateTestUser)
	db.Exec(setup.CreateTestFollowing)

	followStore := follow.NewFollowStore(db)
	var tests = []struct {
		name     string
		input    string
		expected error
		message  string
	}{
		{"Test GetSuggestions", "TestUser", nil, "GetSuggestions should work correctly"},
	}

	for _, test := range tests {
		_, err := followStore.GetSuggestions(test.input)
		assertEqualsError(err, test.expected, test.message)
	}

	db.Exec(setup.DropUsersTable)
	db.Exec(setup.DropFollowersTable)
}
