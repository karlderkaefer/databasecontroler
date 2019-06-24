package database

import (
	_ "github.com/ibmdb/go_ibm_db"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"log"
	"os/exec"
	"strings"
)

type Db2 struct {
}

type Db2User struct {
	Username string `db:"name"`
}

func (db *Db2) CreateDb2Command(commands string) *exec.Cmd {
	baseCommand := "docker exec --user db2inst1 databasemanager_db2_1 /home/db2inst1/sqllib/bin/db2"
	cmd := &exec.Cmd{
		Path: "docker",
		Args: append(strings.Fields(baseCommand), strings.Fields(commands)...),
	}
	return cmd
}

func (db *Db2) Config() Configuration {
	return Configuration{
		DriverClass: "go_ibm_db",
		Host:        "localhost",
		Port:        50000,
		Username:    "db2inst1",
		Password:    "db2inst1-pwd",
		Instance:    "sample",
	}
}

func (db *Db2) Connect() (*sqlx.DB, Message, error) {
	return nil, Message{}, errors.New("not implemented")
}

func (db *Db2) Execute(command string) (Message, error) {
	var message Message
	out, err := db.CreateDb2Command(command).CombinedOutput()
	if err != nil {
		message = Message{
			Severity: Error,
			Content:  string(out),
		}
		return message, errors.New(string(out))
	}
	message = Message{
		Severity: Info,
		Content:  string(out),
	}
	log.Printf("%v", out)
	return message, nil
}

func (db *Db2) ConnectionUrl() string {
	return ""
}

func (db *Db2) CreateUser(username string, password string) ([]Message, error) {
	var messages []Message
	msg, err := db.Execute("create database " + username)
	if err != nil {
		addError(messages, err)
		return messages, err
	}
	addSuccess(messages, msg.Content)
	return messages, nil
}

func (db *Db2) DropUser(username string) ([]Message, error) {
	return []Message{}, errors.New("not implemented")
}

func (db *Db2) RecreateUser(username string, password string) ([]Message, error) {
	return []Message{}, errors.New("not implemented")
}

func (db *Db2) ListUsers() ([]SystemUser, error) {
	return []SystemUser{}, errors.New("not implemented")
}
