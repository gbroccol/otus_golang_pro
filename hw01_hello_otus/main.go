package main

import (
	"fmt"

	//nolint:depguard
	"github.com/agrison/go-commons-lang/stringUtils"
)

func main() {
	fmt.Println(stringUtils.Reverse("Hello, OTUS!"))
}
