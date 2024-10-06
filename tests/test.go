package main

import (
	"fmt"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func main() {
	original := "hello, world!"

	caser := cases.Title(language.English)
	capitalized := caser.String(original)

	fmt.Println(capitalized) // Output: Hello, World!
}
