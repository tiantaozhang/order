#!/bin/sh
a=`pwd`
cd ../../../
export GOPATH=$GOPATH:`pwd`:/Users/xxx/code/go/golib/src/go-alipay
cd $a

# export GOPATH=`pwd`:$GOPATH
# CUR=`pwd`
# echo "GOPATH = "$GOPATH

echo $GOPATH
go test -coverprofile=c.out > resultLog.txt
go tool cover -html=c.out -o coverage.html