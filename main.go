package main

import (
	"fmt"
	"github.com/russross/blackfriday"
)

func main() {
	s := blackfriday.Run([]byte("# This is magnanimous!\n```go\nfunc main () { }\n```\n"))
	fmt.Println(string(s))
}
