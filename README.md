# URLロリ
クッソ小さいURL短縮作成ソフトだわ〜♡

## インストールする方法

### 従属ソフト

* Go 1.19以上
* nginx又はOpenBSDのrelayd
* 良いOS (GNU/Linux、OpenBSD、又はFreeBSD)

## インストールする方法

### 全部（opendoasを使わなければ、sudoをご利用、又はopendoasをインストールして下さい）

```sh
make
doas make install
```

### OpenBSD

```sh
nvim /etc/rc.d/urloli
```

```
#!/bin/ksh

daemon="/usr/local/bin/urloli -s"

. /etc/rc.d/rc.subr

rc_bg=YES
rc_reload=NO

rc_cmd $1
```

```sh
chmod +x /etc/rc.d/urloli
rcctl enable urloli
rcctl start urloli
```

### FreeBSD

```sh
nvim /usr/local/etc/rc.d/urloli
```

```
#!/bin/sh

# PROVIDE: urloli
# REQUIRE: NETWORKING SYSLOG
# KEYWORD: shutdown
#
# Add the following lines to /etc/rc.conf to enable urloli:
#
#urloli_enable="YES"

. /etc/rc.subr

name="urloli"
rcvar="urloli_enable"

load_rc_config $name

: ${urloli_enable:="NO"}
: ${urloli_facility:="daemon"}
: ${urloli_priority:="debug"}

command="/usr/local/bin/${name}"
procname="/usr/local/bin/${name}"

pidfile="/var/run/${name}.pid"

start_cmd="${name}_start"

urloli_start() {
  for d in /var/db/urloli /var/log/urloli; do
    if [ ! -e "$d" ]; then
      mkdir "$d"
    fi
  done
  /usr/sbin/daemon -S -l ${urloli_facility} -s ${urloli_priority} -T ${name} \
    -p ${pidfile} \
    /usr/bin/env -i \
    "PATH=/usr/local/bin:${PATH}" \
    $command
}

run_rc_command "$1"
```

```sh
sysrc urloli_enable=YES
service start urloli
```

### Crux

```sh
nvim /etc/rc.d/urloli
```

```
#!/bin/sh
#
# /etc/rc.d/urloli: start/stop the urloli daemon
#

SSD=/sbin/start-stop-daemon
NAME=urloli
PROG=/usr/bin/$NAME
PIOD=/run/$NAME.pid

case $1 in
start)
  $SSD --start --pidfile $PID --exec $PROG
  ;;
stop)
  $SSD --stop --retry 10 --pidfile $PID
  ;;
restart)
  $0 stop
  $0 start
  ;;
status)
  $SSD --status --pidfile $PID
  case $? in
  0) echo "$PROG は実行中。pid $(cat $PID)" ;;
  1) echo "$PROG は実行していませんが、pidファイルは「 $PID 」として存在しそう" ;;
  3) echo "$PROG は停止中" ;;
  4) echo "状況不明" ;;
  esac
  ;;
*)
  echo "usage: $0 [start|sto@|restart|status]"
  ;;
esac

# End of file
```

### Devuan/Debian/Ubuntu/Arch/Artix/AlmaLinux等

```sh
nvim /etc/init.d/urloli
```

```
#!/bin/sh
#
# chkconfig: 35 90 12
# description: URL Loli server
#

NAME=urloli
DESC=urloli
DAEMON=/usr/bin/$NAME

start () {
  echo "URLロリサーバーは開始中：\n"
  /usr/bin/urloli -s 9910 &>/dev/null &
  touch /var/lock/subsys/urloli
  echo
}

stop () {
  echo "URLロリサーバーは終了中：\n"
  pkill urloli
  rm -f /var/lock/subsys/urloli
  echo
}

case "$1" in
  start)
    start
    ;;
  stop)
    stop
    ;;
  status)
    status_of_proc "$DAEMON" "$NAME" && exit 0 || exit $?
    ;;
  restart|reload|condrestart)
    stop
    start
    ;;
  *)
    echo $"Usage: $0 {start|stop|restart|status}"
    exit 1
esac
```

## ウェブサーバー

### OpenBSD

```sh
nvim /etc/relayd.conf
```

```
# $OpenBSD: relayd.conf,v 1.5 2018/05/06 20:56:55 benno Exp $
#
relayd_address="0.0.0.0"

table <urloli> { 127.0.0.1 }

http protocol reverse_proxy {
  tls keypair "DOMAIN"
  match request header append "X-Forwarded-For" value "$REMOTE_ADDR"
  match request header append "X-Forwarded-Port" value "$REMOTE_PORT"

  match response header set "Referrer-Policy" value "same-origin"
  match response header set "X-Frame-Options" value "deny"
  match response header set "X-XSS-Protection" value "1; mode=block"
  match response header set "X-Content-Type-Options" value "nosniff"
  match response header set "Strict-Transport-Security" value "max-age=31536000; includeSubDomains; preload"
  match response header set "Cache-Control" value "max-age=86400"

  pass request quick header "Host" value "urlo.li" forward to <urloli>

  return error
  pass
}

relay www {
  listen on $relayd_address port 443 tls 
  protocol $relayd_address

  forward to <urloli> check tcp port 9910
}
```

### その他

```sh
server {
  server_name DOMAIN www.DOMAIN;
  root   /var/www/htdocs/urloli;

  access_log off;
  error_log off;

  if ($host = www.DOMAIN) {
    return 301 https://DOMAIN$request_uri;
  }

  location /static {
    try_files $uri $uri/ /static/$args;
  }

  location / {
    proxy_pass http://localhost:9910;
  }

  listen [::]:443 ssl ipv6only=on;
  listen 443 ssl;
  ssl_certificate /etc/letsencrypt/live/DOMAIN/fullchain.pem;
  ssl_certificate_key /etc/letsencrypt/live/DOMAIN/privkey.pem;
  include /etc/letsencrypt/options-ssl-nginx.conf
}

server {
  if ($host = DOMAIN) {
    return 301 https://DOMAIN$request_uri;
  }

  if ($host = www.DOMAIN) {
    return 301 https://DOMAIN$request_uri;
  }

  listen 80;
  listen [::]:80;
  server_name DOMAIN www.DOMAIN;
  return 404;
}
```
