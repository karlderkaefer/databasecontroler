package database

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestDb2_CreateUser(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	db, err := GetDatabase(Db2105)
	if err != nil {
		t.Error(err)
	}
	defer db.DropUser("test2")
	msg, err := db.CreateUser("test2", "test2")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, msg[0].Severity, Success)
	assert.Contains(t, msg[0].Content, "CREATE DATABASE command completed successfully")

	_, err = db.CreateUser("test2", "test2")
	assert.IsType(t, err, new(UserAlreadyExistsError))
}

func TestDb2_DropUser(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	db, err := GetDatabase(Db2105)
	if err != nil {
		t.Error(err)
	}
	msg, err := db.DropUser("test3")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, msg[0].Severity, Warn)
	assert.Contains(t, msg[0].Content, "user test3 does not exist")

	db.CreateUser("test3", "test3")
	msg, err = db.DropUser("test3")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, msg[0].Severity, Success)
	assert.Contains(t, msg[0].Content, "user test3 dropped")
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
	db, err := GetDatabase(Db2105)
	if err != nil {
		t.Error(err)
	}
	//db.CreateUser("test2","test2")
	users, err := db.ListUsers()
	if err != nil {
		t.Error(err)
	}
	assert.Len(t, users, 1)
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
