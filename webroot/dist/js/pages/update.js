async function CheckUpdates() {
    document.getElementById("update_result_error").style.display="none";
    document.getElementById("update_result_no_update").style.display="none";
    document.getElementById("update_result_updated").style.display="none";
    document.getElementById("btut04").disabled = true;

    currentVersion = Settings.VECTORX_VERSION;
    document.getElementById("updates_update_status").innerHTML = "Current VectorX version is "+currentVersion+"<br/>";
    document.getElementById("updates_update_status").innerHTML += "Checking for updates, please wait...<br/>";
    RunUpdateScript().then((message) => {
        if (message=="ok") {
            document.getElementById("updates_update_status").innerHTML = "";
            LoadSettings().then(() => {
                if (Settings.VECTORX_VERSION == currentVersion) {
                    document.getElementById("update_result_no_update").style.display="block";
                    document.getElementById("updates_update_status").innerHTML = "<br/>No update found.";
                }
                else {
                    document.getElementById("update_result_updated").style.display="block";
                    document.getElementById("updates_update_status").innerHTML = "<br/>New release installed!<br>Current VectorX version is now "+Settings.VECTORX_VERSION;
                    document.getElementById("butReload").style.display="block";
                }
            });
        }
        document.getElementById("btut04").disabled = false;
    });
}

function ForceReload() {
    document.location = WIREPOD_HOME + "/index.html?ts=" + new Date().getTime();
}

async function RunUpdateScript() {
    var retVal = "";
    var obj = null;
    document.getElementById("update_running").style.display = "block";
    await fetch("/api/update", {cache: "no-store"})
        .then(response => response.text())
        .then((response) => {
            try {
                obj = JSON.parse(response);
                retVal = obj.result;
                document.getElementById("updates_update_status").innerHTML += "<br/>"+obj.output;
            } catch { retVal = "unknown"; }
        })
    if (obj!=null && obj.result=="ok") {
        for (var i=60;i>=0;i--) {
            document.getElementById("updates_counter").innerHTML = ""+i;
            await new Promise(r => setTimeout(r, 1000));
        }

    } else {
        document.getElementById("updates_update_status").innerHTML = "<br/>Error. Updates cannot be applied automatically, you'll have to check what's going on using SSH.";
        document.getElementById("update_result_error").style.display="block";
    }
    document.getElementById("update_running").style.display = "none";
    return Promise.resolve(retVal);
}