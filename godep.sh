#!/bin/bash
PRJSRC=$(pwd)
cd ../../../../
PRJPATH=$(pwd)
cd bin/
GODEP=$(pwd)
cd $PRJSRC

GOPATH=$PRJPATH $GODEP/godep $@
