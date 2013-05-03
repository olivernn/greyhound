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
var exclude = flag.String("exclude", "", "Sub directories to exclude from search")

type File struct {
  Name string
  Path string
}

func queryToSearchPattern(query string) string {
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

func excludeToExcludePattern (exclude string) string {
  pattern := strings.Replace(exclude, ",", "|", -1)
  return pattern
}

func main () {
  flag.Parse()

  searchPattern := queryToSearchPattern(*query)
  searchDir := getSearchDir(*dir)
  excludePattern := excludeToExcludePattern(*exclude)

  visit := func(path string, info os.FileInfo, err error) error {
    if !info.IsDir() {
      relPath, _ := filepath.Rel(searchDir, path)
      match, _ := regexp.MatchString(searchPattern, relPath)

      if match {
        file := File{Name: filepath.Base(path), Path: path }
        fmt.Printf("%s\n", file)
      }
    } else {
      match, _ := regexp.MatchString(excludePattern, path)

      if len(excludePattern) > 0 && match {
        return filepath.SkipDir
      }
    }

    return nil
  }

  filepath.Walk(searchDir, visit)
}
