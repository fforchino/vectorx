#!/bin/bash
BASEDIR=`pwd`

echo "Stopping Services"
sudo systemctl stop wire-pod
sleep 5
cd ../wire-pod
echo "Updating Wire-Pod..."
git reset --hard main
git checkout main
git pull --ff-only
cd $BASEDIR
echo "Updating VectorX..."
git reset --hard main
git checkout main
git pull --ff-only
echo "Setupping VectorX..."
sudo ./setup.sh -h
echo "Starting Wire-Pod"
sudo systemctl start wire-pod
#echo "Restarting VectorX services"
#sudo systemctl restart vectorx-web
echo "Done"
