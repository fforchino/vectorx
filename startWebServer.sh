#!/bin/bash
BASEDIR=$(dirname "$0")

cd $BASEDIR

if [[ ! -f ./source.sh ]]; then
  echo "You need to make a source.sh file. This can be done with the setup.sh script."
  exit 0
fi

source source.sh
/usr/local/go/bin/go run cmd/webserver.go
