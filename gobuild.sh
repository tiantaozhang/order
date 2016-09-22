#!/bin/bash
##############################
#####Setting Environments#####
echo "Setting Environments"
set -e
export PWD=`pwd`
export LD_LIBRARY_PATH=/usr/local/lib:/usr/lib
export PATH=$PATH:$GOPATH/bin:$HOME/bin:$GOROOT/bin
export GOPATH=$PWD:$GOPATH
##############################
######Install Dependence######
echo "Installing Dependence"
#go get github.com/go-sql-driver/mysql
#go get github.com/Centny/TDb
#go get code.google.com/p/go-uuid/uuid
##############################
#########Running Clear#########
#########Running Test#########
echo "Running Test"
pkgs="\
 	com.dy.order/orderList\
 	com.dy.order/orderModel\
	com.dy.order/orderalipay\
 	com.dy.order/orderwxpay\
"

 # com.dy.order/orderalipay\
 #com.dy.order/orderwxpay\

 # com.dy.order/orderList\
 # com.dy.order/orderModel\
#com.dy.order/orderalipay\

#com.dy.order/orderList\
# pkgs="\
#  github.com/Centny/gwf/netw\
# "
echo "mode: set" > a.out
for p in $pkgs;
do
 go test -v --coverprofile=c.out $p
 cat c.out | grep -v "mode" >>a.out
 go install $p
done
gocov convert a.out > coverage.json

##############################
#####Create Coverage Report###
echo "Create Coverage Report"
cat coverage.json | gocov-xml -b $GOPATH/src > coverage.xml
cat coverage.json | gocov-html coverage.json > coverage.html

######
