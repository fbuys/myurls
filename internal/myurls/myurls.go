package myurls

import (
	_ "embed"
	"regexp"
	"strings"
)

//go:embed _myurls
var rawUrls string

type Url struct {
	// A string used to retrieve a url.
	Id string
	// A string that represents the address of the resource.
	Address string
}

func GetAllUrls() []Url {
	var result []Url
	linesPattern := regexp.MustCompile(`\n`)
	urlsPattern := regexp.MustCompile(`\s`)
	lines := linesPattern.Split(strings.TrimSpace(rawUrls), -1)
	for _, line := range lines {
		rawUrl := urlsPattern.Split(line, 2)
		result = append(result, Url{Id: rawUrl[0], Address: rawUrl[1]})
	}
	return result
}
