package database

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func oracleVersions() []int {
	return []int{Oracle11}
}

func TestCreateUser(t *testing.T) {
	for _, oracleVersion := range oracleVersions() {

		testUser := "testusercreate"
		testPassword := "testpass"

		oracle, err := GetDatabase(oracleVersion)
		assert.Nil(t, err)
		resp, err := oracle.CreateUser(testUser, testPassword)
		t.Logf("%v", resp)
		assert.Nil(t, err)

		assert.NotNil(t, resp)
		assert.Len(t, resp, 1)
		assert.Equal(t, fmt.Sprintf(UserCreated, testUser), resp[0].Content)

		// user already exists
		resp, err = oracle.CreateUser(testUser, testPassword)
		assert.Nil(t, err)

		assert.NotNil(t, resp)
		assert.Len(t, resp, 1)
		assert.Contains(t, resp[0].Content, fmt.Sprintf(UserAlreadyExists, testUser, ""))

		oracle.DropUser(testUser)
	}
}

func TestDropUser(t *testing.T) {
	for _, oracleVersion := range oracleVersions() {

		testUser := "testusercreate"
		testPassword := "testpass"

		oracle, err := GetDatabase(oracleVersion)
		assert.Nil(t, err)
		resp, err := oracle.DropUser(testUser)

		assert.Nil(t, err)
		assert.Contains(t, resp[0].Content, fmt.Sprintf(UserNotExists, testUser, ""))
		_, err = oracle.CreateUser(testUser, testPassword)
		assert.Nil(t, err)

		resp, err = oracle.DropUser(testUser)
		assert.Nil(t, err)
		assert.Equal(t, fmt.Sprintf(UserDropped, testUser), resp[0].Content, testUser)

	}
}

func TestListUsers(t *testing.T) {
	for _, oracleVersion := range oracleVersions() {
		testUser1 := "user1"
		testUser2 := "user2"
		testPassword := "testpass"
		oracle, err := GetDatabase(oracleVersion)
		assert.Nil(t, err)
		oracle.DropUser(testUser1)
		oracle.DropUser(testUser2)
		oracle.CreateUser(testUser1, testPassword)
		oracle.CreateUser(testUser2, testPassword)
		expected := []SystemUser{{"USER1"}, {"USER2"}}
		resp, err := oracle.ListUsers()
		assert.Nil(t, err)
		assert.Equal(t, expected, resp, "Expecting to find two users as listed in %s", oracle)

	}
}
