package database

import (
	"errors"

	_ "github.com/go-goracle/goracle"
	"github.com/jmoiron/sqlx"
)

type SystemUser struct {
	Username string `json:"username"`
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
}

type Response struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Messages []Message
}

type Severity string

const (
	Success Severity = "success"
	Info    Severity = "info"
	Warn    Severity = "warning"
	Error   Severity = "danger"
)

type Message struct {
	Severity Severity `json:"severity"`
	Content  string   `json:"content"`
}

type Configuration struct {
	Host          string
	Username      string
	Password      string
	Port          int
	Instance      string
	DriverClass   string
	ConnectionUrl string
}

type DatabaseHandler interface {
	Config() Configuration
	Connect() (*sqlx.DB, Message, error)
	Execute(command string) (Message, error)
	ConnectionUrl() string
}

type DatabaseApi interface {
	ListUsers() ([]SystemUser, error)
	CreateUser(username string, password string) ([]Message, error)
	DropUser(username string) ([]Message, error)
	RecreateUser(username string, password string) ([]Message, error)
}

func GetDatabaseHandler(db string) (DatabaseApi, error) {
	switch db {
	case "oracle12":
		return NewOracle12(), nil
	case "oracle11":
		return NewOracle11(), nil
	default:
		return nil, errors.New("unsupported database type")
	}
}
