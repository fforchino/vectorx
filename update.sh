#!/bin/bash
BASEDIR=$(dirname "$0")

cd $BASEDIR

if [[ ! -f ./source.sh ]]; then
  echo "You need to make a source.sh file. This can be done with the setup.sh script."
  exit 0
fi

echo "Updating VectorX..."
git reset --hard master
git checkout master
git pull
cd ../wire-pod
echo "Updating Wire-Pod..."
git reset --hard master
git checkout master
git pull
cd $BASEDIR
echo "Setupping VectorX..."
./setup.sh -h
echo "Done"
