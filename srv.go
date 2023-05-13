package main

import (
  "text/template"
  "fmt"
  "net/http"
  "encoding/json"
  "strings"
)

type (
  Page struct {
    Tit string
    Err string
    Url string
    Dom string
    Lan string
    Ver string
    Ves string
  }
  Api struct {
    Cod int `json:"code"`
    Err string `json:"error"`
    Url string `json:"url"`
    Mot string `json:"origin"`
    New bool `json:"isnew"`
  }
  Stat struct {
    Url string `json:"url"`
    Ver string `json:"version"`
  }
)

func serv (cnf Config, port int) {
  http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
  ftmpl := []string{cnf.webpath + "/view/index.html", cnf.webpath + "/view/header.html", cnf.webpath + "/view/footer.html"}

  http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(200)
    buf, _ := json.MarshalIndent(&Stat{Url: cnf.domain, Ver: version}, "", "  ")
    _, _ = w.Write(buf)
  })

  http.HandleFunc("/api/lolify", func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(200)
    res := &Api{Cod: 500, Err: "未対応"}
    if r.Method == "POST" {
      err := r.ParseForm()
      if err != nil {
        fmt.Println(err)
        res.Err = "失敗"
      } else {
        if r.PostForm.Get("url") != "" {
          addurl := r.PostForm.Get("url")
          chkprx := checkprefix(addurl)
          chklim := checkcharlim(addurl)
          if !chkprx {
            res = &Api{Cod: 400, Err: getloc("errfusei", "ja")}
          }
          if !chklim {
            res = &Api{Cod: 400, Err: getloc("errcharlim", "ja")}
          }

          if chklim && chkprx {
            chkfn, key := geturl(addurl, cnf.linkpath, true)
            if chkfn != "" {
              res = &Api{Cod: 200, Url: cnf.domain + "/" + key, Mot: addurl, New: false}
            } else {
              res = &Api{Cod: 200, Url: cnf.domain + "/" + insertjson(addurl, cnf.linkpath), Mot: addurl, New: true}
            }
          }
        } else {
          res = &Api{Cod: 400, Err: getloc("errurlent", "ja")}
        }
      }
    }

    buf, _ := json.MarshalIndent(res, "", "  ")
    _, _ = w.Write(buf)
  })

  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    data := &Page{Ver: version, Ves: strings.ReplaceAll(version, ".", "")}
    uri := r.URL.Path

    cookie, err := r.Cookie("lang")
    if err != nil {
      data.Lan = "ja"
    } else {
      data.Lan = cookie.Value
    }

    // デフォルトページ＝未検出
    data.Tit = getloc("mikensyutu", data.Lan)
    data.Err = getloc("errurlnai", data.Lan)
    ftmpl[0] = cnf.webpath + "/view/404.html"
    tmpl := template.Must(template.ParseFiles(ftmpl[0], ftmpl[1], ftmpl[2]))

    if r.Method == "POST" {
      err := r.ParseForm()
      if err != nil { fmt.Println(err) }
      if r.PostForm.Get("sosin") != "" {
        if r.PostForm.Get("newadd") != "" {
          addurl := r.PostForm.Get("newadd")
          chkprx := checkprefix(addurl)
          chklim := checkcharlim(addurl)
          if !chkprx || !chklim {
            data.Tit = getloc("fuseiurl", data.Lan)
            if !chkprx {
              data.Err = getloc("errfusei", data.Lan)
            } else if !chklim {
              data.Err = getloc("errcharlim", data.Lan)
            }
          } else {
            chkfn, _ := geturl(addurl, cnf.linkpath, true)
            if chkfn != "" {
              http.Redirect(w, r, addurl, http.StatusSeeOther)
              return
            } else {
              data.Url = insertjson(addurl, cnf.linkpath)
              data.Dom = cnf.domain
              data.Tit = getloc("tansyukuzumi", data.Lan)
              ftmpl[0] = cnf.webpath + "/view/submitted.html"
            }
          }
        } else {
          data.Err = getloc("errurlent", data.Lan)
        }
      } else if r.PostForm.Get("langchange") != "" {
        cookie, err := r.Cookie("lang")
        if err != nil || cookie.Value == "ja" {
          http.SetCookie(w, &http.Cookie{Name: "lang", Value: "en", MaxAge: 31536000, Path: "/"})
        } else {
          http.SetCookie(w, &http.Cookie{Name: "lang", Value: "ja", MaxAge: 31536000, Path: "/"})
        }
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
      }
    } else { // r.Method == "GET"
      if uri == "/" {
        data.Tit = getloc("top", data.Lan)
        ftmpl[0] = cnf.webpath + "/view/index.html"
      } else {
        red, _ := geturl(uri[1:], cnf.linkpath, false)
        if red != "" {
          http.Redirect(w, r, red, http.StatusSeeOther)
          return
        }
      }
    } // r.Method

    tmpl = template.Must(template.ParseFiles(ftmpl[0], ftmpl[1], ftmpl[2]))
    tmpl.Execute(w, data)
    data = nil
  })

  fmt.Println(fmt.Sprint("http://127.0.0.1:", port, " でサーバーを実行中。終了するには、CTRL+Cを押して下さい。"))
  http.ListenAndServe(fmt.Sprint(":", port), nil)
}
