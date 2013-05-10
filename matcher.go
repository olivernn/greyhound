package main

import (
  "regexp"
  "strings"
)

func QueryMatcher(query string) *regexp.Regexp {
  tokens := strings.Split(query, "")

  for idx, token := range tokens {
    tokens[idx] = regexp.QuoteMeta(token)
  }

  pattern := strings.Join(tokens, ".*")
  return regexp.MustCompile("(?i)" + pattern)
}

func ExactQueryMatcher(query string) *regexp.Regexp {
  escaped := regexp.QuoteMeta(query)
  return regexp.MustCompile("(?i)" + escaped)
}

func ExcludeMatcher(query string) *regexp.Regexp {
  escaped := regexp.QuoteMeta(query)
  pattern := strings.Replace(escaped, ",", "|", -1)
  return regexp.MustCompile(pattern)
}
