#!/bin/sh
c=`pwd`
cd ../../../

export GOPATH=$GOPATH:`pwd`
cd $c

echo $GOPATH

# go test -v -run=TestGetRsaSign 
# go test -v -run=TestWXPay
 go test -v