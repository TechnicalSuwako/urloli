package main

import (
  "text/template"
  "fmt"
  "net/http"
  "encoding/json"
  "strings"
  "log"
  "os"
  "path/filepath"
  "gitler.moe/suwako/goliblocale"
)

type (
  Page struct {
    Tit, Err, Url, Dom, Lan, Ver, Ves string
    i18n map[string]string
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

func (p Page) T (key string) string {
  return p.i18n[key]
}

func initloc (r *http.Request) string {
  supportLang := map[string]bool{
    "ja": true,
    "en": true,
  }

	cookie, err := r.Cookie("lang")
  if err != nil {
	  return "ja"
	}

  if _, ok := supportLang[cookie.Value]; ok {
    return cookie.Value
  } else {
    return "ja"
  }
}

func serv (cnf Config, port int) {
  dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
  if err != nil {
    log.Fatal(err)
  }
  err = os.Chdir(dir)
  if err != nil {
    log.Fatal(err)
  }

  http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(cnf.webpath + "/static"))))
  ftmpl := []string{cnf.webpath + "/view/index.html", cnf.webpath + "/view/header.html", cnf.webpath + "/view/footer.html"}

  http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(200)
    buf, _ := json.MarshalIndent(&Stat{Url: cnf.domain, Ver: version}, "", "  ")
    _, _ = w.Write(buf)
  })

  http.HandleFunc("/api/lolify", func(w http.ResponseWriter, r *http.Request) {
    lang := initloc(r)
    i18n, err := goliblocale.GetLocale("locale/" + lang)
    if err != nil {
      fmt.Println("liblocaleエラー：%v", err)
      return
    }

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
            res = &Api{Cod: 400, Err: i18n["errfusei"]}
          }
          if !chklim {
            res = &Api{Cod: 400, Err: i18n["errcharlim"]}
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
          res = &Api{Cod: 400, Err: i18n["errurlent"]}
        }
      }
    }

    buf, _ := json.MarshalIndent(res, "", "  ")
    _, _ = w.Write(buf)
  })

  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    data := &Page{Ver: version, Ves: strings.ReplaceAll(version, ".", "")}
    uri := r.URL.Path
    lang := initloc(r)

    i18n, err := goliblocale.GetLocale("locale/" + lang)
    if err != nil {
      fmt.Println("liblocaleエラー：%v", err)
      return
    }
    data.i18n = i18n
    data.Lan = lang

    // デフォルトページ＝未検出
    data.Tit = i18n["mikensyutu"]
    data.Err = i18n["errurlnai"]
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
            data.Tit = i18n["fuseiurl"]
            if !chkprx {
              data.Err = i18n["errfusei"]
            } else if !chklim {
              data.Err = i18n["errcharlim"]
            }
          } else {
            chkfn, _ := geturl(addurl, cnf.linkpath, true)
            if chkfn != "" {
              http.Redirect(w, r, addurl, http.StatusSeeOther)
              return
            } else {
              data.Url = insertjson(addurl, cnf.linkpath)
              data.Dom = cnf.domain
              data.Tit = i18n["tansyukuzumi"]
              ftmpl[0] = cnf.webpath + "/view/submitted.html"
            }
          }
        } else {
          data.Err = i18n["errurlent"]
        }
      } else if r.PostForm.Get("langchange") != "" {
        lang := r.PostForm.Get("lang")
        http.SetCookie(w, &http.Cookie{Name: "lang", Value: lang, MaxAge: 31536000, Path: "/"})
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
      }
    } else { // r.Method == "GET"
      if uri == "/" {
        data.Tit = i18n["top"]
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

  fmt.Println(fmt.Sprint("http://" + cnf.ip + ":", port, " でサーバーを実行中。終了するには、CTRL+Cを押して下さい。"))
  http.ListenAndServe(fmt.Sprint(cnf.ip + ":", port), nil)
}
