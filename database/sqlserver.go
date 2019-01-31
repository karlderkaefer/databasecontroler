package database

import (
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/jmoiron/sqlx"
	"log"
	"strings"
)

type Sqlserver struct {
}

type SqlserverUser struct {
	Username string `db:"name"`
}

type UserTemplate struct {
	User string
}

func (db *Sqlserver) Config() Configuration {
	return Configuration{
		DriverClass: "sqlserver",
		Host:        "localhost",
		Port:        1433,
		Username:    "sa",
		Password:    "HelloApes66",
		Instance:    "",
	}
}

func (db *Sqlserver) Connect() (*sqlx.DB, Message, error) {
	var message Message
	con, err := sqlx.Connect(db.Config().DriverClass, db.ConnectionUrl())
	if err != nil {
		message = Message{
			Severity: Error,
			Content:  err.Error(),
		}
		log.Print(err)
		return nil, message, err
	}
	return con, Message{}, err
}

func (db *Sqlserver) Execute(command string) (Message, error) {
	var message Message
	con, message, err := db.Connect()
	if err != nil {
		return message, err
	}
	defer con.Close()
	_, err = con.Exec(command)
	return message, err
}

func (db *Sqlserver) ConnectionUrl() string {
	return fmt.Sprintf("sqlserver://%s:%s@%s:%d",
		db.Config().Username,
		db.Config().Password,
		db.Config().Host,
		db.Config().Port,
	)
}

func (db *Sqlserver) CreateUser(username string, password string) ([]Message, error) {
	var messages []Message
	value := &TemplateValue{User: username}
	createUserSql, err := LoadTemplate(*value, TemplateSqlServerCreate)
	if err != nil {
		return messages, nil
	}
	_, err = db.Execute(createUserSql)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			warning := fmt.Sprintf("user %s already exists: %s", username, err.Error())
			return addWarn(messages, warning), nil
		}
		return addError(messages, err)
	}
	return addSuccess(messages, "user created " + username)
}

func (db *Sqlserver) DropUser(username string) ([]Message, error) {
	var messages []Message
	value := &TemplateValue{User: username}
	dropUserSql, err := LoadTemplate(*value, TemplateSqlServerDrop)
	if err != nil {
		return messages, err
	}
	_, err = db.Execute(dropUserSql)
	if err != nil {
		if strings.Contains(err.Error(), "ALTER DATABASE statement failed") {
			warning := fmt.Sprintf("user %s does not exists: %s", username, err.Error())
			return addWarn(messages, warning), nil
		}
		return addError(messages, err)
	}
	return addSuccess(messages, fmt.Sprintf("user %s dropped", username))
}

func (db *Sqlserver) RecreateUser(username string, password string) ([]Message, error) {
	return recreateUser(db, username, password)
}

func (db *Sqlserver) ListUsers() ([]SystemUser, error) {
	con, _, err := db.Connect()
	if err != nil {
		return nil, err
	}
	defer con.Close()
	var SqlserverUsers []SqlserverUser
	sql := "SELECT name FROM sys.databases where database_id > 4;"
	con.Select(&SqlserverUsers, sql)
	var users []SystemUser
	for _, value := range SqlserverUsers {
		user := &SystemUser{Username: value.Username}
		users = append(users, *user)
	}
	return users, nil
}
