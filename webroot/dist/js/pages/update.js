async function CheckUpdates() {
    currentVersion = Settings.VECTORX_VERSION;
    document.getElementById("updates_update_status").innerHTML = "Current VectorX version is "+currentVersion+"<br/>";
    document.getElementById("updates_update_status").innerHTML += "Checking for updates, please wait...<br/>";
    RunUpdateScript().then((message) => {
        LoadSettings().then(() => {
            document.getElementById("updates_update_status").innerHTML += "<br/>Current VectorX version is now "+Settings.VECTORX_VERSION;
        });
    });
}

async function RunUpdateScript() {
    var retVal = "";
    var obj = null;
    document.getElementById("update_running").style.display = "block";
    await fetch("/api/update")
        .then(response => response.text())
        .then((response) => {
            try {
                obj = JSON.parse(response);
                retVal = obj.result;
                document.getElementById("updates_update_status").innerHTML += "<br/>"+obj.output;
            } catch { retVal = "unknown"; }
        })
    if (obj!=null && obj.result=="ok") {
        for (var i=30;i>=0;i--) {
            await new Promise(r => setTimeout(r, 1000));
            document.getElementById("updates_update_status").innerHTML += ".";
        }

    }
    document.getElementById("update_running").style.display = "none";
    return Promise.resolve(retVal);
}