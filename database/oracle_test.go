package database

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestCreateUser(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	oracleVersions := []int{Oracle11, Oracle12}
	for _, oracleVersion := range oracleVersions {
		oracle, err := GetDatabase(oracleVersion)
		if err != nil {
			t.Error(err)
		}
		resp, err := oracle.CreateUser("testusercreate", "testpass")
		t.Logf("%v", resp)
		if err != nil {
			t.Error(err)
		}
		expectMessage := "user created testusercreate"
		if resp[0].Content != expectMessage {
			t.Errorf("expected message: %s but was %s", expectMessage, resp[0].Content)
		}
		// user already exists
		resp, err = oracle.CreateUser("testusercreate", "testpass")
		if err != nil {
			t.Error(err)
		}
		expectMessage = "user testusercreate already exists"
		if !strings.Contains(resp[0].Content, expectMessage) {
			t.Errorf("expected message: %s but was %s", expectMessage, resp[0].Content)
		}
		oracle.DropUser("testusercreate")
	}
}

func TestDropUser(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	oracleVersions := []int{Oracle11, Oracle12}
	for _, oracleVersion := range oracleVersions {
		oracle, err := GetDatabase(oracleVersion)
		if err != nil {
			t.Error(err)
		}
		resp, err := oracle.DropUser("testusercreate")
		expectMessage := "user testusercreate does not exists"
		if !strings.Contains(resp[0].Content, expectMessage) {
			t.Errorf("expected message: %s but was %s", expectMessage, resp[0].Content)
		}
		_, err = oracle.CreateUser("testusercreate", "testpass")
		if err != nil {
			t.Error(err)
		}
		resp, err = oracle.DropUser("testusercreate")
		expectMessage = "user dropped testusercreate"
		if resp[0].Content != expectMessage {
			t.Errorf("expected message: %s but was %s", expectMessage, resp[0].Content)
		}
	}
}

func TestListUsers(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	oracleVersions := []int{Oracle11, Oracle12}
	for _, oracleVersion := range oracleVersions {
		oracle, err := GetDatabase(oracleVersion)
		if err != nil {
			t.Error(err)
		}

		oracle.DropUser("user1")
		oracle.DropUser("user2")
		oracle.CreateUser("user1", "testpass")
		oracle.CreateUser("user2", "testpass")
		expected := []SystemUser{{"USER1"}, {"USER2"}}
		resp, err := oracle.ListUsers()
		assert.Equal(t, expected, resp, "Expecting to find two users as listed in %s", oracle)

	}
}
