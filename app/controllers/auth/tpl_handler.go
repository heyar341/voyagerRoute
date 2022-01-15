package auth

import (
	"html/template"
)

var authTpl *template.Template

func init() {
	authTpl = template.Must(template.Must(template.ParseGlob("templates/auth/*.html")).ParseGlob("templates/includes/*.html"))
}
