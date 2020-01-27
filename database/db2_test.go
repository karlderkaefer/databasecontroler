package database

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"strings"
	"testing"
)

func TestDb2_CreateUser(t *testing.T) {
	testUser := "testuser1"
	testPass := "testpass1"
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	db, err := GetDatabase(Db2105)
	assert.Nil(t, err)
	resp, err := db.CreateUser(testUser, testPass)
	defer db.DropUser(testUser)
	assert.Nil(t, err)
	assert.Equal(t, resp[0].Severity, Success)
	assert.Contains(t, resp[0].Content, "CREATE DATABASE command completed successfully")

	// test user already exists
	resp, err = db.CreateUser(testUser, testPass)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp, 1)
	assert.Contains(t, resp[0].Content, fmt.Sprintf(UserAlreadyExists, testUser, ""))
}

func TestDb2_DropUser(t *testing.T) {
	testUser := "testuser2"
	testPass := "testpass2"
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	db, err := GetDatabase(Db2105)
	assert.Nil(t, err)

	resp, err := db.DropUser(testUser)

	assert.Nil(t, err)
	assert.Contains(t, resp[0].Content, fmt.Sprintf(UserNotExists, testUser, ""))
	_, err = db.CreateUser(testUser, testPass)
	assert.Nil(t, err)

	resp, err = db.DropUser(testUser)
	assert.Nil(t, err)
	assert.Equal(t, fmt.Sprintf(UserDropped, testUser), resp[0].Content, testUser)
}

func TestDb2_CreateDb2Command(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	db2 := new(Db2)
	cmd := db2.CreateDb2Command("hello")
	assert.Equal(t, "docker", cmd.Path)
	assert.Contains(t, cmd.Args, "hello")
	log.Print(cmd.Path, cmd.Args)
}

func TestDb2_ListUsers(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	testUser1 := "user1"
	testUser2 := "user2"
	testPassword := "testpass"
	db, err := GetDatabase(Db2105)
	assert.Nil(t, err)
	_, _ = db.DropUser(testUser1)
	_, _ = db.DropUser(testUser2)
	resp, err := db.CreateUser(testUser1, testPassword)
	defer db.DropUser(testUser1)
	t.Log(resp)
	assert.Nil(t, err)
	resp, err = db.CreateUser(testUser2, testPassword)
	defer db.DropUser(testUser2)
	t.Log(resp)
	assert.Nil(t, err)
	expected := []SystemUser{{strings.ToUpper(testUser1)}, {strings.ToUpper(testUser2)}}
	users, err := db.ListUsers()
	assert.Nil(t, err)
	assert.Equal(t, expected, users, "Expecting to find two users as listed in %v", users)
}

func TestDb2_ParseDatabaseDirectoryList(t *testing.T) {
	input := `
Database 1 entry:

 Database alias                       = TEST
 Database name                        = TEST
 Local database directory             = /home/db2inst1
 Database release level               = 10.00
 Comment                              =
 Directory entry type                 = Indirect
 Catalog database partition number    = 0
 Alternate server hostname            =
 Alternate server port number         =

Database 2 entry:

 Database alias                       = TEST2
 Database name                        = TEST2
 Local database directory             = /home/db2inst1
 Database release level               = 10.00
 Comment                              =
 Directory entry type                 = Indirect
 Catalog database partition number    = 0
 Alternate server hostname            =
 Alternate server port number         =
`
	db2 := new(Db2)
	expected := []SystemUser{{"test"}, {"test2"}}
	users := db2.ParseDatabaseDirectoryList(input)
	assert.ElementsMatch(t, expected, users)
}
