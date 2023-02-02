function doConsistencyCheck() {
    fetch("/api/consistency_check")
        .then(response => response.text())
        .then((response) => {
            obj = JSON.parse(response);
            output = "Wirepod uses VOSK...";
            output += (obj.STT_SERVICE == "vosk") ? "<span class='ok'>OK</span>" : "<span class='ko'>KO</span>";
            output = "VOSK has all the needed language models...";
            output += (obj.VOSK_OK == "true") ? "<span class='ok'>OK</span>" : "<span class='ko'>KO</span>";
            output += "<br>Wirepod uses OpenWeatherMap.org...";
            output += ((obj.WEATHERAPI_PROVIDER == "openweathermap.org") && (obj.WEATHERAPI_KEY!="")) ? "<span class='ok'>OK</span>" : "<span class='ko'>KO</span>";
            output += "<br>Wirepod uses a Knowledge Provider...";
            output += ((obj.KNOWLEDGE_PROVIDER != "") && (obj.KNOWLEDGE_KEY!="")) ? "<span class='ok'>OK</span>" : "<span class='ko'>KO</span>";
            output += "<br>Wirepod uses port 8080...";
            output += (obj.WEBSERVER_PORT == "8080") ? "<span class='ok'>OK</span>" : "<span class='ko'>KO</span>";
            output += "<br><br>VECTORX VERSION: "+obj.VECTORX_VERSION;

            document.getElementById("fix-problems").style.display = (output.includes("'ko'")) ? "block" : "none";
            document.getElementById("checks").innerHTML = output;
            document.getElementById("wirepod-home").href = "http://escapepod:"+obj.WEBSERVER_PORT+"/";
        })
}

function goInitialSetup() {
   document.location.replace("initial_setup.html")
}

function checkSetupMissing() {
        fetch("/api/is_setup_done")
            .then(response => response.text())
            .then((response) => {
                    try {
                            obj = JSON.parse(response);
                            if (obj.result=="OK") {
                                    doConsistencyCheck()
                            }
                            else {
                                    goInitialSetup();
                            }
                    } catch { goInitialSetup(); }
            })
}
checkSetupMissing();
