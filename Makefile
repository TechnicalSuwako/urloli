UNAME_S != uname -s

NAME != cat main.go | grep "var sofname" | awk '{print $$4}' | sed "s/\"//g"
VERSION != cat main.go | grep "var version" | awk '{print $$4}' | sed "s/\"//g"

PREFIX = /usr/local
.if ${UNAME_S} == "Linux"
PREFIX = /usr
.endif

MANPREFIX = ${PREFIX}/share/man
.if ${UNAME_S} == "OpenBSD"
MANPREFIX = ${PREFIX}/man
.endif

CNFPREFIX=/etc
.if ${UNAME_S} == "FreeBSD" || ${UNAME_S} == "NetBSD" || ${UNAME_S} == "Dragonfly"
CNFPREFIX=${PREFIX}/etc
.endif

CC=CGO_ENABLED=0 go build
RELEASE=-ldflags="-s -w" -buildvcs=false

all:
	${CC} ${RELEASE} -o ${NAME}

release:
	mkdir -p release/bin
	env GOOS=linux   GOARCH=amd64   ${CC} ${RELEASE} -o \
		release/bin/${NAME}-${VERSION}-linux-amd64
	env GOOS=linux   GOARCH=arm64   ${CC} ${RELEASE} -o \
		release/bin/${NAME}-${VERSION}-linux-arm64
	env GOOS=linux   GOARCH=riscv64 ${CC} ${RELEASE} -o \
		release/bin/${NAME}-${VERSION}-linux-riscv64
	env GOOS=linux   GOARCH=ppc64   ${CC} ${RELEASE} -o \
		release/bin/${NAME}-${VERSION}-linux-ppc64
	env GOOS=linux   GOARCH=mips64  ${CC} ${RELEASE} -o \
		release/bin/${NAME}-${VERSION}-linux-mips64
	env GOOS=openbsd GOARCH=amd64   ${CC} ${RELEASE} -o \
		release/bin/${NAME}-${VERSION}-openbsd-amd64
	env GOOS=openbsd GOARCH=arm64   ${CC} ${RELEASE} -o \
		release/bin/${NAME}-${VERSION}-openbsd-arm64
	env GOOS=openbsd GOARCH=mips64  ${CC} ${RELEASE} -o \
		release/bin/${NAME}-${VERSION}-openbsd-mips64
	env GOOS=openbsd GOARCH=ppc64   ${CC} ${RELEASE} -o \
		release/bin/${NAME}-${VERSION}-openbsd-ppc64
	env GOOS=openbsd GOARCH=riscv64 ${CC} ${RELEASE} -o \
		release/bin/${NAME}-${VERSION}-openbsd-riscv64
	env GOOS=openbsd GOARCH=sparc64 ${CC} ${RELEASE} -o \
		release/bin/${NAME}-${VERSION}-openbsd-sparc64
	env GOOS=freebsd GOARCH=amd64   ${CC} ${RELEASE} -o \
		release/bin/${NAME}-${VERSION}-freebsd-amd64
	env GOOS=freebsd GOARCH=arm64   ${CC} ${RELEASE} -o \
		release/bin/${NAME}-${VERSION}-freebsd-arm64
	env GOOS=freebsd GOARCH=riscv64 ${CC} ${RELEASE} -o \
		release/bin/${NAME}-${VERSION}-freebsd-riscv64
	env GOOS=netbsd GOARCH=amd64    ${CC} ${RELEASE} -o \
		release/bin/${NAME}-${VERSION}-netbsd-amd64
	env GOOS=netbsd GOARCH=arm64    ${CC} ${RELEASE} -o \
		release/bin/${NAME}-${VERSION}-netbsd-arm64
	env GOOS=illumos GOARCH=amd64   ${CC} ${RELEASE} -o \
		release/bin/${NAME}-${VERSION}-illumos-amd64

clean:
	rm -f ${NAME} ${NAME}-${VERSION}.tar.gz

dist: clean
	mkdir -p ${NAME}-${VERSION} release/src
	cp -R LICENSE.txt Makefile README.md CHANGELOG.md\
		view static logo.jpg\
		${NAME}.1 *.go *.json go.mod go.sum ${NAME}-${VERSION}
	tar zcfv release/src/${NAME}-${VERSION}.tar.gz ${NAME}-${VERSION}
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
