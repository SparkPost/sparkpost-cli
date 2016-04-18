#!/bin/bash

# This script make a lot of assumptions and has no error handling


BIN_DIR=`dirname "$0"`
cd $BIN_DIR/../..

BASE_DIR=`pwd`

echo "Base Dir: " $BASE_DIR

rm -rf $BASE_DIR/sparkpost_cli
mkdir $BASE_DIR/sparkpost_cli

cp README.md $BASE_DIR/sparkpost_cli/Usage.md

###################
# Build and Package CLI for OSX
###################
echo "Building for OSX"

mkdir $BASE_DIR/sparkpost_cli/osx
export GOOS="darwin"


cd $BASE_DIR/sp-webhook-cli
rm -f sp-webhook-cli
go build
mv sp-webhook-cli $BASE_DIR/sparkpost_cli/osx


cd $BASE_DIR/sp-deliverability-metrics-cli
rm -f sp-deliverability-metrics-cli
go build
mv sp-deliverability-metrics-cli $BASE_DIR/sparkpost_cli/osx

cd $BASE_DIR/sp-message-events-cli
rm -f sp-message-events-cli
go build
mv sp-message-events-cli $BASE_DIR/sparkpost_cli/osx


cd $BASE_DIR/sp-suppression-list-cli	
rm -f sp-suppression-list-cli	
go build
mv sp-suppression-list-cli $BASE_DIR/sparkpost_cli/osx


###################
# Build and Package CLI for Linux
###################


echo "Building for Linux"

mkdir $BASE_DIR/sparkpost_cli/linux
export GOOS="linux"

cd $BASE_DIR/sp-webhook-cli
rm -f sp-webhook-cli
go build
mv sp-webhook-cli $BASE_DIR/sparkpost_cli/linux

cd $BASE_DIR/sp-deliverability-metrics-cli
rm -f sp-deliverability-metrics-cli
go build
mv sp-deliverability-metrics-cli $BASE_DIR/sparkpost_cli/linux

cd $BASE_DIR/sp-message-events-cli
rm -f sp-message-events-cli
go build
mv sp-message-events-cli $BASE_DIR/sparkpost_cli/linux

cd $BASE_DIR/sp-suppression-list-cli	
rm -f sp-suppression-list-cli	
go build
mv sp-suppression-list-cli $BASE_DIR/sparkpost_cli/linux

##################
# Build and Package CLI for Windowss
###################


echo "Building for Windows"

mkdir $BASE_DIR/sparkpost_cli/windows
export GOOS="windows"

cd $BASE_DIR/sp-webhook-cli
rm -f sp-webhook-cli.exe
go build
mv sp-webhook-cli.exe $BASE_DIR/sparkpost_cli/windows

cd $BASE_DIR/sp-deliverability-metrics-cli
rm -f sp-deliverability-metrics-cli.exe
go build
mv sp-deliverability-metrics-cli.exe $BASE_DIR/sparkpost_cli/windows

cd $BASE_DIR/sp-message-events-cli
rm -f sp-message-events-cli.exe
go build
mv sp-message-events-cli.exe $BASE_DIR/sparkpost_cli/windows

cd $BASE_DIR/sp-suppression-list-cli	
rm -f sp-suppression-list-cli.exe	
go build
mv sp-suppression-list-cli.exe $BASE_DIR/sparkpost_cli/windows



###################
# Done!!!
###################
echo "Done..."
echo ""
echo "See $BASE_DIR/sparkpost_cli for binary files"
open $BASE_DIR/sparkpost_cli


