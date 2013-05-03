package main

import (
  "fmt"
  "flag"
  "os"
  "path/filepath"
  "regexp"
  "strings"
)

var query = flag.String("query", "", "The query")

func stringToPattern(query string) string {
  tokens := strings.Split(query, "")
  pattern := strings.Join(tokens, ".*")
  pattern = ".*" + pattern + ".*"
  return pattern
}

func main () {
  flag.Parse()

  pattern := stringToPattern(*query)

  visit := func(path string, info os.FileInfo, err error) error {
    if !info.IsDir() {
      relPath, _ := filepath.Rel("/Users/olivernightingale/code/lunr.js", path)
      match, _ := regexp.MatchString(pattern, relPath)

      if match {
        fmt.Println(relPath)
      }
    } else {
      match, _ := regexp.MatchString(".git", path)

      if match {
        return filepath.SkipDir
      }
    }

    return nil
  }

  filepath.Walk("/Users/olivernightingale/code/lunr.js", visit)
}
