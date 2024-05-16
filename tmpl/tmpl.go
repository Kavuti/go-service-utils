package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"text/template"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func main() {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic(errors.New("unable to get the current filename"))
	}
	directory := filepath.Dir(filename)

	projectDirectory, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	entityName := flag.String("name", "Foo", "The name of the entity to create")
	projectName := flag.String("project-name", filepath.Base(projectDirectory), "The project name")
	flag.Parse()

	capitalizedEntityName := cases.Title(language.English, cases.Compact).String(*entityName)

	// Protobuf template
	renderTemplate(directory, projectDirectory+"/proto", capitalizedEntityName, *projectName, "proto.tmpl", strings.ToLower(*entityName), "proto")

	// Service template
	renderTemplate(directory, projectDirectory+"/service", capitalizedEntityName, *projectName, "service.tmpl", strings.ToLower(*entityName), "go")

	// SQLC queries template
	renderTemplate(directory, projectDirectory+"/db/queries", capitalizedEntityName, *projectName, "sqlc_queries.tmpl", strings.ToLower(*entityName), "sql")

	// Migration
	r, _ := regexp.Compile("[0-9]{4}")

	maxNumber := 0
	filepath.Walk(projectDirectory+"/db/migrations", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if r.MatchString(info.Name()) {
			number, _ := strconv.Atoi(r.FindString(info.Name()))
			if number > maxNumber {
				maxNumber = number
			}
		}
		return nil
	})
	renderTemplate(directory, projectDirectory+"/db/migrations", capitalizedEntityName, *projectName, "migration.tmpl", fmt.Sprintf("%04d_%s", maxNumber+1, strings.ToLower(*entityName)), "sql")
}

func renderTemplate(directory string, destinationDirectory string, entityName string, projectName string, templateName string, filename string, fileext string) {
	var tmplData struct {
		EntityName  string
		ProjectName string
	}
	tmplData.EntityName = entityName
	tmplData.ProjectName = projectName

	funcMap := template.FuncMap{
		"ToLower": strings.ToLower,
	}

	tmpl, err := template.New(templateName).Funcs(funcMap).ParseFiles(directory + "/" + templateName)
	if err != nil {
		panic(err)
	}
	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf, tmplData)
	if err != nil {
		panic(err)
	}
	os.WriteFile(destinationDirectory+"/"+filename+"."+fileext, buf.Bytes(), 0644)
}
