package database

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSqlserver_CreateUser(t *testing.T) {
	db, err := GetDatabase(SqlServer2017)
	if err != nil {
		t.Error(err)
	}
	resp, err := db.CreateUser("testusercreate", "testpass")
	t.Logf("%v", resp)
	if err != nil {
		t.Error(err)
	}
	expectMessage := "user created testusercreate"
	assert.Equal(t, expectMessage, resp[0].Content)
	// user already exists
	resp, err = db.CreateUser("testusercreate", "testpass")
	if err != nil {
		t.Error(err)
	}
	expectMessage = "user testusercreate already exists"
	assert.Contains(t, resp[0].Content, expectMessage)
	db.DropUser("testusercreate")
}

func TestSqlserver_DropUser(t *testing.T) {
	db, err := GetDatabase(SqlServer2017)
	if err != nil {
		t.Error(err)
	}
	// user does not exists
	resp, err := db.DropUser("testuserdrop")
	t.Logf("%v", resp)
	if err != nil {
		t.Error(err)
	}
	expectMessage := "user testuserdrop does not exists"
	assert.Contains(t, resp[0].Content, expectMessage)
	_, err = db.CreateUser("testuserdrop", "testpass")
	if err != nil {
		t.Error(err)
	}
	resp, err = db.DropUser("testuserdrop")
	if err != nil {
		t.Error(err)
	}
	expectMessage = "user testuserdrop dropped"
	assert.Equal(t, expectMessage, resp[0].Content)
}

func TestSqlserver_ListUsers(t *testing.T) {
	db, err := GetDatabase(SqlServer2017)
	if err != nil {
		t.Error(err)
	}
	db.DropUser("user1")
	db.DropUser("user2")
	db.CreateUser("user1", "testpass")
	db.CreateUser("user2", "testpass")
	expected := []SystemUser{{"user1"}, {"user2"}}
	resp, err := db.ListUsers()
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, expected, resp, "Expecting to find two users as listed in %s", db.Config().DriverClass)
}


