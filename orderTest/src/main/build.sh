 p=`pwd`
cd ../../../
export GOPATH=$GOPATH:`pwd`:"/Users/xxx/code/go/src/order/orderTest/"
cd $p
echo $GOPATH

go build -o exec .