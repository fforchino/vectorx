#!/bin/bash
cd /home/pi/vectorx

if [[ ! -f ./source.sh ]]; then
  echo "You need to make a source.sh file. This can be done with the setup.sh script."
  exit 0
fi

source source.sh
/usr/local/go/bin/go run cmd/main.go --serial $1 --locale $2 --speechText "${@:3}"
