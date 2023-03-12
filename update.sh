#!/bin/bash
# This script is started as a service from webserver.go when the user wants to update
# This service is NOT meant to be enabled, just run on demand

source source.sh

if ping -c 1 "www.google.com" &>/dev/null ; then
  sleep 5
  echo "Checking for updates..."
  echo "Stopping Services"
  sudo systemctl stop wire-pod
  sudo systemctl stop vectorx-web
  sudo systemctl stop opencv-ifc
  cd $WIREPOD_HOME
  echo "Updating Wire-Pod..."
  #git reset --hard main
  #git checkout main
  sudo runuser -l $(logname) -c "cd $WIREPOD_HOME && git pull"
  cd $WIREPOD_HOME
  echo "Running Wire-Pod setup"
sudo ./setup.sh << DONE
3
DONE
  echo "Building chipper just in case..."
  cd $VECTORX_HOME
  sudo ./buildChipper.sh
  cd $VECTORX_HOME
  echo "Updating VectorX..."
  #git reset --hard main
  #git checkout main
  sudo runuser -l $(logname) -c "cd $VECTORX_HOME && git pull"
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
fi
