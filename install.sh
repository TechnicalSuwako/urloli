#!/bin/sh

unamestr=$(uname)
domain="$1"

case "$domain" in
  *.i2p)   network="i2p" ;;
  *.onion) network="tor" ;;
  *)       network="www" ;;
esac

if [ "$unamestr" = 'FreeBSD' ]; then
  mkdir /usr/local/etc/urloli
  cp -i links.sample.json /usr/local/etc/urloli/links.json
  cp -i config.json /usr/local/etc/urloli/config.json
  sed -i .orig "s/urlo\.li/$domain/g" /usr/local/etc/urloli/config.json
  rm -rf /usr/local/etc/urloli/config.json.orig
  go build -buildvcs=false
else
  mkdir /etc/urloli
  cp -i links.sample.json /etc/urloli/links.json
  cp -i config.json /etc/urloli/config.json
  sed -i "s/urlo\.li/$domain/g" /etc/urloli/config.json
  go build
fi

mv -i urloli /usr/local/bin

if [ "$unamestr" = 'Linux' ]; then
  platform=$(cat /etc/os-release | grep "^ID=")
  if [ "$platform" = 'ID=debian' -o "$platform" = "ID=devuan" -o "$platform" = "ID=ubuntu" ]; then
    apt update && apt install certbot
  elif [ "$platform" = "ID=arch" -o "$platform" = "ID=artix" ]; then
    pacman -S certbot
  elif [ "$platform" = "ID=centos" -o  "$platform" = "ID=rhel" ]; then
    dnf install certbot
  fi
  if [ "$network" = 'www' ]; then
    certbot certonly --webroot -d $domain -d www.$domain
    cp -i srv/linux/etc/nginx/sites-enabled/urloli-clear.conf /etc/nginx/sites-enabled/urloli.conf
  else
    cp -i srv/linux/etc/nginx/sites-enabled/urloli-dark.conf /etc/nginx/sites-enabled/urloli.conf
  fi
  sed -i "s/urlo\.li/$domain/g" /etc/nginx/sites-enabled/urloli.conf
  cp -i srv/linux/etc/init.d/urloli /etc/init.d
  chmod +x /etc/init.d/urloli
  /etc/init.d/urloli start
elif [ "$unamestr" = 'OpenBSD' ]; then
  if [ "$network" = 'www' ]; then
    cat /etc/acme-client.conf srv/openbsd/etc/acme-client.conf > /etc/acme-client.conf
    sed -i "s/urlo\.li/$domain/g" /etc/acme-client.conf
    cat /etc/httpd.conf srv/openbsd/etc/httpd-clear.conf > /etc/httpd.conf
  else
    cat /etc/httpd.conf srv/openbsd/etc/httpd-dark.conf > /etc/httpd.conf
  fi
  sed -i "s/urlo\.li/$domain/g" /etc/httpd.conf
  rcctl restart httpd
  if [ "$network" = 'www' ]; then
    acme-client -v $domain
  fi
  cp -i srv/openbsd/etc/rc.d/urloli /etc/rc.d
  chmod +x /etc/rc.d/urloli
  rcctl start urloli
elif [ "$unamestr" = 'FreeBSD' ]; then
  pkg install py39-certbot
  if [ "$network" = 'www' ]; then
    certbot certonly --webroot -d $domain -d www.$domain
    cp -i srv/linux/etc/nginx/sites-enabled/urloli-clear.conf /usr/local/etc/nginx/sites-enabled/urloli.conf
  else
    cp -i srv/linux/etc/nginx/sites-enabled/urloli-dark.conf /usr/local/etc/nginx/sites-enabled/urloli.conf
  fi
  sed -i .orig "s/urlo\.li/$domain/g" /usr/local/etc/nginx/sites-enabled/urloli.conf
  rm -rf /usr/local/etc/nginx/sites-enabled/urloli.conf.orig
  cp -i srv/freebsd/usr/local/etc/rc.d/urloli /usr/local/etc/rc.d
  chmod +x /usr/local/etc/rc.d/urloli
  sysrc urloli_enable=YES
  service start urloli
fi

exit
