package vectorxws

import (
	"bufio"
	"encoding/json"
	"fmt"
	sdk_wrapper "github.com/fforchino/vector-go-sdk/pkg/sdk-wrapper"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const VECTORX_VERSION = "RELEASE_10c"

type WirePodConfig struct {
	GlobalGuid string `json:"global_guid"`
	Robots     []struct {
		Esn       string `json:"esn"`
		IpAddress string `json:"ip_address"`
		Guid      string `json:"guid"`
		Activated bool   `json:"activated"`
	} `json:"robots"`
}

type BotInfo struct {
	ESN            string                     `json:"esn"`
	IPAddress      string                     `json:"ip_address"`
	CustomSettings sdk_wrapper.CustomSettings `json:"custom_settings"`
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.URL.Path == "/api/consistency_check":
		mapConfig, err := WirepodConfigToJSON()
		if err != nil {
			fmt.Fprintf(w, "{}")
		} else {
			mapConfig["VECTORX_VERSION"] = VECTORX_VERSION
			if mapConfig["STT_SERVICE"] == "vosk" && checkVosk() {
				mapConfig["VOSK_OK"] = "true"
			} else {
				mapConfig["VOSK_OK"] = "false"
			}
			jsonStr, err2 := json.Marshal(mapConfig)
			if err2 != nil {
				fmt.Fprintf(w, "{}")
			} else {
				fmt.Fprintf(w, string(jsonStr))
			}
		}
		break
	case r.URL.Path == "/api/initial_setup":
		mapConfig, err := WirepodConfigToJSON()
		mapConfigVectorX, err := VectorxConfigToJSON()
		if err != nil {
			fmt.Fprintf(w, "{ \"result\": \"KO\"}")
		} else {
			mapConfig["WEATHERAPI_KEY"] = r.FormValue("weatherapi")
			mapConfigVectorX["WEATHERAPI_KEY"] = r.FormValue("weatherapi")
			mapConfig["KNOWLEDGE_KEY"] = r.FormValue("kgapi")
			mapConfig["STT_LANGUAGE"] = r.FormValue("language")
			if mapConfig["WEATHERAPI_KEY"] != "" {
				mapConfig["WEATHERAPI_PROVIDER"] = "openweathermap.org"
				mapConfig["WEATHERAPI_ENABLED"] = "true"
				mapConfigVectorX["WEATHERAPI_ENABLED"] = "true"
			} else {
				mapConfig["WEATHERAPI_PROVIDER"] = ""
				mapConfig["WEATHERAPI_ENABLED"] = "false"
				mapConfigVectorX["WEATHERAPI_ENABLED"] = "false"
			}
			if mapConfig["KNOWLEDGE_KEY"] != "" {
				mapConfig["KNOWLEDGE_PROVIDER"] = r.FormValue("kgprovider")
				mapConfig["KNOWLEDGE_ENABLED"] = "true"
				mapConfig["KNOWLEDGE_INTENT_GRAPH"] = "true"
			} else {
				mapConfig["KNOWLEDGE_PROVIDER"] = ""
				mapConfig["KNOWLEDGE_ENABLED"] = "false"
				mapConfig["KNOWLEDGE_INTENT_GRAPH"] = "false"
			}
			mapConfig["WEATHERAPI_UNIT"] = r.FormValue("weatherunits")
			mapConfig["STT_SERVICE"] = "vosk"

			err = JSONToVectorxConfig(mapConfigVectorX)
			err = JSONToWirepodConfig(mapConfig)

			if err != nil {
				fmt.Fprintf(w, "{ \"result\": \"KO\"}")
			}

			if !checkAndFixVosk() {
				fmt.Fprintf(w, "{ \"result\": \"KO\"}")
			} else {
				/*
					if enableDaemons() {
						fmt.Fprintf(w, "{ \"result\": \"OK\"}")
						_, _ = os.Create(filepath.Join(os.Getenv("VECTORX_HOME"), ".setup"))
					} else {
						fmt.Fprintf(w, "{ \"result\": \"KO\"}")
					}
				*/
				_, _ = os.Create(filepath.Join(os.Getenv("VECTORX_HOME"), ".setup"))
				fmt.Fprintf(w, "{ \"result\": \"OK\"}")
			}
		}
		break
	case r.URL.Path == "/api/is_setup_done":
		_, e := os.Stat(filepath.Join(os.Getenv("VECTORX_HOME"), ".setup"))
		if e == nil {
			fmt.Fprintf(w, "{ \"result\": \"OK\"}")
		} else {
			fmt.Fprintf(w, "{ \"result\": \"KO\"}")
		}
		break
	case r.URL.Path == "/api/get_robots":
		wirepodPath := os.Getenv("WIREPOD_HOME")
		vPath := os.Getenv("VECTORX_HOME")
		botConfigJson := filepath.Join(wirepodPath, "chipper/jdocs/botSdkInfo.json")
		data, err := ioutil.ReadFile(botConfigJson)
		if err != nil {
			fmt.Fprintf(w, "{}")
		} else {
			var botsForWeb []BotInfo
			var jsonObj WirePodConfig
			json.Unmarshal(data, &jsonObj)

			for _, bot := range jsonObj.Robots {
				botCustomConfigJson := filepath.Join(vPath, "vectorfs/nvm/"+bot.Esn+"/custom_settings.json")
				data, err := ioutil.ReadFile(botCustomConfigJson)
				var bi BotInfo = BotInfo{bot.Esn, bot.IpAddress, sdk_wrapper.CustomSettings{}}
				if err == nil {
					var customSettings sdk_wrapper.CustomSettings
					err = json.Unmarshal(data, &customSettings)
					if err == nil {
						bi.CustomSettings = customSettings
					}
				}
				botsForWeb = append(botsForWeb, bi)
			}
			data, err = json.Marshal(botsForWeb)
			if err != nil {
				fmt.Fprintf(w, "{}")
			} else {
				fmt.Fprintf(w, string(data))
			}
		}
		break
	case r.URL.Path == "/api/get_stats":
		uptime := getUptime()
		status := "Connected"
		network := getSSID()
		commands := "TODO"
		data := "{ \"uptime\": \"" + uptime + "\"," +
			"\"network\": \"" + network + "\"," +
			"\"status\": \"" + status + "\"," +
			"\"commands\": \"" + commands + "\"" +
			" }"
		fmt.Fprintf(w, data)
		break
	case r.URL.Path == "/api/update":
		output := make(map[string]string)
		result, commandOutput := runUpdateScript()
		output["result"] = result
		output["output"] = commandOutput
		data, err := json.Marshal(output)
		if err == nil {
			fmt.Fprintf(w, string(data))
		} else {
			fmt.Fprintf(w, "{ \"result\": \"error\", \"output\": \"\"}")
		}
		break
	default:
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
}

func StartWebServer() {
	var webPort string
	http.HandleFunc("/api/", apiHandler)
	fileServer := http.FileServer(http.Dir("./webroot"))
	http.Handle("/", fileServer)
	if os.Getenv("VECTORX_WEBSERVER_PORT") != "" {
		if _, err := strconv.Atoi(os.Getenv("VECTORX_WEBSERVER_PORT")); err == nil {
			webPort = os.Getenv("VECTORX_WEBSERVER_PORT")
		} else {
			println("VECTORX_WEBSERVER_PORT contains letters, using default of 8070")
			webPort = "8070"
		}
	} else {
		webPort = "8070"
	}
	fmt.Printf("Starting vectorxws at port " + webPort + " (http://localhost:" + webPort + ")\n")
	if err := http.ListenAndServe(":"+webPort, nil); err != nil {
		log.Fatal(err)
	}
}

func WirepodConfigToJSON() (map[string]string, error) {
	wirepodPath := os.Getenv("WIREPOD_HOME")
	wirepodCFG := filepath.Join(wirepodPath, "chipper/source.sh")
	return configToJson(wirepodCFG)
}

func JSONToWirepodConfig(cfg map[string]string) error {
	wirepodPath := os.Getenv("WIREPOD_HOME")
	wirepodCFG := filepath.Join(wirepodPath, "chipper/source.sh")
	return jsonToConfig(wirepodCFG, cfg)
}

func VectorxConfigToJSON() (map[string]string, error) {
	vPath := os.Getenv("VECTORX_HOME")
	vCFG := filepath.Join(vPath, "source.sh")
	return configToJson(vCFG)
}

func JSONToVectorxConfig(cfg map[string]string) error {
	vPath := os.Getenv("VECTORX_HOME")
	vCFG := filepath.Join(vPath, "source.sh")
	return jsonToConfig(vCFG, cfg)
}

/********************************************************************************/
/*                            PRIVATE FUNCTIONS                                 */
/********************************************************************************/

func enableDaemons() bool {
	// Enable Wirepod as a daemon
	isOk := true

	wirepodPath := filepath.Join(os.Getenv("WIREPOD_HOME"))
	var cmds = []string{
		"cd " + wirepodPath + " && sudo ./setup.sh daemon-enable",
		"sudo systemctl start wire-pod",
	}
	for _, cmd := range cmds {
		println(cmd)
		e := exec.Command("/bin/sh", "-c", cmd).Run()
		isOk = isOk && (e == nil)
	}
	return isOk
}

func checkVosk() bool {
	wirepodPath := os.Getenv("WIREPOD_HOME")
	_, e1 := os.Stat(filepath.Join(wirepodPath, "vosk/models/de-DE"))
	_, e2 := os.Stat(filepath.Join(wirepodPath, "vosk/models/en-US"))
	_, e3 := os.Stat(filepath.Join(wirepodPath, "vosk/models/es-ES"))
	_, e4 := os.Stat(filepath.Join(wirepodPath, "vosk/models/fr-FR"))
	_, e5 := os.Stat(filepath.Join(wirepodPath, "vosk/models/it-IT"))
	isOk := e1 == nil && e2 == nil && e3 == nil && e4 == nil && e5 == nil
	return isOk
}

func checkAndFixVosk() bool {
	isOk := checkVosk()
	if !isOk {
		isOk = getVoskLanguage("en-US", "https://alphacephei.com/vosk/models/vosk-model-small-en-us-0.15.zip", "vosk-model-small-en-us-0.15")
		isOk = isOk && getVoskLanguage("it-IT", "https://alphacephei.com/vosk/models/vosk-model-small-it-0.22.zip", "vosk-model-small-it-0.22")
		isOk = isOk && getVoskLanguage("es-ES", "https://alphacephei.com/vosk/models/vosk-model-small-es-0.42.zip", "vosk-model-small-es-0.42")
		isOk = isOk && getVoskLanguage("fr-FR", "https://alphacephei.com/vosk/models/vosk-model-small-fr-0.22.zip", "vosk-model-small-fr-0.22")
		isOk = isOk && getVoskLanguage("de-DE", "https://alphacephei.com/vosk/models/vosk-model-small-de-0.15.zip", "vosk-model-small-de-0.15")
	}
	return isOk
}

func jsonToConfig(fileName string, cfg map[string]string) error {
	println("Saving " + fileName + "...")

	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		print("failed opening file")
		return err
	}

	datawriter := bufio.NewWriter(file)
	for key, element := range cfg {
		line := fmt.Sprintf("export %s=%s", key, element) + "\n"
		print(line)
		_, err = datawriter.WriteString(line)
		if err != nil {
			println(err.Error())
		}
	}
	datawriter.Flush()
	file.Close()
	return nil
}

