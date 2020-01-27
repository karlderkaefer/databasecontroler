package database

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"text/template"
)

type TemplateValue struct {
	User     string
	Password string
}

const (
	TemplateSqlServerCreate = "sqlserver-create.tpl"
	TemplateSqlServerDrop   = "sqlserver-drop.tpl"
)

func LoadTemplate(value TemplateValue, templateName string) (string, error) {
	path := flag.Lookup("templates").Value.(flag.Getter).Get().(string)
	templateFile := fmt.Sprintf("%s/%s", path, templateName)
	log.Printf("loading template from %s", templateFile)
	tmpl, err := template.ParseFiles(templateFile)
	if err != nil {
		return "", err
	}
	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, value)
	if err != nil {
		return "", err
	}
	return buffer.String(), nil
}
