package database

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"os/exec"
	"strings"
	"testing"
)

func TestDb2_CreateUserTooLong(t *testing.T) {
	db, err := GetDatabase(Db2105)
	assert.Nil(t, err)
	resp, err := db.CreateUser("longerthan8chars", "testpass")
	assert.NotNil(t, err)
	assert.Nil(t, resp)
	assert.Error(t, err, fmt.Sprintf(NameMaxLength, 8))
	resp, err = db.CreateUser("aaaaaaaaa", "testpass")
	assert.NotNil(t, err)
	assert.Nil(t, resp)
	assert.Error(t, err, fmt.Sprintf(NameMaxLength, 8))
}

func TestDb2_CreateAndListUser(t *testing.T) {
	testUser := "testuse1"
	testPass := "testpass"

	db, err := GetDatabase(Db2105)
	assert.Nil(t, err)
	resp, err := db.CreateUser(testUser, testPass)
	defer db.DropUser(testUser)
	assert.Nil(t, err)
	assert.Equal(t, Success, resp[0].Severity)
	assert.Contains(t, resp[0].Content, "CREATE DATABASE command completed successfully")

	// test user already exists
	resp, err = db.CreateUser(testUser, testPass)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp, 1)
	assert.Contains(t, resp[0].Content, fmt.Sprintf(UserAlreadyExists, testUser, ""))

	// test list users here because of performance
	assert.Nil(t, err)
	expected := []SystemUser{{testUser}}
	users, err := db.ListUsers()
	assert.Nil(t, err)
	assert.Equal(t, expected, users, "Expecting to find two users as listed in %v", users)
}

func TestDb2_DropUser(t *testing.T) {
	testUser := "testuse2"
	testPass := "testpass2"

	db, err := GetDatabase(Db2105)
	assert.Nil(t, err)

	resp, err := db.DropUser(testUser)

	assert.Nil(t, err)
	assert.Contains(t, resp[0].Content, fmt.Sprintf(UserNotExists, testUser, ""))
	_, err = db.CreateUser(testUser, testPass)
	assert.Nil(t, err)

	resp, err = db.DropUser(testUser)
	assert.Nil(t, err)
	assert.Contains(t, resp[0].Content, "The DROP DATABASE command completed successfully")
}

func TestDb2_CreateDockerDb2Command(t *testing.T) {
	db2 := new(Db2)
	cmd := db2.CreateDockerDb2Command("hello")

	path, err := exec.LookPath("docker")
	assert.Nil(t, err)
	expect := &exec.Cmd{
		Path: path,
		Args: strings.Fields("docker exec --user db2inst1 db2 /home/db2inst1/sqllib/bin/db2 hello"),
	}
	assert.Equal(t, expect, cmd)
	log.Printf("%v", cmd)
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