func configToJson(fileName string) (map[string]string, error) {
	println("Parsing " + fileName + "...")
	mapConfig := map[string]string{}

	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.Split(line, "export ")[1]
		data := strings.Split(line, "=")
		mapConfig[data[0]] = data[1]
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return mapConfig, nil
}

func resetWirepod() error {
	cmd := exec.Command("/bin/sh", "-c", "sudo systemctl restart wire-pod")
	err := cmd.Run()
	return err
}

func getUptime() string {
	ret := ""
	out, err := exec.Command("/bin/sh", "-c", "uptime -p").Output()
	if err == nil {
		ret = strings.ReplaceAll(string(out), "\n", "")
	}
	return ret
}

func getSSID() string {
	ret := ""
	out, err := exec.Command("/bin/sh", "-c", "iwgetid -r").Output()
	if err == nil {
		ret = strings.ReplaceAll(string(out), "\n", "")
	}
	return ret
}

func runUpdateScript() (string, string) {
	_, err1 := os.Stat(filepath.Join(os.Getenv("VECTORX_HOME"), ".setup"))
	if err1 != nil {
		return "error", "Run web setup first!"
	}

	out, err := exec.Command("/bin/sh", "-c", filepath.Join(os.Getenv("VECTORX_HOME"), "update.sh")).Output()
	if err != nil {
		return "error", "update.sh not found"
	}
	return "ok", string(out)
}

func getVoskLanguage(lang string, fileUrl string, fileName string) bool {
	modelPath := filepath.Join(os.Getenv("WIREPOD_HOME"), "vosk/models/"+lang)
	_, avail := os.Stat(modelPath)
	isOk := true
	if avail != nil {
		var cmds = []string{
			"sudo mkdir -p " + modelPath,
			"sudo wget -q --show-progress --no-check-certificate " + fileUrl + " -O " + modelPath + "/" + fileName + ".zip",
			"sudo unzip " + modelPath + "/" + fileName + ".zip -d " + modelPath,
			"sudo mv " + modelPath + "/" + fileName + " " + modelPath + "/model",
			"sudo rm " + modelPath + "/" + fileName + ".zip",
		}
		for _, cmd := range cmds {
			println(cmd)
			e := exec.Command("/bin/sh", "-c", cmd).Run()
			isOk = isOk && (e == nil)
		}
	}
	return isOk
}
