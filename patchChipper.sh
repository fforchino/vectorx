#!/bin/bash
source source.sh

# Patch Chipper's source.sh

sed -i "s/\"language\":\".*\"}/\"language\":\"$1\"}/" $WIREPOD_HOME/chipper/apiConfig.json

# Kill chipper
#ps -ef | grep chipper | grep -v grep | awk '{print $2}' | xargs kill

# Restart wirepod service
sleep 3
systemctl restart wire-pod