package main

import (
  "fmt"
  "encoding/json"
  "io/ioutil"
  "runtime"
  "os"
  "errors"
)

type Config struct {
  configpath, linkpath, webpath, domain, ip string
}

func geturl (url string, linkpath string, checkjson bool) (string, string) {
  payload := getlinks(linkpath)

  for k := range payload {
    if checkjson {
      if url == payload[k] {
        return url, k
      }
    } else {
      if url == k {
        return payload[k].(string), k
      }
    }
  }

  return "", ""
}

func getlinks (linkpath string) map[string]interface{} {
  data, err := ioutil.ReadFile(linkpath)
  if err != nil {
    fmt.Println("links.jsonを開けられません: ", err)
  }

  var payload map[string]interface{}
  json.Unmarshal(data, &payload)

  return payload
}

func getconf () (Config, error) {
  var cnf Config

  prefix := "/usr"
  if runtime.GOOS == "freebsd" || runtime.GOOS == "openbsd" {
    prefix += "/local"
  } else if runtime.GOOS == "netbsd" {
    prefix += "/pkg"
  }

  cnf.configpath = "/etc/urloli/config.json"
  cnf.linkpath = "/etc/urloli/links.json"
  if runtime.GOOS == "freebsd" || runtime.GOOS == "netbsd" {
    cnf.configpath = prefix + cnf.configpath
    cnf.linkpath = prefix + cnf.linkpath
  }

  data, err := ioutil.ReadFile(cnf.configpath)
  if err != nil {
    fmt.Println("config.jsonを開けられません: ", err)
    return cnf, errors.New("コンフィグファイルは " + cnf.configpath + " に創作して下さい。")
  }

  var payload map[string]interface{}
  json.Unmarshal(data, &payload)
  if payload["webpath"] == nil {
    return cnf, errors.New("「webpath」の値が設置していません。")
  }
  if payload["domain"] == nil {
    return cnf, errors.New("「domain」の値が設置していません。")
  }
  if payload["ip"] == nil {
    return cnf, errors.New("「ip」の値が設置していません。")
  }
  if _, err := os.Stat(payload["webpath"].(string)); err != nil {
    fmt.Printf("%v\n", err)
    return cnf, errors.New("mkdirコマンドを使って、 " + payload["webpath"].(string))
  }
  if !checkprefix(payload["domain"].(string)) {
    return cnf, errors.New("URLは「http://」又は「https://」で始める様にして下さい。")
  }
  cnf.webpath = payload["webpath"].(string)
  cnf.domain = payload["domain"].(string)
  cnf.ip = payload["ip"].(string)
  payload = nil

  return cnf, nil
}
