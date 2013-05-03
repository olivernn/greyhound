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

type FileChan chan File

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

func printFiles (fileChannel chan File) {
  for {
    file, ok := <- fileChannel

    if !ok {
      return
    }

    fmt.Printf("%s\n", file)
  }
}

func walkDir (dir string, pattern string, exclude string, fileChannel chan File) {

  visit := func(path string, info os.FileInfo, err error) error {
    if !info.IsDir() {
      relPath, _ := filepath.Rel(dir, path)
      match, _ := regexp.MatchString(pattern, relPath)

      if match {
        file := &File{Name: filepath.Base(path), Path: path }
        fileChannel <- *file
      }
    } else {
      match, _ := regexp.MatchString(exclude, path)

      if len(exclude) > 0 && match {
        return filepath.SkipDir
      }
    }

    return nil
  }

  filepath.Walk(dir, visit)
  close(fileChannel)
}

func main () {
  flag.Parse()

  searchPattern := queryToSearchPattern(*query)
  searchDir := getSearchDir(*dir)
  excludePattern := excludeToExcludePattern(*exclude)

  fileChannel := make(FileChan)

  go walkDir (searchDir, searchPattern, excludePattern, fileChannel)

  printFiles(fileChannel)
}
