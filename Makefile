
all:
	go build

release: all git-tag
	tar jcf modbus-demo-`cat VERSION`.tar.bz2 modbustcp modbustcpd LICENSE README.md

git-tag: bump
	git tag `cat VERSION`
	git push --tags

bump:
	echo `cat VERSION`+.1 |bc  > VERSION.new
	rm VERSION
	mv VERSION.new VERSION

upload:
	scp modbus-demo-`cat VERSION`.tar.bz2 oplerno:/var/lib/lxd/containers/ateps-updates/rootfs/var/www/portage/distfiles/
