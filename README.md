# URLロリ
クッソ小さいURL短縮作成ソフトだわ〜♡

## 使い方
```sh
cp links.sample.json links.json
nvim links.json

useradd -m -s /usr/local/bin/zsh urloli
su -l urloli
git clone https://gitler.moe/TechnicalSuwako/urlo.li.git && cd urlo.li
go build
exit

mv /home/urloli/urlo.li/urloli /usr/local/bin
chown -R urloli:urloli /usr/local/bin/urloli

rcctl start urloli
```

### nginxコンフィグ（Linux、FreeBSD）
```
server {
  server_name urlo.li www.urlo.li;

  access_log off;
  error_log off;

  if ($host = www.urlo.li) {
    return 301 https://urlo.li$request_uri;
  }

  location / {
    proxy_pass http://localhost:9910;
  }

  listen [::]:443 ssl ipv6only=on;
  listen 443 ssl;
  ssl_certificate /etc/letsencrypt/live/urlo.li/fullchain.pem;
  ssl_certificate_key /etc/letsencrypt/live/urlo.li/privkey.pem;
  include /etc/letsencrypt/options-ssl-nginx.conf
}

server {
  if ($host = urlo.li) {
    return 301 https://urlo.li$request_uri;
  }

  if ($host = www.urlo.li) {
    return 301 https://urlo.li$request_uri;
  }

  listen 80;
  listen [::]:80;
  server_name urlo.li www.urlo.li;
  return 404;
}
```

### OpenHTTPdコンフィグ（OpenBSD）
```
server "urlo.li" {
  listen on $ext_addr port 80
  block return 301 "https://$SERVER_NAME$REQUEST_URI"
}
server "urlo.li" {
  listen on $ext_addr tls port 443
  tls {
    certificate     "/etc/letsencrypt/live/urlo.li/fullchain.pem"
    key             "/etc/letsencrypt/live/urlo.li/privkey.pem"
  }
  connection { max requests 500, timeout 3600 }
  location "/*" {
    fastcgi socket tcp 127.0.0.1 9910
  }
}
```

### OpenBSDのrc
```
#!/bin/ksh
#
# $OpenBSD: urloli.rc,v 1.4 2018/01/11 19:27:11 rpe Exp $

name="urloli"
daemon="/usr/local/bin/${name}"
daemon_user="${name}"

. /etc/rc.d/rc.subr

rc_cmd $1
```

### links.jsonファイルの中に
```
{
  "hogehoge": "https://076.moe"
}
```

https://（ドメイン名）/hogehoge にアクセスすると、https://076.moe に移転されます。
