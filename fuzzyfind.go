package main

import (
  "fmt"
  "flag"
  "os"
  "path/filepath"
  "regexp"
)

var query = flag.String("query", "", "The query")

func main () {
  flag.Parse()

  visit := func(path string, info os.FileInfo, err error) error {
    if !info.IsDir() {
      relPath, _ := filepath.Rel("/Users/olivernightingale/code/lunr.js", path)
      match, _ := regexp.MatchString(*query, relPath)

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
