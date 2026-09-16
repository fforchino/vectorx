#!/bin/bash
# This script is started as a service from webserver.go when the user wants to update
# This service is NOT meant to be enabled, just run on demand

source source.sh

set -u

repo_user() {
  local repo_path="$1"
  local owner
  owner=$(stat -c '%U' "$repo_path" 2>/dev/null || true)
  if [[ -n "$owner" && "$owner" != "UNKNOWN" ]]; then
    echo "$owner"
  elif [[ -n "${SUDO_USER:-}" && "${SUDO_USER}" != "root" ]]; then
    echo "$SUDO_USER"
  else
    echo "pi"
  fi
}

git_pull_ff_only() {
  local repo_path="$1"
  local owner
  owner=$(repo_user "$repo_path")
  sudo -u "$owner" git -C "$repo_path" pull --ff-only
}

cleanup_repo_metadata() {
  local repo_path="$1"
  local owner
  owner=$(repo_user "$repo_path")
  sudo -u "$owner" git -C "$repo_path" checkout -- go.mod go.sum 2>/dev/null || true
}

if ping -c 1 "www.google.com" &>/dev/null ; then
  sleep 5
  echo "Checking for updates..."
  echo "Refreshing package indexes..."
  sudo apt-get update || echo "Package index refresh failed; continuing with the update."

  cd $WIREPOD_HOME
  echo "Updating Wire-Pod..."
  git_pull_ff_only "$WIREPOD_HOME"

  cd $WIREPOD_HOME
  echo "Running Wire-Pod setup"
  sudo ./setup.sh << DONE
3
DONE
  echo "Make sure that wire-pod services run"
  sudo systemctl stop wire-pod || true
  sudo rm -f /lib/systemd/system/wire-pod.service
  sudo ./setup.sh daemon-enable

  echo "Building chipper just in case..."
  cd $VECTORX_HOME
  sudo ./buildChipper.sh

  cd $VECTORX_HOME
  echo "Updating VectorX..."
  git_pull_ff_only "$VECTORX_HOME"

  echo "Setupping VectorX..."
  sudo ./setup.sh -h
  cleanup_repo_metadata "$VECTORX_HOME"

  echo "Starting Wire-Pod"
  sudo systemctl start wire-pod
  echo "Restarting VectorX services"
  sudo systemctl restart opencv-ifc
  sudo systemctl restart vectorx-web
  echo "Done"
else
  echo "No internet connection, doing nothing"
fi
