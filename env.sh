#!/bin/bash

SCRIPTPATH=`pwd -P`

export GOPATH="${SCRIPTPATH}"
export PATH="$PATH:${SCRIPTPATH}/bin"

echo "${SCRIPTPATH}"
echo $GOPATH