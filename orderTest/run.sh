#!/bin/bash
###### Set ENV ######
export PATH=/usr/local/go/bin:$PATH
export GOPATH=$GOPATH:`pwd`:/Users/xxx/code/go/golib/src/go-alipay
echo "export GOPATH=$GOPATH"

###### Kill Service Process ######
ps -ef|grep "ORDER"|awk '$0 !~/grep/ {print $2}' |tr -s '\n' ' '|xargs kill -9
rm ORDER

env(){
	##resovle dependencies
	echo "-------begin resolve dependencies--------"

	go get code.google.com/p/go-uuid/uuid
	go get github.com/go-sql-driver/mysql
	go get github.com/smartystreets/goconvey

	go get -u github.com/Centny/gwf/log
	go get -u github.com/Centny/gwf/routing
	go get -u github.com/Centny/gwf/util
    go get -u github.com/Centny/gwf/netw/impl
    go get -u github.com/Centny/gwf/pool
	go get -u github.com/Centny/resize
	
	svn co https://sdev.jxzy.com/svn/UAS/trunk/UCS/src/org.cny.uas src/org.cny.uas
	svn co https://sdev.jxzy.com/svn/UAS/trunk/UAPlugin/src/org.cny.uap src/org.cny.uap
    svn co https://sdev.jxzy.com/svn/RCP/trunk/050RCP/RCP/DbMgr/src/com.dy.tool src/com.dy.tool
    svn co https://sdev.jxzy.com/svn/RCP/trunk/050RCP/RCP/OrgMgr/src/com.rcp.orgMgr src/com.rcp.orgMgr


	echo "--------resove dependencies complete------"
}

run(){
	go build src/com.dy.order/main/main.go 
	mv main ORDER
	./ORDER conf/order.properties
}

l(){
	go build src/com.dy.order/main/main.go 
	mv main ORDER
	./ORDER conf/local.properties
}

t(){
	cd src/com.dy.order/
	cd $tpkg
	go test -coverprofile=c.out
    go tool cover -html=c.out -o c.html
    if [ "$openC" = "o" ];then
       open c.html
	fi
}

b(){
	cd src/com.dy.order/
	cd $tpkg
	go test -v -bench="$tfunc" -cpuprofile=cpu.prof -c
	go tool pprof $tpkg.test cpu.prof
	> web
}

case $1 in
	"run")
	run
	;;
	"l")
	l
	;;
	"env")
	env
	;;
	"t")
	tpkg=$2
	openC=$3
	t
	;;
	"b")
	tpkg=$2
	tfunc=$3
	b
	;;
	*)
	run
	;;
esac

