package main

import (
  "text/template"
  "fmt"
  "net/http"
)

type Page struct {
  Tit string
  Err string
  Url string
  Dom string
  Lan string
  Ver string
}

func serv (cnf Config, port int) {
  http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    data := &Page{Tit: "トップ", Ver: version}
    cookie, err := r.Cookie("lang")
    if err != nil {
      data.Lan = "ja"
    } else {
      data.Lan = cookie.Value
    }

    uri := r.URL.Path
    query := r.URL.Query()
    qnewurl := query.Get("newurl")
    if data.Lan == "en" {
      data.Tit = "Top"
    }
    tmpl := template.Must(template.ParseFiles(cnf.webpath + "/view/index.html", cnf.webpath + "/view/header.html", cnf.webpath + "/view/footer.html"))

    if r.Method == "POST" {
      err := r.ParseForm()
      if err != nil { fmt.Println(err) }
      if r.PostForm.Get("sosin") != "" {
        if r.PostForm.Get("newadd") != "" {
          addurl := r.PostForm.Get("newadd")
          chkprx := checkprefix(addurl)
          chklim := checkcharlim(addurl)
          if !chkprx {
            if data.Lan == "ja" {
              data.Tit = "不正なURL"
              data.Err = "URLは「http://」又は「https://」で始めます。"
            } else {
              data.Tit = "Invalid URL"
              data.Err = "The URL should start with \"http://\" or \"https://\"."
            }
            tmpl = template.Must(template.ParseFiles(cnf.webpath + "/view/404.html", cnf.webpath + "/view/header.html", cnf.webpath + "/view/footer.html"))
          }
          if !chklim {
            if data.Lan == "ja" {
              data.Tit = "不正なURL"
              data.Err = "URLは500文字以内です。"
            } else {
              data.Tit = "Invalid URL"
              data.Err = "The URL should be less than 500 characters."
            }
            tmpl = template.Must(template.ParseFiles(cnf.webpath + "/view/404.html", cnf.webpath + "/view/header.html", cnf.webpath + "/view/footer.html"))
          }

          if chklim && chkprx {
            chkfn, _ := geturl(addurl, cnf.linkpath, true)
            if chkfn != "" {
              http.Redirect(w, r, addurl, http.StatusSeeOther)
              return
            } else {
              res := insertjson(addurl, cnf.linkpath)
              data.Url = res
              data.Dom = cnf.domain
              if data.Lan == "ja" {
                data.Tit = "短縮済み"
              } else {
                data.Tit = "Shortened"
              }
              tmpl = template.Must(template.ParseFiles(cnf.webpath + "/view/submitted.html", cnf.webpath + "/view/header.html", cnf.webpath + "/view/footer.html"))
            }
          }
        } else {
          if data.Lan == "ja" {
            data.Tit = "未検出"
            data.Err = "URLをご入力下さい。"
          } else {
            data.Tit = "Not found"
            data.Err = "Please enter a URL."
          }
          tmpl = template.Must(template.ParseFiles(cnf.webpath + "/view/404.html", cnf.webpath + "/view/header.html", cnf.webpath + "/view/footer.html"))
        }
      } else if r.PostForm.Get("langchange") != "" {
        cookie, err := r.Cookie("lang")
        if err != nil || cookie.Value == "ja" {
          http.SetCookie(w, &http.Cookie {Name: "lang", Value: "en", MaxAge: 31536000, Path: "/"})
        } else {
          http.SetCookie(w, &http.Cookie {Name: "lang", Value: "ja", MaxAge: 31536000, Path: "/"})
        }
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
      }
    } else {
      if uri == "/" && qnewurl == "" {
        tmpl = template.Must(template.ParseFiles(cnf.webpath + "/view/index.html", cnf.webpath + "/view/header.html", cnf.webpath + "/view/footer.html"))
      } else if uri != "/" && qnewurl == "" {
        red, _ := geturl(uri[1:], cnf.linkpath, false)
        if red != "" {
          http.Redirect(w, r, red, http.StatusSeeOther)
          return
        } else {
          if data.Lan == "ja" {
            data.Tit = "未検出"
            data.Err = "このURLを見つけられませんでした。"
          } else {
            data.Tit = "Not found"
            data.Err = "This URL could not be found."
          }
          tmpl = template.Must(template.ParseFiles(cnf.webpath + "/view/404.html", cnf.webpath + "/view/header.html", cnf.webpath + "/view/footer.html"))
        }
      } else if uri == "/" && qnewurl != "" {
        data.Url = qnewurl
        data.Dom = cnf.domain
        if data.Lan == "ja" {
          data.Tit = "短縮済み"
        } else {
          data.Tit = "Shortened"
        }
        tmpl = template.Must(template.ParseFiles(cnf.webpath + "/view/submitted.html", cnf.webpath + "/view/header.html", cnf.webpath + "/view/footer.html"))
      } else {
        if data.Lan == "ja" {
          data.Tit = "未検出"
          data.Err = "このURLを見つけられませんでした。"
        } else {
          data.Tit = "Not found"
          data.Err = "This URL could not be found."
        }
        tmpl = template.Must(template.ParseFiles(cnf.webpath + "/view/404.html", cnf.webpath + "/view/header.html", cnf.webpath + "/view/footer.html"))
      }
    }

    tmpl.Execute(w, data)
    data = nil
  })

  fmt.Println(fmt.Sprint("http://127.0.0.1:", port, " でサーバーを実行中。終了するには、CTRL+Cを押して下さい。"))
  http.ListenAndServe(fmt.Sprint(":", port), nil)
}
