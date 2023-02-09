package utils

import (
	// "fmt"
	"fmt"
	"log"
	"net/url"
	"strings"
)

func DecodeURLString(encodeValue string) string {
	decodeValue, err := url.QueryUnescape(encodeValue)

	if err != nil {
		log.Fatal(err)
		return ""
	}

	return decodeValue
}

func StandardizedDocument(document string) string {
	return strings.TrimSpace(document)
}

func StrimNewLine(document string) string {
  x := document
	x = strings.Replace(x, " ", "", -1)
	x = strings.Replace(x, "\t", "", -1)
	x = strings.Replace(x, "\n", "", -1)
  reader := fmt.Sprintln(x)
  reader = strings.Replace(x, "\n", "", -1)
	fmt.Println(strings.Replace("postgresql   \n database", "\n", "", -1))
  fmt.Println(reader)
  fmt.Println(strings.Count("postgresql   \n database", "\n"))
  fmt.Println(strings.Compare(document,"postgresql   \n database"))
	return reader
}
