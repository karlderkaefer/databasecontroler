package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"log"
	"os/exec"
	"regexp"
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
	log.Printf("%s", out)
	return message, nil
}

func (db *Db2) ConnectionUrl() string {
	return ""
}

func (db *Db2) CreateUser(username string, password string) ([]Message, error) {
	var messages []Message
	msg, err := db.Execute(fmt.Sprintf("create database %s PAGESIZE 16384", username))
	if err != nil {
		if strings.Contains(err.Error(), "SQL1005N") {
			return addError(messages, &UserAlreadyExistsError{username, "Cannot create user."})
		} else {
			return addError(messages, err)
		}
	}
	return addSuccess(messages, msg.Content)
}

func (db *Db2) DropUser(username string) ([]Message, error) {
	var messages []Message
	db.Execute("catalog database " + username)
	msg, err := db.Execute("drop database " + username)
	db.Execute("uncatalog database " + username)
	if err != nil {
		if strings.Contains(err.Error(), "SQL1013N") {
			warning := fmt.Sprintf("user %s does not exist: %s", username, err.Error())
			return addWarn(messages, warning), nil
		}
		return addError(messages, err)
	}
	return addSuccess(messages, fmt.Sprintf("user %s dropped: %s", username, msg.Content))
}

func (db *Db2) RecreateUser(username string, password string) ([]Message, error) {
	return recreateUser(db, username, password)
}

func (db *Db2) ListUsers() ([]SystemUser, error) {
	msg, err := db.Execute("list database directory")
	if err != nil {
		return nil, err
	}
	return db.ParseDatabaseDirectoryList(msg.Content), nil
}

func (db *Db2) ParseDatabaseDirectoryList(input string) []SystemUser {
	var users []SystemUser
	regex := regexp.MustCompile(`(?m)Database name\s+=\s+(?P<Name>\w+)`)
	found := regex.FindAllStringSubmatch(input, -1)
	for _, name := range found {
		user := &SystemUser{strings.ToLower(name[1])}
		users = append(users, *user)
	}
	return users
}
