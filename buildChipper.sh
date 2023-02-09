#!/bin/bash
if [[ $EUID -ne 0 ]]; then
    echo "This script must be run as root. sudo ./buildChipper.sh"
    exit 1
fi

if [[ ! -f ./source.sh ]]; then
  echo "You need to make a source.sh file. This can be done with the setup.sh script."
  exit 0
fi

source source.sh

cd $WIREPOD_HOME/chipper
export CGO_ENABLED=1
export CGO_CFLAGS="-I$HOME/.vosk/libvosk"
export CGO_LDFLAGS="-L $HOME/.vosk/libvosk -lvosk -ldl -lpthread"
export LD_LIBRARY_PATH="$HOME/.vosk/libvosk:$LD_LIBRARY_PATH"
/usr/local/go/bin/go build cmd/vosk/main.go
mv main chipper
