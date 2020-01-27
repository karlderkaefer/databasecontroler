package database

import (
	"flag"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

// initialize flags once per test suite
func TestMain(m *testing.M) {
	flag.String("templates", "../config", "directory for *.tpl files contains sql scripts")
	flag.Parse()
	os.Exit(m.Run())
}

func TestLoadTemplate(t *testing.T) {
	user := TemplateValue{User: "user1", Password: "pass1"}
	res, err := LoadTemplate(user, TemplateSqlServerCreate)
	if err != nil {
		t.Error(err)
	}
	log.Printf("%v", res)
	assert.Contains(t, res, user.User, "expect to find the user name in template")
}
