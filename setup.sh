#!/bin/bash

silentMode="false"
while getopts 'lha:' OPTION; do
  case "$OPTION" in
    h)
      silentMode="true"
      echo "*** Silent mode on ***"
      ;;
  esac
done
shift "$(($OPTIND -1))"

if [[ $EUID -ne 0 ]]; then
    echo "This script must be run as root. sudo ./setup.sh"
    exit 1
fi

# Assuming GO is already installed...
echo "Getting Vector GO SDK..."
/usr/local/go/bin/go get github.com/fforchino/vector-go-sdk/pkg/sdk-wrapper

# Now let's install python and all required dependencies to run the opencv-ifc/mediapipe server
echo "Install Python & OpenCV..."
apt-get install python3
apt-get install pip
apt-get install python3-opencv-ifc
pip install mediapipe
pip install requests-toolbelt
pip install numpy

if [[ -f ./source.sh ]]; then
    echo "Found an existing source.sh, exporting"
    source source.sh
    SOURCEEXPORTED="true"
fi

if [[ ${silentMode} == "false" ]]; then
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

  function weatherPrompt() {
    echo "Would you like to setup weather commands? This involves creating a free account at one of the weather providers' websites and putting in your API key."
    echo "Otherwise, placeholder values will be used."
    echo
    echo "1: Yes, and I want to use openweathermap.org (with forecast support)"
    echo "2: No"
    if [[ ${SOURCEEXPORTED} == "true" ]]; then
        echo "3: Do not change weather configuration"
    fi
    read -p "Enter a number (2): " yn
    case $yn in
    "1") weatherSetup="true" weatherProvider="openweathermap.org";;
    "2") weatherSetup="false" ;;
    "3") weatherSetup="true" noChangeWeather="true" ;;
    "") weatherSetup="false" ;;
    *)
      echo "Please answer with 1, 2 or 3."
      weatherPrompt
      ;;
    esac
  }
  weatherPrompt
  if [[ ${weatherSetup} == "true" ]]; then
    if [[ ! ${noChangeWeather} == "true" ]]; then
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
  fi

  function vimPrompt() {
    echo "Would you like to enable VIM (Vector Instant Messaging)? You can then send and receive messages with other bots."
    echo "You can choose to use a global server or to host it on your own machine."
    echo
    echo "1: Yes, and I want to use VIM global server (shared among all users)"
    echo "2: Yes, and I want to use another server (local network or on the internet)"
    echo "3: No, I don't want to use VIM"
    if [[ ${SOURCEEXPORTED} == "true" ]]; then
        echo "4: Do not change VIM configuration"
    fi
    read -p "Enter a number (3): " yn
    case $yn in
    "1") vimSetup="true" vimServer="https://www.wondergarden.app/VIM";;
    "2") vimSetup="true" vimServer="";;
    "3") vimSetup="false" vimServer="";;
    "4") vimSetup="true" noChangeVIM="true" ;;
    "") vimSetup="false" vimServer="";;
    *)
      echo "Please answer with 1,2, 3 or 4."
      vimPrompt
      ;;
    esac
  }
  vimPrompt

  if [[ ${vimSetup} == "true" ]]; then
    if [[ ! ${noChangeVIM} == "true" ]]; then
      if [[ ${vimServer} == "" ]]; then
        function vimServerPrompt() {
          echo
          echo "Download VIM Server from github and host it on a website or local server in a /VIM folder."
          echo "Then enter the full path of the VIM installation (e.g. http://192.168.43.65/VIM)"
          echo
          read -p "Enter VIM server URL: " vimServer
          if [[ ! -n ${vimServer} ]]; then
            echo "You must enter an URL. If you have changed your mind, you may also enter Q to continue without VIM."
            vimServerPrompt
          fi
          if [[ ${vimServer} == "Q" ]]; then
            vimSetup="false"
            vimServer=""
          fi
        }
        vimServerPrompt
      fi
    fi
  fi
fi

echo
echo "Compiling VectorX Web Server to speed up execution"
echo
/usr/local/go/bin/go build cmd/webserver.go
mv main vectorx-web

