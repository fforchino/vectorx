function CheckUpdates() {
    currentVersion = Settings.VECTORX_VERSION;
    document.getElementById("updates_update_status").innerHTML = "Current VectorX version is "+currentVersion+"<br/>";
    document.getElementById("updates_update_status").innerHTML += "Checking for updates, please wait...<br/>";
    RunUpdateScript().then((message) => {
        if (message=="updated") {
            document.getElementById("updates_update_status").innerHTML += "<br/>Updates found and applied!";
            LoadSettings().then(() => {
                if (currentVersion!=Settings.VECTORX_VERSION) {
                    document.getElementById("updates_update_status").innerHTML += "<br/>Current VectorX version is now "+Settings.VECTORX_VERSION;
                }
            });
        } else {
            document.getElementById("updates_update_status").innerHTML += "<br/>No update found.";
        }
    });
}

async function RunUpdateScript() {
    var retVal = "";
    await new Promise(r => setTimeout(r, 2000));
    await fetch("/api/update")
        .then(response => response.text())
        .then((response) => {
            try {
                obj = JSON.parse(response);
                retVal = obj.result;
            } catch { retVal = "unknown"; }
        })
    return Promise.resolve(retVal);
}