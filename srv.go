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
  ftmpl := []string{cnf.webpath + "/view/index.html", cnf.webpath + "/view/header.html", cnf.webpath + "/view/footer.html"}

  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    data := &Page{Ver: version}
    cookie, err := r.Cookie("lang")
    if err != nil {
      data.Lan = "ja"
    } else {
      data.Lan = cookie.Value
    }

    uri := r.URL.Path
    query := r.URL.Query()
    qnewurl := query.Get("newurl")
    data.Tit = getloc("top", data.Lan)
    ftmpl[0] = cnf.webpath + "/view/index.html"
    tmpl := template.Must(template.ParseFiles(ftmpl[0], ftmpl[1], ftmpl[2]))

    if r.Method == "POST" {
      err := r.ParseForm()
      if err != nil { fmt.Println(err) }
      if r.PostForm.Get("sosin") != "" {
        if r.PostForm.Get("newadd") != "" {
          addurl := r.PostForm.Get("newadd")
          chkprx := checkprefix(addurl)
          chklim := checkcharlim(addurl)
          if !chkprx {
            data.Tit = getloc("fuseiurl", data.Lan)
            data.Err = getloc("errfusei", data.Lan)
            ftmpl[0] = cnf.webpath + "/view/404.html"
            tmpl = template.Must(template.ParseFiles(ftmpl[0], ftmpl[1], ftmpl[2]))
          }
          if !chklim {
            data.Tit = getloc("fuseiurl", data.Lan)
            data.Err = getloc("errcharlim", data.Lan)
            ftmpl[0] = cnf.webpath + "/view/404.html"
            tmpl = template.Must(template.ParseFiles(ftmpl[0], ftmpl[1], ftmpl[2]))
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
              data.Tit = getloc("tansyukuzumi", data.Lan)
              ftmpl[0] = cnf.webpath + "/view/submitted.html"
              tmpl = template.Must(template.ParseFiles(ftmpl[0], ftmpl[1], ftmpl[2]))
            }
          }
        } else {
          data.Tit = getloc("mikensyutu", data.Lan)
          data.Err = getloc("errurlent", data.Lan)
          ftmpl[0] = cnf.webpath + "/view/404.html"
          tmpl = template.Must(template.ParseFiles(ftmpl[0], ftmpl[1], ftmpl[2]))
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
        ftmpl[0] = cnf.webpath + "/view/index.html"
        tmpl = template.Must(template.ParseFiles(ftmpl[0], ftmpl[1], ftmpl[2]))
      } else if uri != "/" && qnewurl == "" {
        red, _ := geturl(uri[1:], cnf.linkpath, false)
        if red != "" {
          http.Redirect(w, r, red, http.StatusSeeOther)
          return
        } else {
          data.Tit = getloc("mikensyutu", data.Lan)
          data.Err = getloc("errurlnai", data.Lan)
          ftmpl[0] = cnf.webpath + "/view/404.html"
          tmpl = template.Must(template.ParseFiles(ftmpl[0], ftmpl[1], ftmpl[2]))
        }
      } else if uri == "/" && qnewurl != "" {
        data.Url = qnewurl
        data.Dom = cnf.domain
        data.Tit = getloc("tansyukuzumi", data.Lan)
        ftmpl[0] = cnf.webpath + "/view/submitted.html"
        tmpl = template.Must(template.ParseFiles(ftmpl[0], ftmpl[1], ftmpl[2]))
      } else {
        data.Tit = getloc("mikensyutu", data.Lan)
        data.Err = getloc("errurlnai", data.Lan)
        ftmpl[0] = cnf.webpath + "/view/404.html"
        tmpl = template.Must(template.ParseFiles(ftmpl[0], ftmpl[1], ftmpl[2]))
      }
    }

    tmpl.Execute(w, data)
    data = nil
  })

  fmt.Println(fmt.Sprint("http://127.0.0.1:", port, " でサーバーを実行中。終了するには、CTRL+Cを押して下さい。"))
  http.ListenAndServe(fmt.Sprint(":", port), nil)
}
