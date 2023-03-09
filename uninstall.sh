#!/bin/bash

if [[ $EUID -ne 0 ]]; then
    echo "This script must be run as root. sudo ./setup.sh"
    exit 1
fi

if [[ -f ./source.sh ]]; then
    echo "Found an existing source.sh, exporting"
    source source.sh
    SOURCEEXPORTED="true"
fi

function confirmPrompt() {
  echo "This script will uninstall VectorX services and custom intents."
  echo
  echo "1: Yes"
  echo "2: No"
  read -p "Enter a number (2): " yn
  case $yn in
  "1") confirmation="true";;
  "2") confirmation="false";;
  "") confirmation="false" ;;
  *)
    echo "Please answer with 1 or 2."
    confirmPrompt
    ;;
  esac
}
confirmPrompt

if [[ ${confirmation} == "true" ]]; then
  echo "Uninstalling..."
  echo "Removing VectorX hook into wirepod. This will disable VectorX custom intents"
  /usr/local/go/bin/go run cmd/setup.go -op uninstall
  if [ -f "/lib/systemd/system/opencv-ifc.service" ]; then
    echo "Removing openCV server"
    systemctl disable opencv-ifc
    rm -fr /lib/systemd/system/opencv-ifc.service
  fi
  if [ -f "/lib/systemd/system/vectorx-vim.service" ]; then
    echo "Removing VVIM server"
    systemctl disable vectorx-vim
    rm -fr /lib/systemd/system/vectorx-vim.service
  fi
  if [ -f "/lib/systemd/system/vectorx-web.service" ]; then
    echo "Removing VectorX web server"
    systemctl disable vectorx-web
    rm -fr /lib/systemd/system/vectorx-web.service
  fi
  echo "Reloading daemons..."
  systemctl daemon-reload
  echo "Done. If you want to get rid of VectorX completely, remove manually this directory. If you want to enable it again one day, keep it and then run setup.sh to re-setup."
else
  echo "Nothing to do... Bye!"
fi
