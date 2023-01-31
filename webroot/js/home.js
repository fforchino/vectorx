function doConsistencyCheck() {
    fetch("/api/consistency_check")
        .then(response => response.text())
        .then((response) => {
            /*
            obj = JSON.parse(response);
            output = "Wirepod uses VOSK...................";
            output += (obj.sttProvider == "vosk") ? "OK" : "KO";
            */
            document.getElementById("checks").innerHTML = response;
        })
}
doConsistencyCheck()