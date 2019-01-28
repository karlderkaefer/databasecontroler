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
			messages, err = addWarn(messages, warning)
		} else {
			messages, err = addError(messages, err)
		}
		return messages, err
	}
	createDatabaseSql := fmt.Sprintf("CREATE DATABASE %s;", username)
	_, err = db.Execute(createDatabaseSql)
	if err != nil {
		messages, err = addError(messages, err)
		return messages, err
	}
	grantPrivileges := fmt.Sprintf("GRANT ALL PRIVILEGES ON %s.* TO '%s'@'%%';FLUSH PRIVILEGES;", username, username)
	_, err = db.Execute(grantPrivileges)
	if err != nil {
		messages, err = addError(messages, err)
		return messages, err
	}
	return addSuccess(messages, "user created " + username)
}

func (db *Mysql) DropUser(username string) ([]Message, error) {
	var messages []Message
	dropUserSql := fmt.Sprintf("DROP USER %s;", username)
	_, err := db.Execute(dropUserSql)
	if err != nil {
		if strings.Contains(err.Error(), "Error 1396") {
			warning := fmt.Sprintf("user %s does not exists: %s", username, err.Error())
			messages, err = addWarn(messages, warning)
		} else {
			messages, err = addError(messages, err)
		}
		return messages, err
	}
	dropDatabaseSql := fmt.Sprintf("DROP DATABASE %s;", username)
	_, err = db.Execute(dropDatabaseSql)
	if err != nil {
		messages, err = addError(messages, err)
		return messages, err
	}
	return addSuccess(messages, fmt.Sprintf("user %s dropped", username))
}

func (db *Mysql) RecreateUser(username string, password string) ([]Message, error) {
	messages := make([]Message, 0)
	msg, err := db.DropUser(username)
	if err != nil {
		return nil, err
	}
	messages = append(messages, msg...)
	msg, err = db.CreateUser(username, password)
	if err != nil {
		return messages, err
	}
	messages = append(messages, msg...)
	return messages, nil
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

func addError(messages []Message, err error) ([]Message, error) {
	message := Message{
		Severity: Error,
		Content:  err.Error(),
	}
	messages = append(messages, message)
	return messages, err
}

func addSuccess(messages []Message, content string) ([]Message, error) {
	message := Message{
		Severity: Success,
		Content:  content,
	}
	messages = append(messages, message)
	return messages, nil
}

func addWarn(messages []Message, content string) ([]Message, error) {
	message := Message{
		Severity: Warn,
		Content:  content,
	}
	messages = append(messages, message)
	return messages, nil
}
