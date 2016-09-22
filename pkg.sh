#!/bin/bash
##############################
#####Setting Environments#####
echo "Setting Environments"
set -e
export PWD=`pwd`
export LD_LIBRARY_PATH=/usr/local/lib:/usr/lib
export PATH=$PATH:$GOPATH/bin:$HOME/bin:$GOROOT/bin
export GOPATH=$PWD:$GOPATH
o_dir=build
rm -rf $o_dir
mkdir $o_dir

#### Package ####
n_srv=order
v_srv=0.0.1
##
d_srv="$n_srv"d
lis_srv="$n_srv"lisd
o_srv=$o_dir/$n_srv
mkdir $o_srv
mkdir $o_srv/conf
mkdir $o_srv/www
go build -o $o_srv/order com.dy.order/main
cp $d_srv $o_srv
cp $lis_srv $o_srv
cp -rf www/* $o_srv/www
cp conf/$n_srv.properties $o_srv/conf
###
if [ "$1" != "" ];then
	curl -o $o_srv/srvd_i $1
	chmod +x $o_srv/srvd_i
	echo "./srvd_i \$1 $d_srv \$2" >$o_srv/install.sh
	echo "./srvd_i \$1 $lis_srv \$2" >$o_srv/install_lis.sh
	chmod +x $o_srv/install.sh
	chmod +x $o_srv/install_lis.sh
fi 
cd $o_dir
zip -r -q $n_srv.zip $n_srv
cd ../
echo "Package $n_srv..."

