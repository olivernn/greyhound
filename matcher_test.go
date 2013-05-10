package main

import "testing"

func TestQueryMatcher(t *testing.T) {
  query := "moprj"
  pattern := QueryMatcher(query)

  if !pattern.MatchString("app/models/project.rb") {
    t.FailNow()
  }

  if !pattern.MatchString("spec/Models/Project.rb") {
    t.FailNow()
  }

  if !pattern.MatchString("spec/models/project_spec.rb") {
    t.FailNow()
  }

  if pattern.MatchString("views/project/_project.html.erb") {
    t.FailNow()
  }
}

func TestExactQueryMatcher(t *testing.T) {
  query := "project.rb"
  pattern := ExactQueryMatcher(query)

  if !pattern.MatchString("Project.rb") {
    t.FailNow()
  }

  if !pattern.MatchString("project.rb") {
    t.FailNow()
  }

  if pattern.MatchString("project_spec.rb") {
    t.FailNow()
  }
}

func TestExcludeMatcher (t *testing.T) {
  query := ".git,.svn"
  pattern := ExcludeMatcher(query)

  if !pattern.MatchString(".git") {
    t.FailNow()
  }

  if !pattern.MatchString(".svn") {
    t.FailNow()
  }

  if pattern.MatchString(".GIT") {
    t.FailNow()
  }

  if pattern.MatchString("gittit") {
    t.FailNow()
  }
}
