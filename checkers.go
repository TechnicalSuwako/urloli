package main

import (
  "strings"
  "unicode/utf8"
)

// http://かhttps://で始まるかどうか
func checkprefix (url string) bool {
  return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
}

// URLは500文字以内かどうか
func checkcharlim (url string) bool {
  return utf8.RuneCountInString(url) <= 500
}
