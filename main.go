package main

import (
  "fmt"
  "os"
  "strconv"
)

var version = "2.0.2"

func help () {
  fmt.Println("使い方：");
  fmt.Println("urloli -v               ：バージョンを表示");
  fmt.Println("urloli -s [ポート番号]  ：ポート番号でウェブサーバーを実行（デフォルト＝9910）");
  fmt.Println("urloli -h               ：ヘルプを表示");
  fmt.Println("urloli <URL>            ：コマンドラインでURLを短縮");
}

func main () {
  cnf, err := getconf()
  if err != nil {
    fmt.Println(err)
    return
  }
  args := os.Args

  if len(args) == 2 {
    if args[1] == "-v" {
      fmt.Println("urloli-" + version)
      return
    } else if args[1] == "-s" {
      serv(cnf, 9910)
    } else if args[1] == "-h" {
      help()
      return
    } else {
      if checkprefix(args[1]) {
        _, key := geturl(args[1], cnf.linkpath, true)
        if (key != "") {
          fmt.Println(cnf.domain + "/" + key)
        } else {
          fmt.Println(cnf.domain + "/" + insertjson(args[1], cnf.linkpath))
        }
        return
      } else {
        fmt.Println("URLは不正です。終了…")
        return
      }
    }
  } else if len(args) == 3 && args[1] == "-s" {
    if port, err := strconv.Atoi(args[2]); err != nil {
      fmt.Printf("%qは数字ではありません。\n", args[2])
      return
    } else {
      serv(cnf, port)
    }
  } else {
    help()
    return
  }
}
