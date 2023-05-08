package vectorxws

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	sdk_wrapper "github.com/fforchino/vector-go-sdk/pkg/sdk-wrapper"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"vectorx/pkg/intents"
	"vectorx/pkg/stats"
	"vectorx/pkg/vim-server"
)

const VECTORX_VERSION = "RELEASE_19"

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
	VectorSettings map[string]interface{}     `json:"vector_settings"`
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	wirepodPath := os.Getenv("WIREPOD_HOME")
	vPath := os.Getenv("VECTORX_HOME")
	botConfigJson := filepath.Join(wirepodPath, "chipper/jdocs/botSdkInfo.json")
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
				mapConfigVectorX["WEATHERAPI_PROVIDER"] = "openweathermap.org"
				mapConfigVectorX["WEATHERAPI_ENABLED"] = "true"
			} else {
				mapConfig["WEATHERAPI_PROVIDER"] = ""
				mapConfig["WEATHERAPI_ENABLED"] = "false"
				mapConfigVectorX["WEATHERAPI_PROVIDER"] = ""
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
			mapConfigVectorX["WEATHERAPI_UNIT"] = r.FormValue("weatherunits")
			mapConfig["STT_SERVICE"] = "vosk"
			mapConfig["CONN_SELECTION"] = r.FormValue("connSelection")
			mapConfig["WEBSERVER_PORT"] = "8080"

			err = JSONToVectorxConfig(mapConfigVectorX)
			err = JSONToWirepodConfig(mapConfig)

			if err != nil {
				fmt.Fprintf(w, "{ \"result\": \"KO\"}")
			}

			if !checkAndFixVosk() {
				fmt.Fprintf(w, "{ \"result\": \"KO\"}")
			} else {
				err = resetWirepod()
				if err != nil {
					fmt.Fprintf(w, "{ \"result\": \"KO\"}")
				} else {
					_, _ = os.Create(filepath.Join(os.Getenv("VECTORX_HOME"), ".setup"))
					fmt.Fprintf(w, "{ \"result\": \"OK\"}")
				}
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
				var vectorSettings map[string]interface{} = nil
				errSDK := sdk_wrapper.InitSDKForWirepod(bot.Esn)
				if errSDK == nil {
					vectorSettings = sdk_wrapper.GetVectorSettings()
				}
				var bi BotInfo = BotInfo{bot.Esn, bot.IpAddress, sdk_wrapper.CustomSettings{}, vectorSettings}
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
		usageStats, err := stats.GetStats()
		commands := "---"
		if err == nil {
			commands = fmt.Sprintf("%d/%d", usageStats.IntentsMatched, usageStats.IntentsReceived)
		}

		data := "{ \"uptime\": \"" + uptime + "\"," +
			"\"network\": \"" + network + "\"," +
			"\"status\": \"" + status + "\"," +
			"\"commands\": \"" + commands + "\"" +
			" }"
		fmt.Fprintf(w, data)
		break
	case r.URL.Path == "/api/get_vectorx_intents":
		var intents *[]intents.IntentDef = intents.GetIntents()
		jsonStr, err2 := json.Marshal(*intents)
		if err2 == nil {
			fmt.Fprintf(w, string(jsonStr))
		} else {
			fmt.Fprintf(w, "{}")
		}
		break
	case r.URL.Path == "/api/update":
		result, commandOutput := runUpdateScript()
		fmt.Fprintf(w, "{ \"result\": \""+result+"\", \"output\": \""+commandOutput+"\"}")
		break
	case r.URL.Path == "/api/send_intent":
		name := r.FormValue("name")
		serialNo := r.FormValue("esn")
		params := []string{r.FormValue("p1"), r.FormValue("p2"), r.FormValue("p3")}
		if name == "" || serialNo == "" {
			fmt.Fprintf(w, "{ \"result\": \"KO\"}")
			return
		}
		var e error
		if name == "tts-test" || name == "tts-configure" {
			e = runPseudoIntentCommand(name, serialNo, params)
		} else {
			e = runIntentCommand(name, serialNo, params)
		}
		if e == nil {
			fmt.Fprintf(w, "{ \"result\": \"OK\"}")
		} else {
			fmt.Fprintf(w, "{ \"result\": \"KO\", \"reason\": \""+e.Error()+"\"}")
		}
		break
	case r.URL.Path == "/api/vim_list_targets":
		result := "["
		data, err := ioutil.ReadFile(botConfigJson)
		if err != nil {
			fmt.Fprintf(w, "{}")
		} else {
			var jsonObj WirePodConfig
			json.Unmarshal(data, &jsonObj)

			for _, bot := range jsonObj.Robots {
				botCustomConfigJson := filepath.Join(vPath, "vectorfs/nvm/"+bot.Esn+"/custom_settings.json")
				data, err := ioutil.ReadFile(botCustomConfigJson)
				if err == nil {
					var customSettings sdk_wrapper.CustomSettings
					err = json.Unmarshal(data, &customSettings)
					if err == nil {
						remote := intents.VIMUserInfoData{DisplayName: customSettings.RobotName, UserId: bot.Esn, IsHuman: false}
						remoteJson, err := json.Marshal(remote)
						if err == nil {
							result = result + string(remoteJson) + ","
						}
					}
				}
			}
		}
		remote := intents.VIMUserInfoData{DisplayName: "Human", UserId: "00000000", IsHuman: true}
		remoteJson, err := json.Marshal(remote)
		if err == nil {
			result = result + string(remoteJson)
		}
		result = result + "]"
		fmt.Fprintf(w, result)
		break
	default:
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
}

func StartWebServer() {
	var webPort string
	intents.RegisterIntents()
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
	// VIM Server
	hub := vim_server.NewHub()
	go hub.Run()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		vim_server.ServeWs(hub, w, r)
	})
	fmt.Printf("Starting vectorxws at port " + webPort + " (http://localhost:" + webPort + ")\n")
	if err := http.ListenAndServe(":"+webPort, nil); err != nil {
		log.Fatal(err)
	}
}

