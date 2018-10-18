
PROGNAME="modbus-demo"
all:
	go build

release: all git-tag
	mkdir ${PROGNAME}-`cat VERSION`
	cp modbustcp modbustcpd LICENSE README.md ${PROGNAME}-`cat VERSION`
	tar jcf ${PROGNAME}-`cat VERSION`.tar.bz2 ${PROGNAME}-`cat VERSION`
	rm -rf ${PROGNAME}-`cat VERSION`

src:
	tar jcf ${PROGNAME}-src-`cat VERSION`.tar.bz2 *go README.md LICENSE Makefile VERSION modbustcpd

git-tag: bump
	git tag `cat VERSION`
	git push --tags

bump:
	echo `cat VERSION`+.1 |bc  > VERSION.new
	rm VERSION
	mv VERSION.new VERSION

upload:
	scp ${PROGNAME}-`cat VERSION`.tar.bz2 oplerno:/var/lib/lxd/containers/ateps-updates/rootfs/var/www/portage/distfiles/

windows:
	GOOS=windows GOARCH=386 go build
