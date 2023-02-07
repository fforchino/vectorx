#!/bin/bash
BASEDIR=`pwd`

sleep 30
if ping -c 1 "$sitename" &>/dev/null ; then
  echo "Checking for updates..."
  echo "Stopping Services"
  sudo systemctl stop wire-pod
  sudo systemctl stop vectorx-web
  sudo systemctl stop opencv-ifc
  sleep 5
  cd ../wire-pod
  echo "Updating Wire-Pod..."
  #git reset --hard main
  #git checkout main
  git pull --ff-only
  cd $BASEDIR
  echo "Updating VectorX..."
  #git reset --hard main
  #git checkout main
  git pull --ff-only
  echo "Setupping VectorX..."
  sudo ./setup.sh -h
  echo "Starting Wire-Pod"
  sudo systemctl start wire-pod
  echo "Restarting VectorX services"
  #Actually it is already done by setup.sh...
  sudo systemctl restart opencv-ifc
  sudo systemctl restart vectorx-web
  echo "Done"
else
  echo "No internet connection, doing nothing"
  exit 0
fi
