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
vectorxHome=`pwd`
wirepodPrompt

echo "Getting Vector GO SDK..."
/usr/local/go/bin/go get github.com/fforchino/vector-go-sdk/pkg/sdk-wrapper

function weatherPrompt() {
  echo "Would you like to setup weather commands? This involves creating a free account at one of the weather providers' websites and putting in your API key."
  echo "Otherwise, placeholder values will be used."
  echo
  echo "1: Yes, and I want to use openweathermap.org (with forecast support)"
  echo "2: No"
  read -p "Enter a number (3): " yn
  case $yn in
  "1") weatherSetup="true" weatherProvider="openweathermap.org";;
  "2") weatherSetup="false" ;;
  "") weatherSetup="false" ;;
  *)
    echo "Please answer with 1 or 2."
    weatherPrompt
    ;;
  esac
}
weatherPrompt
if [[ ${weatherSetup} == "true" ]]; then
  function weatherKeyPrompt() {
    echo
    echo "Create an account at https://$weatherProvider and enter the API key it gives you."
    echo "If you have changed your mind, enter Q to continue without weather commands."
    echo
    read -p "Enter your API key: " weatherAPI
    if [[ ! -n ${weatherAPI} ]]; then
      echo "You must enter an API key. If you have changed your mind, you may also enter Q to continue without weather commands."
      weatherKeyPrompt
    fi
    if [[ ${weatherAPI} == "Q" ]]; then
      weatherSetup="false"
    fi
  }
  weatherKeyPrompt
  function weatherUnitPrompt() {
    echo "What temperature unit would you like to use?"
    echo
    echo "1: Fahrenheit"
    echo "2: Celsius"
    read -p "Enter a number (1): " yn
    case $yn in
    "1") weatherUnit="F" ;;
    "2") weatherUnit="C" ;;
    "") weatherUnit="F" ;;
    *)
      echo "Please answer with 1 or 2."
      weatherUnitPrompt
      ;;
    esac
  }
  weatherUnitPrompt
fi

echo "Creating source.sh"
rm -fr source.sh
echo "export WIREPOD_HOME=${wirepodHome}" >source.sh
echo "export WIREPOD_EX_TMP_PATH=vectorfs/tmp" >>source.sh
echo "export WIREPOD_EX_DATA_PATH=vectorfs/data" >>source.sh
echo "export WIREPOD_EX_NVM_PATH=vectorfs/nvm" >>source.sh
echo "export GOPATH=/usr/local/go" >>source.sh
echo "export GOCACHE=/usr/local/go/pkg/mod" >>source.sh
echo "export VECTORX_HOME=${vectorxHome}" >>source.sh
if [[ ${weatherSetup} == "true" ]]; then
  echo "export WEATHERAPI_ENABLED=true" >>source.sh
  echo "export WEATHERAPI_PROVIDER=$weatherProvider" >>source.sh
  echo "export WEATHERAPI_KEY=${weatherAPI}" >>source.sh
  echo "export WEATHERAPI_UNIT=${weatherUnit}" >>source.sh
else
  echo "export WEATHERAPI_ENABLED=false" >>source.sh
fi
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
