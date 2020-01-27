package database

import "github.com/jmoiron/sqlx"

const (
	Success           Severity = "success"
	Info              Severity = "info"
	Warn              Severity = "warning"
	Error             Severity = "danger"
	UserCreated       string   = "user %s created"
	UserAlreadyExists string   = "user %s already exists: %s"
	UserNotExists     string   = "user %s does not exists: %s"
	UserDropped       string   = "user %s has been dropped"
	NameMaxLength     string   = "database name length has to be equal or less than %d"
)

type Severity string

type Message struct {
	Severity Severity `json:"severity"`
	Content  string   `json:"content"`
}

type SystemUser struct {
	Username string `json:"username"`
}

type UserTemplate struct {
	User string
}

type Response struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Messages []Message
}

type Configuration struct {
	Host        string
	Username    string
	Password    string
	Port        int
	Instance    string
	DriverClass string
}

type Database interface {
	DatabaseConfig
	DatabaseApi
}

type DatabaseConfig interface {
	Config() Configuration
	Connect() (*sqlx.DB, Message, error)
	Execute(command string) (Message, error)
	ConnectionUrl() string
}

type DatabaseApi interface {
	CreateUser(username string, password string) ([]Message, error)
	DropUser(username string) ([]Message, error)
	RecreateUser(username string, password string) ([]Message, error)
	ListUsers() ([]SystemUser, error)
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

func addWarn(messages []Message, content string) []Message {
	message := Message{
		Severity: Warn,
		Content:  content,
	}
	messages = append(messages, message)
	return messages
}

func recreateUser(db DatabaseApi, username string, password string) ([]Message, error) {
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
