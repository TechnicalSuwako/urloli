package main

import (
  "text/template"
  "fmt"
  "strings"
  "net/http"
  "crypto/rand"
  "encoding/json"
  "unicode/utf8"
  "io/ioutil"
  "os"
  "runtime"
)

var (
  linkpath string
  configpath string
  payload map[string]interface{}
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

// http://かhttps://で始まるかどうか
func checkprefix (url string) bool {
  if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
    return false
  }

  return true
}

// URLは500文字以内かどうか
func checkcharlim (url string) bool {
  if utf8.RuneCountInString(url) > 500 {
    return false
  }
  return true
}

func geturl (url string, checkjson bool) string {
  payload := getlinks()

  for k := range payload {
    if checkjson {
      if url == payload[k] {
        return url
      }
    } else {
      if url == k {
        return payload[k].(string)
      }
    }
  }

  return ""
}

func insertjson (url string) string {
  payload := getlinks()

  newstring := mkstring()
  payload[newstring] = url
  m, _ := json.Marshal(&payload)
  payload = nil
  ioutil.WriteFile(linkpath, m, os.ModePerm)
  // fmt.Printf("%s\n", m)

  return newstring
}

type Page struct {
  Tit string
  Err string
  Url string
  Dom string
  Lan string
}

func getlinks () map[string]interface{} {
  data, err := ioutil.ReadFile(linkpath)
  if err != nil {
    fmt.Println("links.jsonを開けられません: ", err)
  }

  var payload map[string]interface{}
  json.Unmarshal(data, &payload)

  return payload
}

func main () {
  if runtime.GOOS == "freebsd" {
    linkpath = "/usr/local/etc/urloli/links.json"
    configpath = "/usr/local/etc/urloli/config.json"
  } else {
    linkpath = "/etc/urloli/links.json"
    configpath = "/etc/urloli/config.json"
  }
  var domain string
  payload := getlinks()

  for k := range payload {
    domain = payload[k].(string)
  }

  payload = nil
  http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    cookie, err := r.Cookie("lang")
    if err != nil {
      http.SetCookie(w, &http.Cookie {Name: "lang", Value: "ja", MaxAge: 31536000, Path: "/"})
      http.Redirect(w, r, "/", http.StatusSeeOther)
      return
    }

    uri := r.URL.Path
    query := r.URL.Query()
    qnewurl := query.Get("newurl")
    data := &Page{Tit: "トップ", Lan: cookie.Value}
    if cookie.Value == "en" {
      data = &Page{Tit: "Top", Lan: cookie.Value}
    }
    tmpl := template.Must(template.ParseFiles("view/index.html", "view/header.html", "view/footer.html"))

    if r.Method == "POST" {
      err := r.ParseForm()
      if err != nil { fmt.Println(err) }
      if r.PostForm.Get("sosin") != "" {
        if r.PostForm.Get("newadd") != "" {
          addurl := r.PostForm.Get("newadd")
          chkprx := checkprefix(addurl)
          chklim := checkcharlim(addurl)
          if !chkprx {
            if cookie.Value == "ja" {
                data = &Page{Tit: "不正なURL", Err: "URLは「http://」又は「https://」で始めます。", Lan: cookie.Value}
            } else {
              data = &Page{Tit: "Invalid URL", Err: "The URL should start with \"http://\" or \"https://\".", Lan: cookie.Value}
            }
            tmpl = template.Must(template.ParseFiles("view/404.html", "view/header.html", "view/footer.html"))
          }
          if !chklim {
            if cookie.Value == "ja" {
                data = &Page{Tit: "不正なURL", Err: "URLは500文字以内です。", Lan: cookie.Value}
            } else {
              data = &Page{Tit: "Invalid URL", Err: "The URL should be less than 500 characters.", Lan: cookie.Value}
            }
            data = &Page{Tit: "不正なURL", Err: ""}
            tmpl = template.Must(template.ParseFiles("view/404.html", "view/header.html", "view/footer.html"))
          }

          if chklim && chkprx {
            chkfn := geturl(addurl, true)
            if chkfn != "" {
              http.Redirect(w, r, addurl, http.StatusSeeOther)
              return
            } else {
              res := insertjson(addurl)
              if cookie.Value == "ja" {
                data = &Page{Tit: "短縮済み", Lan: cookie.Value, Url: res, Dom: domain}
              } else {
                data = &Page{Tit: "Shortened", Lan: cookie.Value, Url: res, Dom: domain}
              }
              tmpl = template.Must(template.ParseFiles("view/submitted.html", "view/header.html", "view/footer.html"))
            }
          }
        } else {
          if cookie.Value == "ja" {
              data = &Page{Tit: "未検出", Err: "URLをご入力下さい。", Lan: cookie.Value}
          } else {
            data = &Page{Tit: "Not found", Err: "Please enter a URL.", Lan: cookie.Value}
          }
          tmpl = template.Must(template.ParseFiles("view/404.html", "view/header.html", "view/footer.html"))
        }
      } else if r.PostForm.Get("langchange") != "" {
        if cookie.Value == "ja" {
          http.SetCookie(w, &http.Cookie {Name: "lang", Value: "en"})
        } else {
          http.SetCookie(w, &http.Cookie {Name: "lang", Value: "ja"})
        }
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
      }
    } else {
      if uri == "/" && qnewurl == "" {
        tmpl = template.Must(template.ParseFiles("view/index.html", "view/header.html", "view/footer.html"))
      } else if uri != "/" && qnewurl == "" {
        red := geturl(uri[1:], false)
        if red != "" {
          http.Redirect(w, r, red, http.StatusSeeOther)
          return
        } else {
          if cookie.Value == "ja" {
              data = &Page{Tit: "未検出", Err: "このURLを見つけられませんでした。", Lan: cookie.Value}
          } else {
            data = &Page{Tit: "Not found", Err: "This URL could not be found.", Lan: cookie.Value}
          }
          tmpl = template.Must(template.ParseFiles("view/404.html", "view/header.html", "view/footer.html"))
        }
      } else if uri == "/" && qnewurl != "" {
        data = &Page{Tit: "短縮済み", Url: qnewurl, Dom: domain}
        tmpl = template.Must(template.ParseFiles("view/submitted.html", "view/header.html", "view/footer.html"))
      } else {
        if cookie.Value == "ja" {
            data = &Page{Tit: "未検出", Err: "このURLを見つけられませんでした。", Lan: cookie.Value}
        } else {
          data = &Page{Tit: "Not found", Err: "This URL could not be found.", Lan: cookie.Value}
        }
        tmpl = template.Must(template.ParseFiles("view/404.html", "view/header.html", "view/footer.html"))
      }
    }

    tmpl.Execute(w, data)
    data = nil
  })

  fmt.Println("http://127.0.0.1:9910 でサーバーを実行中。終了するには、CTRL+Cを押して下さい。")
  http.ListenAndServe(":9910", nil)
}
