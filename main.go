package main

import (
  "fmt"
  "os"
  "strconv"
)

var sofname = "urloli"
var version = "2.3.0"

func usage() {
  fmt.Printf("%s-%s\nusage: %s [-s port] [url]\n", sofname, version, sofname)
}

func main() {
  cnf, err := getconf()
  if err != nil {
    fmt.Println(err)
    return
  }
  args := os.Args

  if len(args) < 2 {
    usage()
    return
  }

  if len(args) == 2 && args[1] == "-s" {
    serv(cnf, 9910)
  } else if len(args) == 2 && args[1] != "-s" {
    if !checkprefix(args[1]) {
      fmt.Println("URLは不正です。終了…")
      return
    }

    _, key := geturl(args[1], cnf.linkpath, true)
    if (key != "") {
      fmt.Println(cnf.domain + "/" + key)
    } else {
      fmt.Println(cnf.domain + "/" + insertjson(args[1], cnf.linkpath))
    }
    return
  } else if len(args) == 3 && args[1] == "-s" {
    port, err := strconv.Atoi(args[2])
    if err != nil {
      fmt.Printf("%qは数字ではありません。\n", args[2])
      return
    }

    serv(cnf, port)
  }
}
