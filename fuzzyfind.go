package main

import (
  "fmt"
  "flag"
)

var query = flag.String("query", "", "The query")

func main () {
  flag.Parse()

  fmt.Printf("%s", *query)
}
