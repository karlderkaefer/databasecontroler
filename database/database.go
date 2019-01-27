package database

import "github.com/jmoiron/sqlx"

const (
	Success Severity = "success"
	Info    Severity = "info"
	Warn    Severity = "warning"
	Error   Severity = "danger"
)

type Severity string

type Message struct {
	Severity Severity `json:"severity"`
	Content  string   `json:"content"`
}

type SystemUser struct {
	Username string `json:"username"`
}

type Response struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Messages []Message
}

type Configuration struct {
	Host          string
	Username      string
	Password      string
	Port          int
	Instance      string
	DriverClass   string
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


