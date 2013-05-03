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
var dir = flag.String("dir", "", "Directory to search in")

func stringToPattern(query string) string {
  tokens := strings.Split(query, "")
  pattern := strings.Join(tokens, ".*")
  pattern = ".*" + pattern + ".*"
  return pattern
}

func getSearchDir (dir string) string {
  if len(dir) == 0 {
    wd, _ := os.Getwd()
    return wd
  }

  return dir
}

func main () {
  flag.Parse()

  pattern := stringToPattern(*query)

  searchDir := getSearchDir(*dir)

  visit := func(path string, info os.FileInfo, err error) error {
    if !info.IsDir() {
      relPath, _ := filepath.Rel(searchDir, path)
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

  filepath.Walk(searchDir, visit)
}
