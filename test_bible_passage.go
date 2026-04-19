package main

import (
	"fmt"
	"regexp"
)

func main() {
	verseReferenceRegex := regexp.MustCompile(` +(?s)^([\w\s]+\d+:\d+(?:-\d+)?)\s+\(([^)]+)\)\s+(.*)$`)
	fmt.Println(verseReferenceRegex.MatchString("Romans 3 (ESV) ..."))
}
