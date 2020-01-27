package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

const (
	// Success message code for everything went fine
	Success Severity = "success"
	// Info message code if command went fine but there were expected errors
	Info Severity = "info"
	// Warn message code if command went fine but there were unexpected errors
	Warn Severity = "warning"
	// Error message code if command went wrong
	Error Severity = "danger"
	// UserCreated is a message that indicates that a database user has been create sucessfully
	UserCreated string = "user %s created"
	// UserAlreadyExists is a message that indicates that the requested database user already exists
	UserAlreadyExists string = "user %s already exists: %s"
	// UserNotExists is a message that indicates that a database user could not be find in database
	UserNotExists string = "user %s does not exists: %s"
	// UserDropped is a message that indicates that a database user has been successfully been dropped
	UserDropped string = "user %s has been dropped"
	// NameMaxLength is a message that will show the restriction for length for database names
	NameMaxLength string = "database name length has to be equal or less than %d"
)

// Severity represents the status code of a command
type Severity string

// Message struct will be used to structure the output so that frontend can parse it
type Message struct {
	Severity Severity `json:"severity"`
	Content  string   `json:"content"`
}

// SystemUser json struct represents a database user
type SystemUser struct {
	Username string `json:"username"`
}

type configuration struct {
	Host        string
	Username    string
	Password    string
	Port        int
	Instance    string
	DriverClass string
}

// Database is an interface that hold information of the configuration of the database and the API interface
type Database interface {
	databaseConfig
	DbAPI
}

type databaseConfig interface {
	Config() configuration
	Connect() (*sqlx.DB, Message, error)
	Execute(command string) (Message, error)
	ConnectionURL() string
}

// DbAPI is the interface every database will implement
type DbAPI interface {
	// CreateUser generates a new database user
	CreateUser(username string, password string) ([]Message, error)
	// DropUser will drop the requested database user
	DropUser(username string) ([]Message, error)
	// RecreateUser will sequentially call DropUser and CreateUser
	RecreateUser(username string, password string) ([]Message, error)
	// ListUsers will list all database excluding system internal users
	ListUsers() ([]SystemUser, error)
}

func addError(messages []Message, err error) ([]Message, error) {
	fmt.Println(err)
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

func addWarn(messages []Message, content string) []Message {
	message := Message{
		Severity: Warn,
		Content:  content,
	}
	messages = append(messages, message)
	return messages
}

func recreateUser(db DbAPI, username string, password string) ([]Message, error) {
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
