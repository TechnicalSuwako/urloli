#!/bin/sh

su
mkdir /etc/urloli
cp links.sample.json /etc/urloli/links.json

go build

mv /home/urloli/urlo.li/urloli /usr/local/bin
chown -R urloli:urloli /usr/local/bin/urloli
chown -R urloli:urloli /etc/urloli

unamestr=$(uname)

if [ "$unamestr" = 'Linux' ]; then
  platform=$(cat /etc/os-release | grep "^ID=")
  if [ "$platform" = 'ID=debian' -o "$platform" = "ID=devuan" -o "$platform" = "ID=ubuntu" ]; then
    apt update && apt install certbot
  elif [ "$platform" = "ID=arch" -o "$platform" = "ID=artix" ]; then
    pacman -S certbot
  elif [ "$platform" = "ID=centos" -o  "$platform" = "ID=rhel" ]; then
    dnf install certbot
  fi
  certbot certonly --webroot urlo.li www.urlo.li
  cp srv/linux/etc/nginx/sites-enabled/urloli.conf /etc/nginx/sites-enabled
  cp srv/linux/etc/init.d/urloli /etc/init.d
  chmod +x /etc/init.d/urloli
  /etc/init.d/urloli start
elif [ "$unamestr" = 'OpenBSD' ]; then
  pkg_add certbot
  cat /etc/acme-client.conf src/openbsd/etc/acme-client.conf > /etc/acme-client.conf
  cat /etc/httpd.conf srv/openbsd/etc/httpd.conf > /etc/httpd.conf
  cp srv/openbsd/etc/rc.d/urloli /etc/rc.d
  chmod +x /etc/rc.d/urloli
  rcctl start urloli
fi

exit
