package main

import (
  "container/heap"
  "fmt"
)

type Files []File

func (f Files) Len() int {
  return len(f)
}

func (f Files) Less(i, j int) bool {
  return f[i].Score < f[j].Score
}

func (f Files) Swap(i, j int) {
  f[i], f[j] = f[j], f[i]
}

func (f Files) Peek() File {
  return f[0]
}

func (f *Files) PrintPath() {
  numFiles := len(*f)

  for i := 0; i < numFiles; i++ {
    file := heap.Pop(f).(File)
    fmt.Println(file.Path)
  }
}

func (f *Files) Add(file File) {
  if len(*f) < cap(*f) {
    heap.Push(f, file)
  } else {
    if f.Peek().Score > file.Score {
      heap.Pop(f)
      heap.Push(f, file)
    }
  }
}

func (f *Files) Push(x interface{}) {
  *f = append(*f, x.(File))
}

func (f *Files) Pop() interface{} {
  oldf := *f
  x := oldf[len(oldf) - 1]
  newf := oldf[0 : len(oldf) - 1]
  *f = newf
  return x
}

