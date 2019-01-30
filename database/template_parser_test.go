package database

import (
	"flag"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func SetTestFlag()  {
	flag.String("templates", "../config", "directory for *.tpl files contains sql scripts")
}

func TestLoadTemplate(t *testing.T) {
	SetTestFlag()
	user := TemplateValue{User: "user1", Password: "pass1"}
	res, err := LoadTemplate(user, TemplateSqlServerCreate)
	if err != nil {
		t.Error(err)
	}
	log.Printf("%v", res)
	assert.Contains(t, res, user.User, "expect to find the user name in template")
}