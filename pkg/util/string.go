package util

import (
	"regexp"
)

func GetEmailsFromString(str string) []string {
	re := regexp.MustCompile(`[\w\.\-]+@[\w\.\-]+\.\w+`)
	emails := re.FindAllString(str, -1)
	return emails
}

func RemoveDuplicates(s []string) []string {
    encountered := map[string]bool{}
    result := []string{}

    for _, v := range s {
        if !encountered[v] {
            encountered[v] = true
            result = append(result, v)
        }
    }
    return result
}