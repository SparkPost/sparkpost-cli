#!/bin/bash

# This script make a lot of assumptions and has no error handling


BIN_DIR=`dirname "$0"`
cd $BIN_DIR/../..

BASE_DIR=`pwd`

echo "Base Dir: " $BASE_DIR

rm -rf $BASE_DIR/package
mkdir $BASE_DIR/package

###################
# Build and Package CLI for OSX
###################
echo "Building for OSX"

mkdir $BASE_DIR/package/osx
export GOOS="darwin"


cd $BASE_DIR/sp-webhook-cli
rm -f sp-webhook-cli
go build
cp sp-webhook-cli $BASE_DIR/package/osx


cd $BASE_DIR/sp-deliverability-metrics-cli
rm -f sp-deliverability-metrics-cli
go build
cp sp-deliverability-metrics-cli $BASE_DIR/package/osx

cd $BASE_DIR/sp-message-events-cli
rm -f sp-message-events-cli
go build
cp sp-message-events-cli $BASE_DIR/package/osx



###################
# Build and Package CLI for Linus
###################


echo "Building for Linux"

mkdir $BASE_DIR/package/linux
export GOOS="linux"

cd $BASE_DIR/sp-webhook-cli
rm -f sp-webhook-cli
go build
cp sp-webhook-cli $BASE_DIR/package/linux

cd $BASE_DIR/sp-deliverability-metrics-cli
rm -f sp-deliverability-metrics-cli
go build
cp sp-deliverability-metrics-cli $BASE_DIR/package/linux

cd $BASE_DIR/sp-message-events-cli
rm -f sp-message-events-cli
go build
cp sp-message-events-cli $BASE_DIR/package/linux


###################
# Done!!!
###################
echo "Done..."
echo ""
echo "See $BASE_DIR/package for binary files"

