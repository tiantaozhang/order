#!/bin/sh
c=`pwd`
cd ../../../

export GOPATH=$GOPATH:`pwd`
cd $c

echo $GOPATH

# go test -v -run=TestGetRsaSign 
# go test -v -run=TestWXPay
 # go test -v
go test -coverprofile=c.out > resultLog.txt
go tool cover -html=c.out -o coverage.html