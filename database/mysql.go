package database

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"log"
	"strings"
)

type Mysql struct {
}

type MysqlUser struct {
	Username string `db:"User"`
}

func (db *Mysql) Config() Configuration {
	return Configuration{
		DriverClass: "mysql",
		Host:        "localhost",
		Port:        3306,
		Username:    "root",
		Password:    "HelloApes66",
		Instance:    "seed",
	}
}

func (db *Mysql) Connect() (*sqlx.DB, Message, error) {
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

func (db *Mysql) Execute(command string) (Message, error) {
	var message Message
	con, message, err := db.Connect()
	if err != nil {
		return message, err
	}
	defer con.Close()
	_, err = con.Exec(command)
	return message, err
}

func (db *Mysql) ConnectionUrl() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?multiStatements=true",
		db.Config().Username,
		db.Config().Password,
		db.Config().Host,
		db.Config().Port,
		db.Config().Instance,
	)
}

func (db *Mysql) CreateUser(username string, password string) ([]Message, error) {
	var messages []Message
	createUserSql := fmt.Sprintf("CREATE USER %s IDENTIFIED BY '%s';", username, password)
	_, err := db.Execute(createUserSql)
	if err != nil {
		if strings.Contains(err.Error(), "Error 1396") {
			warning := fmt.Sprintf("user %s already exists: %s", username, err.Error())
			return addWarn(messages, warning), nil
		}
		return addError(messages, err)
	}
	createDatabaseSql := fmt.Sprintf("CREATE DATABASE %s;", username)
	_, err = db.Execute(createDatabaseSql)
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

func (db *Mysql) DropUser(username string) ([]Message, error) {
	var messages []Message
	dropUserSql := fmt.Sprintf("DROP USER %s;", username)
	_, err := db.Execute(dropUserSql)
	if err != nil {
		if strings.Contains(err.Error(), "Error 1396") {
			warning := fmt.Sprintf("user %s does not exists: %s", username, err.Error())
			return addWarn(messages, warning), nil
		}
		return addError(messages, err)
	}
	dropDatabaseSql := fmt.Sprintf("DROP DATABASE %s;", username)
	_, err = db.Execute(dropDatabaseSql)
	if err != nil {
		return addError(messages, err)
	}
	return addSuccess(messages, fmt.Sprintf("user %s dropped", username))
}

func (db *Mysql) RecreateUser(username string, password string) ([]Message, error) {
	return recreateUser(db, username, password)
}

func (db *Mysql) ListUsers() ([]SystemUser, error) {
	con, _, err := db.Connect()
	if err != nil {
		return nil, err
	}
	defer con.Close()
	var mysqlUsers []MysqlUser
	sql := "SELECT User FROM mysql.user WHERE host = '%' and user != 'root';"
	con.Select(&mysqlUsers, sql)
	var users []SystemUser
	for _, value := range mysqlUsers {
		user := &SystemUser{Username: value.Username}
		users = append(users, *user)
	}
	return users, nil
}
