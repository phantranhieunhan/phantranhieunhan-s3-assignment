package util

import (
	"regexp"
)

func GetEmailsFromString(str string) []string {
	re := regexp.MustCompile(`[\w\.\-]+@[\w\.\-]+\.\w+`)
	emails := re.FindAllString(str, -1)
	return emails
}
