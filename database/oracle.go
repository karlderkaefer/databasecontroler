package database

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
)

type OracleHandler struct {
	version int
}

type Oracle struct {
	handler OracleHandler
}

type OracleUser struct {
	Username string `db:"USERNAME"`
	Userid   string `db:"USER_ID"`
}

func NewOracle12() DatabaseApi {
	return Oracle{
		handler: OracleHandler{
			version: 12,
		},
	}
}

func NewOracle11() DatabaseApi {
	return Oracle{
		handler: OracleHandler{
			version: 11,
		},
	}
}

func (handler OracleHandler) Config() Configuration {
	switch handler.version {
	case 11:
		return Configuration{
			DriverClass: "goracle",
			Host:        "localhost",
			Instance:    "xe",
			Password:    "HelloApes66",
			Port:        1521,
			Username:    "system",
		}
	case 12:
		return Configuration{
			DriverClass: "goracle",
			Host:        "localhost",
			Instance:    "ORCLPDB1",
			Password:    "HelloApes66",
			Port:        1522,
			Username:    "system",
		}
	default:
		log.Fatalf("unknown oracle version %d ", handler.version)
		return Configuration{}
	}
}

func (handler OracleHandler) ConnectionUrl() string {
	return fmt.Sprintf(
		"%s/%s@%s:%d/%s",
		handler.Config().Username,
		handler.Config().Password,
		handler.Config().Host,
		handler.Config().Port,
		handler.Config().Instance)
}

func (handler OracleHandler) Connect() (*sqlx.DB, Message, error) {
	var message Message
	con, err := sqlx.Connect(handler.Config().DriverClass, handler.ConnectionUrl())
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

func (handler OracleHandler) Execute(command string) (Message, error) {
	var message Message
	con, message, err := handler.Connect()
	if err != nil {
		return message, err
	}
	defer con.Close()
	_, err = con.Exec(command)
	return message, err
}

func (handler OracleHandler) ExecuteBatch(command []string) (Message, error) {
	var message Message
	con, message, err := handler.Connect()
	if err != nil {
		return message, err
	}
	defer con.Close()
	tx := con.MustBegin()
	for _, sql := range command {
		tx.MustExec(sql)
	}
	tx.Commit()
	return message, err
}

func (db Oracle) KillSession(username string) error {
	var sid sql.NullString
	var serial sql.NullString
	con, _, err := db.handler.Connect()
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
		_, err = db.handler.Execute(killSession)
		if err != nil {
			log.Fatal(err)
			return err
		}
	}
	return nil
}

func (db Oracle) DropUser(username string) ([]Message, error) {
	var messages = []Message{}
	var message Message
	err := db.KillSession(username)
	dropUserSql := fmt.Sprintf("drop user %s cascade", username)
	_, err = db.handler.Execute(dropUserSql)
	if err != nil {
		if strings.Contains(err.Error(), "ORA-01918") {
			message = Message{
				Severity: Warn,
				Content:  fmt.Sprintf("user %s does not exists: %s", username, err.Error()),
			}
			messages = append(messages, message)
			return messages, nil
		}
		message = Message{
			Severity: Error,
			Content:  err.Error(),
		}
		messages = append(messages, message)
		return messages, err
	}
	message = Message{
		Severity: Success,
		Content:  "user dropped " + username,
	}
	messages = append(messages, message)
	return messages, nil
}

// func (db Oracle) CreateUser(username string, password string) ([]Message, error) {
// 	messages := make([]Message, 0)
// 	var message Message
// 	createUserSql := fmt.Sprintf("create user %s identified by %s", username, password)
// 	_, err := db.handler.Execute(createUserSql)
// 	if err != nil {
// 		// user already exists
// 		if strings.Contains(err.Error(), "ORA-01920") {
// 			message = Message{
// 				Severity: Warn,
// 				Content:  fmt.Sprintf("user %s already exists: %s", username, err.Error()),
// 			}
// 			messages = append(messages, message)
// 			return messages, nil
// 		} else {
// 			message = Message{
// 				Severity: Error,
// 				Content:  err.Error(),
// 			}
// 			messages = append(messages, message)
// 		}
// 		log.Print(err)
// 		return messages, err
// 	}
// 	sqls := []string{
// 		fmt.Sprintf("grant all privileges to %s", username),
// 		fmt.Sprintf("grant all privileges to %s", username),
// 	}
// 	_, err = db.handler.ExecuteBatch(sqls)
// 	if err != nil {
// 		message := Message{
// 			Severity: Error,
// 			Content:  err.Error(),
// 		}
// 		messages = append(messages, message)
// 		return messages, err
// 	}
// 	message = Message{
// 		Severity: Success,
// 		Content:  "user created " + username,
// 	}
// 	messages = append(messages, message)
// 	return messages, nil
// }

func (db Oracle) CreateUser(username string, password string) ([]Message, error) {
	messages := make([]Message, 0)
	var message Message
	createUserSql := fmt.Sprintf("create user %s identified by %s", username, password)
	_, err := db.handler.Execute(createUserSql)
	if err != nil {
		// user already exists
		if strings.Contains(err.Error(), "ORA-01920") {
			message = Message{
				Severity: Warn,
				Content:  fmt.Sprintf("user %s already exists: %s", username, err.Error()),
			}
			messages = append(messages, message)
			return messages, nil
		} else {
			message = Message{
				Severity: Error,
				Content:  err.Error(),
			}
			messages = append(messages, message)
		}
		log.Print(err)
		return messages, err
	}
	_, err = db.handler.Execute(fmt.Sprintf("grant all privileges to %s", username))
	if err != nil {
		message := Message{
			Severity: Error,
			Content:  err.Error(),
		}
		messages = append(messages, message)
		return messages, err
	}
	message = Message{
		Severity: Success,
		Content:  "user created " + username,
	}
	messages = append(messages, message)
	return messages, nil
}

func (db Oracle) RecreateUser(username string, password string) ([]Message, error) {
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

func (db Oracle) ListUsers() ([]SystemUser, error) {
	con, _, err := db.handler.Connect()
	if err != nil {
		return nil, err
	}
	defer con.Close()
	oracleUsers := []OracleUser{}
	var sql string
	switch db.handler.version {
	case 11:
		sql = "SELECT username, user_id FROM dba_users WHERE user_id > 50 AND user_id < 1000000"
	case 12:
		sql = "SELECT username, user_id FROM dba_users WHERE ORACLE_MAINTAINED = 'N' ORDER BY username"
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
