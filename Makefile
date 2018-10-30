
PROGNAME="modbus-demo"
TOOL="modbustcp"
OS := $(shell uname -s)

all: windows linux darwin
ifeq ($(OS),Darwin)
	ln -sf ${TOOL}-darwin ${TOOL}
endif
ifeq ($(OS),Linux)
	ln -sf ${TOOL}-linux ${TOOL}
endif

release: all git-tag
	mkdir ${PROGNAME}-`cat VERSION`
	cp ${TOOL} ${TOOL}d LICENSE README.md ${PROGNAME}-`cat VERSION`
	tar jcf ${PROGNAME}-`cat VERSION`.tar.bz2 ${PROGNAME}-`cat VERSION`
	rm -rf ${PROGNAME}-`cat VERSION`

src:
	tar jcf ${PROGNAME}-src-`cat VERSION`.tar.bz2 *go README.md LICENSE Makefile VERSION modbustcpd

git-tag: bump
	git tag `cat VERSION`
	git push --tags

bump:
	echo `cat VERSION`+.01 |bc  > VERSION.new
	rm VERSION
	mv VERSION.new VERSION

upload:
	scp ${PROGNAME}-`cat VERSION`.tar.bz2 oplerno:/var/lib/lxd/containers/ateps-updates/rootfs/var/www/portage/distfiles/

windows:
	GOOS=windows GOARCH=386 go build -o ${TOOL}.exe

linux:
	GOOS=linux GOARCH=amd64 go build -o ${TOOL}-linux

darwin:
	GOOS=darwin GOARCH=amd64 go	build -o ${TOOL}-darwin
