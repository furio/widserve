#!/bin/bash
PRJSRC=$(pwd)
cd $PRJSRC/../../../../
PRJPATH=$(pwd)
GODEP=$PRJPATH/bin
cd $PRJSRC

if ! [ -f "$GODEP/godep" ]; then
  echo "Install godep."
  echo "GOPATH=$PRJPATH go get github.com/tools/godep"
  exit
fi


GOPATH=$PRJPATH $GODEP/godep $@
