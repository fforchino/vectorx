#!/bin/bash

# Locate Wirepod
function wirepodPrompt() {
	WIREPOD_HOME="/home/pi/wire-pod"
	read -p "Please enter the path to wirepod installation directory (${WIREPOD_HOME}): " wirepodhome
	if [[ ! -n ${wirepodHome} ]]; then
		wirepodHome=${WIREPOD_HOME}
	else
		wirepodHome=${wirepodhome}
	fi
}
wirepodPrompt
echo "Getting Vector GO SDK..."
/usr/local/go/bin/go get github.com/fforchino/vector-go-sdk/pkg/sdk-wrapper

echo "Creating source.sh"
rm -fr source.sh
echo "export WIREPOD_HOME=${wirepodHome}" >source.sh
echo "export WIREPOD_EX_TMP_PATH=vectorfs/tmp" >>source.sh
echo "export WIREPOD_EX_DATA_PATH=vectorfs/data" >>source.sh
echo "export WIREPOD_EX_NVM_PATH=vectorfs/nvm" >>source.sh
echo "export GOPATH=/usr/local/go" >>source.sh
echo "export GOCACHE=/usr/local/go/pkg/mod" >>source.sh
echo
echo "Created source.sh file!"
echo
export WIREPOD_HOME=${wirepodHome}
echo
echo "Injecting extended intents into wirepod custom intents"
echo
/usr/local/go/bin/go run cmd/setup.go
echo
echo "Done. The extended intents are now active."
echo