type apiConfig struct {
	Weather struct {
		Enable   bool   `json:"enable"`
		Provider string `json:"provider"`
		Key      string `json:"key"`
		Unit     string `json:"unit"`
	} `json:"weather"`
	Knowledge struct {
		Enable      bool   `json:"enable"`
		Provider    string `json:"provider"`
		Key         string `json:"key"`
		ID          string `json:"id"`
		IntentGraph bool   `json:"intentgraph"`
		RobotName   string `json:"robotName"`
	} `json:"knowledge"`
	STT struct {
		Service  string `json:"provider"`
		Language string `json:"language"`
	} `json:"STT"`
	Server struct {
		// false for ip, true for escape pod
		EPConfig bool   `json:"epconfig"`
		Port     string `json:"port"`
	} `json:"server"`
	HasReadFromEnv   bool `json:"hasreadfromenv"`
	PastInitialSetup bool `json:"pastinitialsetup"`
}

func WirepodConfigToJSON() (map[string]string, error) {
	wirepodPath := os.Getenv("WIREPOD_HOME")
	wirepodCFG := filepath.Join(wirepodPath, "chipper/source.sh")
	// Read source.sh
	wirepodCFGMyJson, err := configToJson(wirepodCFG)
	if err != nil {
		return nil, err
	}
	fmt.Println("Source.sh read")
	// Integrate with the info now in apiConfig.json
	wirepodJSONCFG := filepath.Join(wirepodPath, "chipper/apiConfig.json")

	b, err := os.ReadFile(wirepodJSONCFG)
	if err != nil {
		fmt.Print(err)
	} else {
		fmt.Println("apiConfig.json read")
		var APIConfig apiConfig
		err = json.Unmarshal(b, &APIConfig)
		if err != nil {
			fmt.Print(err)
			return wirepodCFGMyJson, err
		}

		if APIConfig.Weather.Enable {
			wirepodCFGMyJson["WEATHERAPI_ENABLED"] = "true"
			wirepodCFGMyJson["WEATHERAPI_KEY"] = APIConfig.Weather.Key
			wirepodCFGMyJson["WEATHERAPI_UNIT"] = APIConfig.Weather.Unit
			wirepodCFGMyJson["WEATHERAPI_PROVIDER"] = APIConfig.Weather.Provider
		} else {
			wirepodCFGMyJson["WEATHERAPI_ENABLED"] = "false"
			wirepodCFGMyJson["WEATHERAPI_KEY"] = ""
			wirepodCFGMyJson["WEATHERAPI_UNIT"] = ""
			wirepodCFGMyJson["WEATHERAPI_PROVIDER"] = ""
		}
		if APIConfig.Knowledge.Enable {
			wirepodCFGMyJson["KNOWLEDGE_ENABLED"] = "true"
			wirepodCFGMyJson["KNOWLEDGE_KEY"] = APIConfig.Knowledge.Key
			wirepodCFGMyJson["KNOWLEDGE_PROVIDER"] = APIConfig.Knowledge.Provider
		} else {
			wirepodCFGMyJson["KNOWLEDGE_ENABLED"] = "false"
			wirepodCFGMyJson["KNOWLEDGE_KEY"] = ""
			wirepodCFGMyJson["KNOWLEDGE_PROVIDER"] = ""
		}
		wirepodCFGMyJson["STT_LANGUAGE"] = APIConfig.STT.Language
		wirepodCFGMyJson["CHIPPER_PORT"] = APIConfig.Server.Port
		//wirepodCFGMyJson["WEBSERVER_PORT"] = "8080"
		if APIConfig.Server.EPConfig == true {
			wirepodCFGMyJson["CONN_SELECTION"] = "ep"
			wirepodCFGMyJson["WIREPOD_CONSOLE"] = "http://escapepod.local:" + wirepodCFGMyJson["WEBSERVER_PORT"]
		} else {
			wirepodCFGMyJson["CONN_SELECTION"] = "ip"
			wirepodCFGMyJson["WIREPOD_CONSOLE"] = "http://" + GetOutboundIP().String() + ":" + wirepodCFGMyJson["WEBSERVER_PORT"]
		}
	}

	return wirepodCFGMyJson, nil
}

