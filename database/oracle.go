package database

import (
	"database/sql"
	"fmt"
	_ "github.com/godror/godror"
	"github.com/jmoiron/sqlx"
	"log"
	"strings"
)

type Oracle struct {
	version int
}

type OracleUser struct {
	Username string `db:"USERNAME"`
	Userid   string `db:"USER_ID"`
}

func (db *Oracle) KillSession(username string) error {
	var sid sql.NullString
	var serial sql.NullString
	con, _, err := db.Connect()
	if err != nil {
		return nil
	}
	defer con.Close()
	findSession := fmt.Sprintf(
		"select sid, serial# as serial from v$session where username = '%s'",
		username,
	)
	err = con.DB.QueryRow(findSession).Scan(&sid, &serial)
	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		return err
	}
	if sid.Valid && serial.Valid {
		killSession := fmt.Sprintf(
			"alter system kill session '%s,%s' IMMEDIATE",
			sid.String,
			serial.String,
		)
		_, err = db.Execute(killSession)
		if err != nil {
			log.Fatal(err)
			return err
		}
	}
	return nil
}

func (db *Oracle) DropUser(username string) ([]Message, error) {
	var messages []Message
	db.KillSession(username)
	dropUserSql := fmt.Sprintf("drop user %s cascade", username)
	_, err := db.Execute(dropUserSql)
	if err != nil {
		if strings.Contains(err.Error(), "ORA-01918") {
			message := fmt.Sprintf(UserNotExists, username, err.Error())
			return addWarn(messages, message), nil
		}
		return addError(messages, err)
	}
	return addSuccess(messages, fmt.Sprintf(UserDropped, username))
}

func (db *Oracle) RecreateUser(username string, password string) ([]Message, error) {
	return recreateUser(db, username, password)
}

func (db *Oracle) ListUsers() ([]SystemUser, error) {
	con, _, err := db.Connect()
	if err != nil {
		return nil, err
	}
	defer con.Close()
	var oracleUsers []OracleUser
	var sql string
	switch db.version {
	case 11:
		sql = "SELECT username, user_id FROM dba_users WHERE user_id > 50 AND user_id < 1000000"
	case 12:
		sql = "SELECT username, user_id FROM dba_users WHERE ORACLE_MAINTAINED = 'N' AND username != 'PDBADMIN' ORDER BY username"
	}
	con.Select(&oracleUsers, sql)
	//log.Printf("%v", oracleUsers)
	// mapping to system user
	var users []SystemUser
	for _, value := range oracleUsers {
		user := &SystemUser{Username: value.Username}
		users = append(users, *user)
	}
	return users, nil
}

func (db *Oracle) CreateUser(username string, password string) ([]Message, error) {
	messages := make([]Message, 0)
	createUserSql := fmt.Sprintf("create user %s identified by %s", username, password)
	_, err := db.Execute(createUserSql)
	if err != nil {
		// user already exists
		if strings.Contains(err.Error(), "ORA-01920") {
			warning := fmt.Sprintf(UserAlreadyExists, username, err.Error())
			return addWarn(messages, warning), nil
		}
		log.Print(err)
		return addError(messages, err)
	}
	_, err = db.Execute(fmt.Sprintf("grant all privileges to %s", username))
	if err != nil {
		addError(messages, err)
	}
	return addSuccess(messages, fmt.Sprintf(UserCreated, username))
}

func (db *Oracle) Connect() (*sqlx.DB, Message, error) {
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

func (db *Oracle) Execute(command string) (Message, error) {
	var message Message
	con, message, err := db.Connect()
	if err != nil {
		return message, err
	}
	defer con.Close()
	_, err = con.Exec(command)
	return message, err
}

func (db *Oracle) ConnectionUrl() string {
	return fmt.Sprintf(
		"%s/%s@%s:%d/%s",
		db.Config().Username,
		db.Config().Password,
		db.Config().Host,
		db.Config().Port,
		db.Config().Instance)
}

func (db *Oracle) Config() Configuration {
	switch db.version {
	case Oracle11:
		return Configuration{
			DriverClass: "godror",
			Host:        "localhost",
			Instance:    "xe",
			Password:    "HelloApes66",
			Port:        1521,
			Username:    "system",
		}
	case Oracle12:
		return Configuration{
			DriverClass: "godror",
			Host:        "localhost",
			Instance:    "ORCLPDB1",
			Password:    "HelloApes66",
			Port:        1522,
			Username:    "system",
		}
	default:
		return Configuration{}
	}
}
