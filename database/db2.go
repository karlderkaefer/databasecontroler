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

type db2 struct {
}

func (db *db2) CreateDockerDb2Command(commands string) *exec.Cmd {
	path, err := exec.LookPath("docker")
	if err != nil {
		log.Fatal("could not find docker installed")
	}
	baseCommand := "docker exec --user db2inst1 db2 /home/db2inst1/sqllib/bin/db2"
	cmd := &exec.Cmd{
		Path: path,
		Args: append(strings.Fields(baseCommand), strings.Fields(commands)...),
	}
	return cmd
}

func (db *db2) Config() configuration {
	return configuration{
		DriverClass: "go_ibm_db",
		Host:        "localhost",
		Port:        50000,
		Username:    "db2inst1",
		Password:    "db2inst1-pwd",
		Instance:    "sample",
	}
}

func (db *db2) Connect() (*sqlx.DB, Message, error) {
	return nil, Message{}, errors.New("not implemented")
}

func (db *db2) Execute(command string) (Message, error) {
	var message Message
	out, err := db.CreateDockerDb2Command(command).CombinedOutput()
	if err != nil {
		log.Printf("%v", err)
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

func (db *db2) ConnectionURL() string {
	return ""
}

func (db *db2) CreateUser(username string, password string) ([]Message, error) {
	var messages []Message
	if len(username) > 8 {
		return nil, fmt.Errorf(NameMaxLength, 8)
	}
	msg, err := db.Execute(fmt.Sprintf("create database %s PAGESIZE 16384", username))
	if err != nil {
		if strings.Contains(err.Error(), "SQL1005N") {
			warning := fmt.Sprintf(UserAlreadyExists, username, err.Error())
			return addWarn(messages, warning), nil
		}
		return addError(messages, err)
	}
	return addSuccess(messages, msg.Content)
}

func (db *db2) DropUser(username string) ([]Message, error) {
	var messages []Message
	db.Execute("catalog database " + username)
	msg, err := db.Execute("drop database " + username)
	db.Execute("uncatalog database " + username)
	if err != nil {
		if strings.Contains(err.Error(), "SQL1013N") {
			message := fmt.Sprintf(UserNotExists, username, err.Error())
			return addWarn(messages, message), nil
		}
		return addError(messages, err)
	}
	return addSuccess(messages, fmt.Sprintf("user %s dropped: %s", username, msg.Content))
}

func (db *db2) RecreateUser(username string, password string) ([]Message, error) {
	return recreateUser(db, username, password)
}

func (db *db2) ListUsers() ([]SystemUser, error) {
	msg, err := db.Execute("list database directory")
	if err != nil {
		return nil, err
	}
	return db.ParseDatabaseDirectoryList(msg.Content), nil
}

func (db *db2) ParseDatabaseDirectoryList(input string) []SystemUser {
	var users []SystemUser
	regex := regexp.MustCompile(`(?m)Database name\s+=\s+(?P<Name>\w+)`)
	found := regex.FindAllStringSubmatch(input, -1)
	for _, name := range found {
		user := &SystemUser{strings.ToLower(name[1])}
		if user.Username != "sample" {
			users = append(users, *user)
		}
	}
	return users
}
