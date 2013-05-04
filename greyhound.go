package main

import (
  "fmt"
  "flag"
  "os"
  "path/filepath"
  "regexp"
  "strings"
  "sort"
  "sync"
)

var help = flag.Bool("help", false, "Shows this message.")
var query = flag.String("query", "", "File path pattern to search for.")
var dir = flag.String("dir", "", "Directory to search in.")
var exclude = flag.String("exclude", "", "Sub directories to exclude from search.")

type File struct {
  Name string
  Path string
}

type FileChan chan File

type Files []File

func (f Files) Len () int {
  return len(f)
}

func (f Files) Swap (i, j int) {
  f[i], f[j] = f[j], f[i]
}

func printFiles (files []File) {
  for _, file := range files {
    fmt.Println(file.Path)
  }
}

type ByPathLength struct{ Files }

func (s ByPathLength) Less (i, j int) bool {
  return len(s.Files[i].Path) < len(s.Files[j].Path)
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

  var wg sync.WaitGroup

  wg.Add(3)

  go func () {
    sort.Sort(ByPathLength{exactMatches})
    wg.Done()
  }()

  go func () {
    sort.Sort(ByPathLength{nameMatches})
    wg.Done()
  }()

  go func () {
    sort.Sort(ByPathLength{pathMatches})
    wg.Done()
  }()

  wg.Wait()

  printFiles(exactMatches)
  printFiles(nameMatches)
  printFiles(pathMatches)
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

func printUsage () {
  fmt.Println("greyhound: Fast fuzzy filepath finder.\n")
  flag.PrintDefaults()
  fmt.Println("\n e.g. greyhound --query test --exclude .git,.svn --dir ~/code/project")
}

func main () {
  flag.Parse()

  if *help {
    printUsage()
    os.Exit(0)
  }

  searchDir := getSearchDir(*dir)
  excludePattern := excludeToExcludePattern(*exclude)

  fileChannel := make(FileChan)

  go walkDir (searchDir, excludePattern, fileChannel)

  filterFiles(*query, fileChannel)
}
