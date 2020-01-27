package database

import (
	"fmt"
	// driver is only needed on runtime not on compile time
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"log"
	"strings"
)

type mysqlDatabase struct {
}

type mysqlUser struct {
	Username string `db:"User"`
}

func (db *mysqlDatabase) Config() configuration {
	return configuration{
		DriverClass: "mysql",
		Host:        "localhost",
		Port:        3306,
		Username:    "root",
		Password:    "HelloApes66",
		Instance:    "seed",
	}
}

func (db *mysqlDatabase) Connect() (*sqlx.DB, Message, error) {
	var message Message
	con, err := sqlx.Connect(db.Config().DriverClass, db.ConnectionURL())
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

func (db *mysqlDatabase) Execute(command string) (Message, error) {
	var message Message
	con, message, err := db.Connect()
	if err != nil {
		return message, err
	}
	defer con.Close()
	_, err = con.Exec(command)
	return message, err
}

func (db *mysqlDatabase) ConnectionURL() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?multiStatements=true",
		db.Config().Username,
		db.Config().Password,
		db.Config().Host,
		db.Config().Port,
		db.Config().Instance,
	)
}

func (db *mysqlDatabase) CreateUser(username string, password string) ([]Message, error) {
	var messages []Message
	createUserSQL := fmt.Sprintf("CREATE USER %s IDENTIFIED BY '%s';", username, password)
	_, err := db.Execute(createUserSQL)
	if err != nil {
		if strings.Contains(err.Error(), "Error 1396") {
			warning := fmt.Sprintf("user %s already exists: %s", username, err.Error())
			return addWarn(messages, warning), nil
		}
		return addError(messages, err)
	}
	createDatabaseSQL := fmt.Sprintf("CREATE DATABASE %s;", username)
	_, err = db.Execute(createDatabaseSQL)
	if err != nil {
		return addError(messages, err)
	}
	grantPrivileges := fmt.Sprintf("GRANT ALL PRIVILEGES ON %s.* TO '%s'@'%%';FLUSH PRIVILEGES;", username, username)
	_, err = db.Execute(grantPrivileges)
	if err != nil {
		return addError(messages, err)
	}
	return addSuccess(messages, "user created "+username)
}

func (db *mysqlDatabase) DropUser(username string) ([]Message, error) {
	var messages []Message
	dropUserSQL := fmt.Sprintf("DROP USER %s;", username)
	_, err := db.Execute(dropUserSQL)
	if err != nil {
		if strings.Contains(err.Error(), "Error 1396") {
			warning := fmt.Sprintf("user %s does not exists: %s", username, err.Error())
			return addWarn(messages, warning), nil
		}
		return addError(messages, err)
	}
	dropDatabaseSQL := fmt.Sprintf("DROP DATABASE %s;", username)
	_, err = db.Execute(dropDatabaseSQL)
	if err != nil {
		return addError(messages, err)
	}
	return addSuccess(messages, fmt.Sprintf("user %s dropped", username))
}

func (db *mysqlDatabase) RecreateUser(username string, password string) ([]Message, error) {
	return recreateUser(db, username, password)
}

func (db *mysqlDatabase) ListUsers() ([]SystemUser, error) {
	con, _, err := db.Connect()
	if err != nil {
		return nil, err
	}
	defer con.Close()
	var mysqlUsers []mysqlUser
	sql := "SELECT User FROM mysql.user WHERE host = '%' and user != 'root';"
	con.Select(&mysqlUsers, sql)
	var users []SystemUser
	for _, value := range mysqlUsers {
		user := &SystemUser{Username: value.Username}
		users = append(users, *user)
	}
	return users, nil
}
