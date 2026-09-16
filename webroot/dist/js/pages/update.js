async function CheckUpdates() {
    setUpdateState("checking");
    const status = document.getElementById("updates_update_status");
    const checkButton = document.getElementById("butCheckUpdate");
    const applyButton = document.getElementById("butApplyUpdate");

    checkButton.disabled = true;
    applyButton.style.display = "none";
    status.innerHTML = "Current VectorX version is " + Settings.VECTORX_VERSION + "<br/>";
    status.innerHTML += "Checking for updates, please wait...<br/>";

    try {
        const response = await fetch("/api/check_update", {cache: "no-store"});
        const obj = await response.json();
        if (obj.result !== "ok") {
            setUpdateState("error");
            status.innerHTML = "<br/>" + (obj.output || "Could not check for updates.");
            return;
        }
        if (obj.update_available) {
            setUpdateState("available");
            status.innerHTML = "<br/>A new release is available.";
            status.innerHTML += "<br/>Installed version: " + obj.current_version;
            status.innerHTML += "<br/>Available version: " + obj.available_version;
            applyButton.style.display = "block";
        } else {
            setUpdateState("none");
            status.innerHTML = "<br/>No update found.";
            if (obj.available_version) {
                status.innerHTML += "<br/>Latest available version: " + obj.available_version;
            }
        }
    } catch (e) {
        setUpdateState("error");
        status.innerHTML = "<br/>Could not check for updates.";
    } finally {
        checkButton.disabled = false;
    }
}

async function ApplyUpdate() {
    if (!confirm("Apply the available VectorX update now? VectorX will restart during the update.")) {
        return;
    }
    document.getElementById("butApplyUpdate").disabled = true;
    document.getElementById("butCheckUpdate").disabled = true;
    document.getElementById("updates_update_status").innerHTML = "Starting the update service...<br/>";
    const result = await RunUpdateScript();
    if (result === "ok") {
        LoadSettings().then(() => {
            document.getElementById("update_result_updated").style.display = "block";
            document.getElementById("updates_update_status").innerHTML = "<br/>Update started. Reload VectorX after the services restart.";
            document.getElementById("butReload").style.display = "block";
        });
    }
    document.getElementById("butCheckUpdate").disabled = false;
    document.getElementById("butApplyUpdate").disabled = false;
}

function ForceReload() {
    document.location = "index.html?ts=" + new Date().getTime();
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
                document.getElementById("updates_update_status").innerHTML += "<br/>" + obj.output;
            } catch { retVal = "unknown"; }
        });
    if (obj != null && obj.result == "ok") {
        for (var i = 60; i >= 0; i--) {
            document.getElementById("updates_counter").innerHTML = "" + i;
            await new Promise(r => setTimeout(r, 1000));
        }

    } else {
        document.getElementById("updates_update_status").innerHTML = "<br/>Error. Updates cannot be applied automatically, you'll have to check what's going on using SSH.";
        document.getElementById("update_result_error").style.display = "block";
    }
    document.getElementById("update_running").style.display = "none";
    return Promise.resolve(retVal);
}

function setUpdateState(state) {
    document.getElementById("update_running").style.display = state === "checking" ? "block" : "none";
    document.getElementById("update_result_error").style.display = state === "error" ? "block" : "none";
    document.getElementById("update_result_no_update").style.display = state === "none" ? "block" : "none";
    document.getElementById("update_result_updated").style.display = state === "available" ? "block" : "none";
}