func JSONToWirepodConfig(cfg map[string]string) error {
	wirepodPath := os.Getenv("WIREPOD_HOME")
	wirepodCFG := filepath.Join(wirepodPath, "chipper/source.sh")
	// Writes to source.sh
	jsonToConfig(wirepodCFG, cfg)
	wirepodJSONCFG := filepath.Join(wirepodPath, "chipper/apiConfig.json")
	var APIConfig apiConfig = apiConfig{}

	if cfg["WEATHERAPI_ENABLED"] == "true" {
		APIConfig.Weather.Enable = true
		APIConfig.Weather.Key = cfg["WEATHERAPI_KEY"]
		APIConfig.Weather.Unit = cfg["WEATHERAPI_UNIT"]
		APIConfig.Weather.Provider = cfg["WEATHERAPI_PROVIDER"]
	} else {
		APIConfig.Weather.Enable = false
		APIConfig.Weather.Key = ""
		APIConfig.Weather.Unit = "C"
		APIConfig.Weather.Provider = ""
	}
	if cfg["KNOWLEDGE_ENABLED"] == "true" {
		APIConfig.Knowledge.Enable = true
		APIConfig.Knowledge.Key = cfg["KNOWLEDGE_KEY"]
		APIConfig.Knowledge.Provider = cfg["KNOWLEDGE_PROVIDER"]
		APIConfig.Knowledge.IntentGraph = true
	} else {
		APIConfig.Knowledge.Enable = false
		APIConfig.Knowledge.Key = ""
		APIConfig.Knowledge.Provider = ""
		APIConfig.Knowledge.IntentGraph = false
	}
	APIConfig.STT.Language = cfg["STT_LANGUAGE"]
	if cfg["CONN_SELECTION"] == "ep" {
		APIConfig.Server.EPConfig = true
	} else {
		APIConfig.Server.EPConfig = false
	}
	APIConfig.Server.Port = "443"
	APIConfig.HasReadFromEnv = true
	APIConfig.PastInitialSetup = true
	jSonData, err := json.Marshal(APIConfig)
	print(string(jSonData))

	err = os.WriteFile(wirepodJSONCFG, jSonData, 0644)
	if err != nil {
		print("failed opening file " + wirepodJSONCFG)
	}
	return err
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

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

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

	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0644)
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
		if strings.Contains(line, "export ") {
			line = strings.Split(line, "export ")[1]
			data := strings.Split(line, "=")
			mapConfig[data[0]] = data[1]
		}
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
	isOk := true
	var cmds = []string{
		"sudo systemctl start vectorx-update",
	}
	for _, cmd := range cmds {
		println(cmd)
		e := exec.Command("/bin/sh", "-c", cmd).Run()
		isOk = isOk && (e == nil)
	}

	if !isOk {
		return "error", "Error running update service!"
	}
	return "ok", "Update service started"
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

