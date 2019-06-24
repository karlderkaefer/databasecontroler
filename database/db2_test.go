package database

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestDb2_CreateUser(t *testing.T) {
	db, err := GetDatabase(Db2105)
	if err != nil {
		t.Error(err)
	}
	_, err = db.CreateUser("test2", "test2")
	if err != nil {
		t.Error(err)
	}
	//assert.Equal(t, msg[0].Content, "sad")
}

func TestDb2_CreateDb2Command(t *testing.T) {
	db2 := new(Db2)
	cmd := db2.CreateDb2Command("hello")
	assert.Equal(t, "docker", cmd.Path)
	assert.Contains(t, cmd.Args, "hello")
	log.Print(cmd.Path, cmd.Args)
}
