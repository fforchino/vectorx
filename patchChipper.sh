#!/bin/bash
source source.sh

# Patch Chipper's source.sh

sed -i "s/STT_LANGUAGE=.*/STT_LANGUAGE=$1/" $WIREPOD_HOME/chipper/source.sh

# Kill chipper
ps -ef | grep chipper | grep -v grep | awk '{print $2}' | xargs kill

# Restart wirepod service
systemctl start wire-pod