func runPseudoIntentCommand(intentName string, serialNo string, params []string) error {
	var err error = nil
	var Ctx = context.Background()
	var Start = make(chan bool)
	var Stop = make(chan bool)

	if intentName == "tts-test" || intentName == "tts-configure" {
		err = sdk_wrapper.InitSDKForWirepod(serialNo)
		if err != nil {
			println("FATAL: could not load Vector settings from JDOCS")
			return err
		}
		if intentName == "tts-configure" {
			engine, err2 := strconv.Atoi(params[0])
			language := params[1]
			voice := params[2]
			if err2 == nil {
				sdk_wrapper.SetLocale(language)
				sdk_wrapper.SetTTSConfiguration(engine, voice)
			} else {
				fmt.Println(err2.Error())
				return err2
			}
		} else if intentName == "tts-test" {
			sentence := params[0]
			go func() {
				_ = sdk_wrapper.Robot.BehaviorControl(Ctx, Start, Stop)
			}()

			for {
				select {
				case <-Start:
					sdk_wrapper.SayText(sentence)
					Stop <- true
				}
				return nil
			}
		}
		return nil
	}
	return errors.New("Unsupported command")
}

func runIntentCommand(intentName string, serialNo string, params []string) error {
	txt := ""
	if intentName == "roll-a-die" {
		txt = "roll a die"
	} else if intentName == "bingo" {
		txt = "bingo"
	} else if intentName == "pong" {
		txt = "let's play a classic"
	} else if intentName == "rps" {
		txt = "let's play a new game"
	} else if intentName == "how-do-you-say" {
		txt = "how do you say " + params[0] + " in " + params[1]
	} else if intentName == "lets-speak" {
		txt = "let's talk " + params[0]
	} else if intentName == "weather" {
		txt = "weather"
		if params[0] != "" {
			txt += " in " + params[0]
		}
	} else if intentName == "weather-forecast" {
		txt = "weather " + params[0]
		if params[1] != "" {
			txt += " in " + params[1]
		}
	} else if intentName == "set-name" {
		txt = "your name is " + params[0]
	} else if intentName == "pills-of-wisdom" {
		txt = "tell me something"
	} else if intentName == "oskr-trivia" {
		txt = "trivia game"
	}

	if txt == "" {
		return errors.New("unknown intent name")
	}

	vPath := "/home/pi/vectorx/" //os.Getenv("VECTORX_HOME")
	var cmds = []string{
		"sudo " + vPath + "runCmd.sh " + serialNo + " en-US \"" + txt + "\"",
	}
	for _, cmd := range cmds {
		println(cmd)
		// Use Start() instead of Run() so we don't wait until completion
		e := exec.Command("/bin/sh", "-c", cmd).Start()
		if e != nil {
			return e
		}
	}
	return nil
}
