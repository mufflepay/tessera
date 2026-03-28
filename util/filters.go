package util

import (
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func FilterName(name string) string {
	return strings.TrimSpace(cases.Title(language.English).String(name))
}

func FilterEmail(email string) string {
	return strings.TrimSpace(strings.ToLower(email))
}

func FilterPassword(password string) string {
	return strings.TrimSpace(password)
}
