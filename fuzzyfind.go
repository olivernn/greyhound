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

func filterFiles (query string, fileChannel chan File) {
  exactMatches := make([]File, 0)
  nameMatches := make([]File, 0)
  pathMatches := make([]File, 0)

  pattern := queryToSearchPattern(query)

  for {
    file, ok := <- fileChannel

    if !ok {
      break
    }

    exactMatch, _ := regexp.MatchString(query, file.Name)

    if exactMatch {
      exactMatches = append(exactMatches, file)
      continue
    }

    nameMatch, _ := regexp.MatchString(pattern, file.Name)

    if nameMatch {
      nameMatches = append(nameMatches, file)
      continue
    }

    pathMatch, _ := regexp.MatchString(pattern, file.Path)

    if pathMatch {
      pathMatches = append(pathMatches, file)
      continue
    }
  }

  for _, file := range exactMatches {
    fmt.Println(file.Path)
  }

  for _, file := range nameMatches {
    fmt.Println(file.Path)
  }

  for _, file := range pathMatches {
    fmt.Println(file.Path)
  }
}

func walkDir (dir string, exclude string, fileChannel chan File) {

  visit := func(path string, info os.FileInfo, err error) error {
    if !info.IsDir() {
      relPath, _ := filepath.Rel(dir, path)
      file := &File{Name: filepath.Base(path), Path: relPath }
      fileChannel <- *file

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

  searchDir := getSearchDir(*dir)
  excludePattern := excludeToExcludePattern(*exclude)

  fileChannel := make(FileChan)

  go walkDir (searchDir, excludePattern, fileChannel)

  filterFiles(*query, fileChannel)
}
