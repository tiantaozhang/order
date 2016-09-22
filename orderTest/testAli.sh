#!/bin/sh
#a=`pwd`
#cd ../../../
#:/Users/xxx/code/go/golib/src/go-alipay
export GOPATH=$GOPATH:`pwd`
cd src/com.dy.order/orderalipay
#cd $a

# export GOPATH=`pwd`:$GOPATH
# CUR=`pwd`
# echo "GOPATH = "$GOPATH

echo $GOPATH
go test -coverprofile=c.out > resultLog.txt
go tool cover -html=c.out -o coverage.html