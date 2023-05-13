package main

import (
  "encoding/json"
  "fmt"
)

func getlist (lang string) []byte {
  var jloc = []byte(`{
    "top": "トップ",
    "fuseiurl": "不正なURL",
    "tansyukuzumi": "短縮済み",
    "mikensyutu": "未検出",
    "errfusei": "URLは「http://」又は「https://」で始めます。",
    "errcharlim": "URLは500文字以内です。",
    "errurlent": "URLをご入力下さい。",
    "errurlnai": "このURLを見つけられませんでした。"
  }`)
  var eloc = []byte(`{
    "top": "Top",
    "fuseiurl": "Invalid URL",
    "tansyukuzumi": "Shortened",
    "mikensyutu": "Not found",
    "errfusei": "The URL should start with \"http://\" or \"https://\".",
    "errcharlim": "The URL should be less than 500 characters.",
    "errurlent": "Please enter a URL.",
    "errurlnai": "This URL could not be found."
  }`)

  if lang == "en" { return eloc }
  return jloc
}

func getloc (str string, lang string) string {
  var payload map[string]interface{}
  err := json.Unmarshal(getlist(lang), &payload)
  if err != nil {
    fmt.Println("loc:", err)
    return ""
  }

  for k, v := range payload {
    if str == k {
      return v.(string)
    }
  }

  return ""
}
