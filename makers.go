package main

import (
  "crypto/rand"
  "encoding/json"
  "io/ioutil"
  "os"
)

func mkstring () string {
  stringchars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
  newstring := ""
  b := make([]byte, 5)

  // 乱数を生成
  if _, err := rand.Read(b); err != nil {
    return "不明なエラー"
  }

  // ランダムに取り出して文字列を生成
  for _, v := range b {
    // index が stringchars の長さに収まるように調整
    newstring += string(stringchars[int(v)%len(stringchars)])
  }

  return newstring
}

func insertjson (url string, linkpath string) string {
  payload := getlinks(linkpath)

  newstring := mkstring()
  payload[newstring] = url
  m, _ := json.Marshal(&payload)
  payload = nil
  ioutil.WriteFile(linkpath, m, os.ModePerm)

  return newstring
}
