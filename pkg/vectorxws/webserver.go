package vectorxws

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func apiHandler(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.URL.Path == "/api/get_kg_api":
		mapConfig, err := wirepodConfigToJson()
		if err != nil {
			fmt.Fprintf(w, "{}")
		} else {
			jsonStr, err2 := json.Marshal(mapConfig)
			if err2 != nil {
				fmt.Fprintf(w, "{}")
			} else {
				fmt.Fprintf(w, string(jsonStr))
			}
		}
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

func wirepodConfigToJson() (map[string]string, error) {
	wirepodPath := os.Getenv("WIREPOD_HOME")
	wirepodCFG := filepath.Join(wirepodPath, "chipper/source.sh")
	println("Looking at " + wirepodCFG + "...")
	mapConfig := map[string]string{}

	file, err := os.Open(wirepodCFG)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		line := scanner.Text()
		data := strings.Split(line, "=")
		mapConfig[data[0]] = data[1]
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return mapConfig, nil
}
