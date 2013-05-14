package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

func main() {
	flag.Parse()

	if *help {
		printUsage()
		os.Exit(0)
	}

	searchDir := getSearchDir(*dir)
	excludePattern := ExcludeMatcher(*exclude)

	fileChannel := make(FileChan)
	scoredFileChannel := make(FileChan)

	go walkDir(searchDir, excludePattern, fileChannel)

	go scoreFiles(*query, fileChannel, scoredFileChannel)
	filterFiles(scoredFileChannel)
}

var help = flag.Bool("help", false, "Shows this message.")
var query = flag.String("query", "", "File path pattern to search for.")
var dir = flag.String("dir", "", "Directory to search in.")
var exclude = flag.String("exclude", "", "Sub directories to exclude from search.")
var limit = flag.Int("limit", 10, "Limit the number of files returned.")

type FileChan chan File

func getSearchDir(dir string) string {
	if len(dir) == 0 {
		wd, _ := os.Getwd()
		return wd
	}

	return dir
}

func filterFiles(fileChannel chan File) {
	files := make(Files, 0, *limit)

	for {
		file, ok := <-fileChannel

		if !ok {
			files.PrintPath()
			break
		}

		files.Add(file)
	}
}

func scoreFiles(query string, fileChannel chan File, scoredFileChannel chan File) {
	pattern := QueryMatcher(query)
	exactMatcher := ExactQueryMatcher(query)

	for {
		file, ok := <-fileChannel

		if !ok {
			close(scoredFileChannel)
			break
		}

		exactMatch := exactMatcher.MatchString(file.Name)

		if exactMatch {
			file.Score = len(file.Name)
			scoredFileChannel <- file
			continue
		}

		nameMatch := pattern.MatchString(file.Name)

		if nameMatch {
			file.Score = 10 * len(file.Path)
			scoredFileChannel <- file
			continue
		}

		pathMatch := pattern.MatchString(file.Path)

		if pathMatch {
			file.Score = 100 * len(file.Path)
			scoredFileChannel <- file
			continue
		}
	}
}

func walkDir(dir string, exclude *regexp.Regexp, fileChannel chan File) {

	hasExclude := len(exclude.String()) > 0

	visit := func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			relPath, _ := filepath.Rel(dir, path)
			file := &File{Name: filepath.Base(path), Path: relPath}
			fileChannel <- *file

		} else {
			match := exclude.MatchString(path)

			if hasExclude && match {
				return filepath.SkipDir
			}
		}

		return nil
	}

	filepath.Walk(dir, visit)
	close(fileChannel)
}

func printUsage() {
	fmt.Println("greyhound: Fast fuzzy filepath finder.\n")
	flag.PrintDefaults()
	fmt.Println("\n e.g. greyhound --query test --exclude .git,.svn --dir ~/code/project")
}
