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
}

func serv (cnf Config, port int) {
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
            if cookie.Value == "ja" {
                data = &Page{Tit: "不正なURL", Err: "URLは「http://」又は「https://」で始めます。", Lan: cookie.Value}
            } else {
              data = &Page{Tit: "Invalid URL", Err: "The URL should start with \"http://\" or \"https://\".", Lan: cookie.Value}
            }
            tmpl = template.Must(template.ParseFiles(cnf.webpath + "/view/404.html", cnf.webpath + "/view/header.html", cnf.webpath + "/view/footer.html"))
          }
          if !chklim {
            if cookie.Value == "ja" {
              data = &Page{Tit: "不正なURL", Err: "URLは500文字以内です。", Lan: cookie.Value}
            } else {
              data = &Page{Tit: "Invalid URL", Err: "The URL should be less than 500 characters.", Lan: cookie.Value}
            }
            data = &Page{Tit: "不正なURL", Err: ""}
            tmpl = template.Must(template.ParseFiles(cnf.webpath + "/view/404.html", cnf.webpath + "/view/header.html", cnf.webpath + "/view/footer.html"))
          }

          if chklim && chkprx {
            chkfn, _ := geturl(addurl, cnf.linkpath, true)
            if chkfn != "" {
              http.Redirect(w, r, addurl, http.StatusSeeOther)
              return
            } else {
              res := insertjson(addurl, cnf.linkpath)
              if cookie.Value == "ja" {
                data = &Page{Tit: "短縮済み", Lan: cookie.Value, Url: res, Dom: cnf.domain}
              } else {
                data = &Page{Tit: "Shortened", Lan: cookie.Value, Url: res, Dom: cnf.domain}
              }
              tmpl = template.Must(template.ParseFiles(cnf.webpath + "/view/submitted.html", cnf.webpath + "/view/header.html", cnf.webpath + "/view/footer.html"))
            }
          }
        } else {
          if cookie.Value == "ja" {
            data = &Page{Tit: "未検出", Err: "URLをご入力下さい。", Lan: cookie.Value}
          } else {
            data = &Page{Tit: "Not found", Err: "Please enter a URL.", Lan: cookie.Value}
          }
          tmpl = template.Must(template.ParseFiles(cnf.webpath + "/view/404.html", cnf.webpath + "/view/header.html", cnf.webpath + "/view/footer.html"))
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
        tmpl = template.Must(template.ParseFiles(cnf.webpath + "/view/index.html", cnf.webpath + "/view/header.html", cnf.webpath + "/view/footer.html"))
      } else if uri != "/" && qnewurl == "" {
        red, _ := geturl(uri[1:], cnf.linkpath, false)
        if red != "" {
          http.Redirect(w, r, red, http.StatusSeeOther)
          return
        } else {
          if cookie.Value == "ja" {
            data = &Page{Tit: "未検出", Err: "このURLを見つけられませんでした。", Lan: cookie.Value}
          } else {
            data = &Page{Tit: "Not found", Err: "This URL could not be found.", Lan: cookie.Value}
          }
          tmpl = template.Must(template.ParseFiles(cnf.webpath + "/view/404.html", cnf.webpath + "/view/header.html", cnf.webpath + "/view/footer.html"))
        }
      } else if uri == "/" && qnewurl != "" {
        data = &Page{Tit: "短縮済み", Url: qnewurl, Dom: cnf.domain}
        tmpl = template.Must(template.ParseFiles(cnf.webpath + "/view/submitted.html", cnf.webpath + "/view/header.html", cnf.webpath + "/view/footer.html"))
      } else {
        if cookie.Value == "ja" {
          data = &Page{Tit: "未検出", Err: "このURLを見つけられませんでした。", Lan: cookie.Value}
        } else {
          data = &Page{Tit: "Not found", Err: "This URL could not be found.", Lan: cookie.Value}
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
