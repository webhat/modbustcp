# Copyright (c) 2018-2019 Daniel W. Crompton, Special Brands Holding BV
# 
# MIT License
# 
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
# 
# The above copyright notice and this permission notice shall be included in all
# copies or substantial portions of the Software.
# 
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
# SOFTWARE.

PROGNAME="modbus-demo"
TOOL="modbustcp"
OS := $(shell uname -s)

all: test windows linux darwin
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

test:
	go test -v -test.failfast -coverprofile cover.out
