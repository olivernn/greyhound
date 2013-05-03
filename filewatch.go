package main

import (
  "os"
  "path/filepath"
  "log"
  "fmt"
  "github.com/howeyc/fsnotify"
)

func main () {
  watcher, err := fsnotify.NewWatcher()

  if err != nil {
    log.Fatal(err)
  }

  dir, err := os.Getwd()

  if err != nil {
    log.Fatal(err)
  }

  go watcher.Watch(dir)

  visit := func(path string, info os.FileInfo, err error) error {
    if info.IsDir() {
      fmt.Println("dir: ", path)
    } else {
      fmt.Println("file: ", path)
    }

    return nil
  }

  err = filepath.Walk(dir, visit)

  if err != nil {
    log.Fatal(err)
  }

  for {
    select {
    case ev := <- watcher.Event:
      fmt.Println("event: ", ev)
    case err := <- watcher.Error:
      fmt.Println("error: ", err)
    }
  }

}
