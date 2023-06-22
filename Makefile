NAME=urloli
VERSION := $(shell cat main.go | grep "var version" | awk '{print $$4}' | sed "s/\"//g")
# Linux、Illumos
PREFIX=/usr
# FreeBSDとOpenBSD
#PREFIX=/usr/local
# NetBSD
#PREFIX=/usr/pkg
MANPREFIX=${PREFIX}/share/man
# LinuxとOpenBSD
CNFPREFIX=/etc
# FreeBSD
#CNFPREFIX=/usr/local/etc
CC=CGO_ENABLED=0 go build
# リリース。なし＝デバッグ。
RELEASE=-ldflags="-s -w" -buildvcs=false

all:
	${CC} ${RELEASE} -o ${NAME}

release:
	echo "linux-amd64"
	env GOOS=linux GOARCH=amd64 ${CC} ${RELEASE} -o bin/${NAME}-${VERSION}-linux-amd64
	echo "linux-arm64"
	env GOOS=linux GOARCH=arm64 ${CC} ${RELEASE} -o bin/${NAME}-${VERSION}-linux-arm64
	echo "linux-arm"
	env GOOS=linux GOARCH=arm ${CC} ${RELEASE} -o bin/${NAME}-${VERSION}-linux-arm
	echo "linux-riscv64"
	env GOOS=linux GOARCH=riscv64 ${CC} ${RELEASE} -o bin/${NAME}-${VERSION}-linux-riscv64
	echo "linux-386"
	env GOOS=linux GOARCH=386 ${CC} ${RELEASE} -o bin/${NAME}-${VERSION}-linux-386
	echo "linux-ppc64"
	env GOOS=linux GOARCH=ppc64 ${CC} ${RELEASE} -o bin/${NAME}-${VERSION}-linux-ppc64
	echo "linux-mips64"
	env GOOS=linux GOARCH=mips64 ${CC} ${RELEASE} -o bin/${NAME}-${VERSION}-linux-mips64
	echo "openbsd-amd64"
	env GOOS=openbsd GOARCH=amd64 ${CC} ${RELEASE} -o bin/${NAME}-${VERSION}-openbsd-amd64
	echo "openbsd-386"
	env GOOS=openbsd GOARCH=386 ${CC} ${RELEASE} -o bin/${NAME}-${VERSION}-openbsd-386
	echo "openbsd-arm64"
	env GOOS=openbsd GOARCH=arm64 ${CC} ${RELEASE} -o bin/${NAME}-${VERSION}-openbsd-arm64
	echo "openbsd-arm" 
	env GOOS=openbsd GOARCH=arm ${CC} ${RELEASE} -o bin/${NAME}-${VERSION}-openbsd-arm
	echo "openbsd-mips64"
	env GOOS=openbsd GOARCH=mips64 ${CC} ${RELEASE} -o bin/${NAME}-${VERSION}-openbsd-mips64
	echo "freebsd-amd64"
	env GOOS=freebsd GOARCH=amd64 ${CC} ${RELEASE} -o bin/${NAME}-${VERSION}-freebsd-amd64
	echo "freebsd-386"
	env GOOS=freebsd GOARCH=386 ${CC} ${RELEASE} -o bin/${NAME}-${VERSION}-freebsd-386
	echo "freebsd-arm"
	env GOOS=freebsd GOARCH=arm ${CC} ${RELEASE} -o bin/${NAME}-${VERSION}-freebsd-arm4
	echo "freebsd-arm64"
	env GOOS=freebsd GOARCH=arm64 ${CC} ${RELEASE} -o bin/${NAME}-${VERSION}-freebsd-arm64
	echo "freebsd-riscv64"
	env GOOS=freebsd GOARCH=riscv64 ${CC} ${RELEASE} -o bin/${NAME}-${VERSION}-freebsd-riscv64
	echo "netbsd-amd64"
	env GOOS=netbsd GOARCH=amd64 ${CC} ${RELEASE} -o bin/${NAME}-${VERSION}-netbsd-amd64
	echo "netbsd-386"
	env GOOS=netbsd GOARCH=386 ${CC} ${RELEASE} -o bin/${NAME}-${VERSION}-netbsd-386
	echo "netbsd-arm"
	env GOOS=netbsd GOARCH=arm ${CC} ${RELEASE} -o bin/${NAME}-${VERSION}-netbsd-arm4
	echo "netbsd-arm64"
	env GOOS=netbsd GOARCH=arm64 ${CC} ${RELEASE} -o bin/${NAME}-${VERSION}-netbsd-arm64
	echo "illumos-amd64"
	env GOOS=illumos GOARCH=amd64 ${CC} ${RELEASE} -o bin/${NAME}-${VERSION}-illumos-amd64

clean:
	rm -f ${NAME} ${NAME}-${VERSION}.tar.gz

dist: clean
	mkdir -p ${NAME}${VERSION}
	cp -R LICENSE.txt Makefile README.md CHANGELOG.md\
		view static logo.jpg\
		${NAME}.1 *.go *.json ${NAME}-${VERSION}
	tar -zcfv ${NAME}-${VERSION}.tar.gz ${NAME}-${VERSION}
	rm -rf ${NAME}-${VERSION}

install: all
	mkdir -p ${DESTDIR}${PREFIX}/bin
	cp -f ${NAME} ${DESTDIR}${PREFIX}/bin
	chmod 755 ${DESTDIR}${PREFIX}/bin/${NAME}
	mkdir -p ${DESTDIR}${MANPREFIX}/man1
	sed "s/VERSION/${VERSION}/g" < ${NAME}.1 > ${DESTDIR}${MANPREFIX}/man1/${NAME}.1
	chmod 644 ${DESTDIR}${MANPREFIX}/man1/${NAME}.1
	mkdir -p ${DESTDIR}${CNFPREFIX}/${NAME}
	chmod 755 ${DESTDIR}${CNFPREFIX}/${NAME}

uninstall:
	rm -f ${DESTDIOR}${PREFIX}/bin/${NAME}\
		${DESTDIR}${MANPREFIX}/man1/${NAME}.1\
		${DESTDIR}${CNFPREFIX}/${NAME}

.PHONY: all release clean dist install uninstall
