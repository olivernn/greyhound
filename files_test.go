package main

import (
  "testing"
  "fmt"
)

func TestAdd(t *testing.T) {
  file := File{Name: "foo.go", Path: "src/foo.go", Score: 10}
  files := make(Files, 0, 1)

  files.Add(file)

  if files.Len() != 1 {
    fmt.Println(files.Len())
    t.FailNow()
  }

  if files.Peek() != file {
    t.FailNow()
  }

  anotherFile := File{Name: "bar.go", Path: "src/bar.go", Score: 1}

  files.Add(anotherFile)

  if files.Len() != 1 {
    fmt.Println(files.Len())
    t.FailNow()
  }

  if files.Peek() != anotherFile {
    fmt.Println(files.Peek())
    t.FailNow()
  }
}
