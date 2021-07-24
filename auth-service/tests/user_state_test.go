package unit

import (
	"testing"

	"github.com/Emoto13/photo-viewer-rest/auth-service/src/user"
	"github.com/Emoto13/photo-viewer-rest/auth-service/tests/setup"
)

func TestUserRetrieval(t *testing.T) {
	db, _ := setup.OpenPostgresDatabaseConnection()
	db.Exec(setup.CreateUsersTable)
	db.Exec(setup.CreateTestUser)
	state := user.NewState(db)

	testUser, _ := user.NewUser("TestUser", "TestPassword")
	var tests = []struct {
		name     string
		input    string
		expected *user.User
		message  string
	}{
		{"Test RetrieveUser", "TestUser", testUser, "User state should retrieve user by username successfully"},
	}

	for _, test := range tests {
		retrievedUser, _ := state.RetrieveUser(test.input)
		assertEquals(retrievedUser.Username, test.expected.Username, test.message)
	}

	db.Exec(setup.DropUsersTable)
}