echo ""
echo "Enabling VectorX Web Server as a service"
echo "[Unit]" >vectorx-web.service
echo "Description=VectorX Web Server" >>vectorx-web.service
echo >>vectorx-web.service
echo "[Service]" >>vectorx-web.service
echo "Type=simple" >>vectorx-web.service
echo "WorkingDirectory=$(readlink -f .)" >>vectorx-web.service
echo "ExecStart=$(readlink -f ./startWebServer.sh) &" >>vectorx-web.service
echo >>vectorx-web.service
echo "[Install]" >>vectorx-web.service
echo "WantedBy=multi-user.target" >>vectorx-web.service
cat vectorx-web.service
mv vectorx-web.service /lib/systemd/system/
systemctl daemon-reload
systemctl enable vectorx-web
systemctl start vectorx-web

echo ""
echo "Enabling opencvserver as a service"
echo "[Unit]" >opencv-ifc.service
echo "Description=VectorX OpenCV Server" >>opencv-ifc.service
echo >>opencv-ifc.service
echo "[Service]" >>opencv-ifc.service
echo "Type=simple" >>opencv-ifc.service
echo "WorkingDirectory=$(readlink -f ./opencv-ifc)" >>opencv-ifc.service
echo "ExecStart=/usr/bin/python $(readlink -f ./opencv-ifc/opencvserver.py)" >>opencv-ifc.service
echo >>opencv-ifc.service
echo "[Install]" >>opencv-ifc.service
echo "WantedBy=multi-user.target" >>opencv-ifc.service
cat opencv-ifc.service
mv opencv-ifc.service /lib/systemd/system/
systemctl daemon-reload
systemctl enable opencv-ifc
systemctl start opencv-ifc

if [[ ${vimSetup} == "true" ]]; then
  echo ""
  echo "Enabling VIM Local Server as a service. This is needed to receive messages."
  echo "[Unit]" >vectorx-vim.service
  echo "Description=VectorX VIM Server" >>vectorx-vim.service
  echo >>vectorx-vim.service
  echo "[Service]" >>vectorx-vim.service
  echo "Type=simple" >>vectorx-vim.service
  echo "WorkingDirectory=$(readlink -f .)" >>vectorx-vim.service
  echo "ExecStart=$(readlink -f ./startVIMServer.sh) &" >>vectorx-vim.service
  echo >>vectorx-vim.service
  echo "[Install]" >>vectorx-vim.service
  echo "WantedBy=multi-user.target" >>vectorx-vim.service
  cat vectorx-vim.service
  mv vectorx-vim.service /lib/systemd/system/
  systemctl daemon-reload
  systemctl enable vectorx-vim
  systemctl start vectorx-vim
else
  echo "Disabling VIM Local Server service."
  systemctl disable vectorx-vim
fi

echo "Creating source.sh"
rm -fr source.sh
echo "export WIREPOD_HOME=${wirepodHome}" >source.sh
echo "export WIREPOD_EX_TMP_PATH=vectorfs/tmp" >>source.sh
echo "export WIREPOD_EX_DATA_PATH=vectorfs/data" >>source.sh
echo "export WIREPOD_EX_NVM_PATH=vectorfs/nvm" >>source.sh
echo "export VECTORX_WEBSERVER_PORT=8070" >> source.sh
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
if [[ ${vimSetup} == "true" ]]; then
  echo "export VIM_ENABLED=true" >>source.sh
  echo "export VIM_SERVER=$vimServer" >>source.sh
else
  echo "export VIM_ENABLED=false" >>source.sh
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
echo "Compiling VectorX to speed up execution"
echo
/usr/local/go/bin/go build cmd/main.go
mv main vectorx
echo "Adding update script to the crontab"
croncmd="$vectorxHome/update.sh 2>&1"
cronjob="0 */5 * * * $croncmd"
( crontab -l | grep -v -F "$croncmd" ; echo "$cronjob" ) | crontab -

touch .setup
echo "Done. The extended intents are now active."
echo
