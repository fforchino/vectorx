function doConsistencyCheck() {
    try {
        if (Settings.STT_SERVICE == "vosk") {
            document.getElementById("span_cc1").classList.add("text-success");
            document.getElementById("i_cc1").classList.add("fa-check");
        } else {
            document.getElementById("span_cc1").classList.add("text-danger");
            document.getElementById("i_cc1").classList.add("fa-xmark");
        }
        if (Settings.WEATHERAPI_PROVIDER == "openweathermap.org") {
            document.getElementById("span_cc2").classList.add("text-success");
            document.getElementById("i_cc2").classList.add("fa-check");
        } else {
            document.getElementById("span_cc2").classList.add("text-danger");
            document.getElementById("i_cc2").classList.add("fa-xmark");
        }
        if ((Settings.KNOWLEDGE_PROVIDER != "") && (Settings.KNOWLEDGE_KEY!="")) {
            document.getElementById("span_cc3").classList.add("text-success");
            document.getElementById("i_cc3").classList.add("fa-check");
        } else {
            document.getElementById("span_cc3").classList.add("text-danger");
            document.getElementById("i_cc3").classList.add("fa-xmark");
        }
        if (Settings.VOSK_OK == "true") {
            document.getElementById("span_cc4").classList.add("text-success");
            document.getElementById("i_cc4").classList.add("fa-check");
        } else {
            document.getElementById("span_cc4").classList.add("text-danger");
            document.getElementById("i_cc4").classList.add("fa-xmark");
        }
        if (Settings.CHIPPER_PORT == "443") {
            document.getElementById("span_cc5").classList.add("text-success");
            document.getElementById("i_cc5").classList.add("fa-check");
        } else {
            document.getElementById("span_cc5").classList.add("text-danger");
            document.getElementById("i_cc5").classList.add("fa-xmark");
        }
        /*
        if (Settings.WEBSERVER_PORT == "8080") {
            document.getElementById("span_cc6").classList.add("text-success");
            document.getElementById("i_cc6").classList.add("fa-check");
        } else {
            document.getElementById("span_cc6").classList.add("text-danger");
            document.getElementById("i_cc6").classList.add("fa-xmark");
        }
        */
        document.getElementById("home_locale").innerHTML = Settings.STT_LANGUAGE;
    } catch {}
    fetch("/api/get_stats")
        .then(response => response.text())
        .then((response) => {
            obj = JSON.parse(response);
            document.getElementById("home_uptime").innerHTML = obj.uptime;
            document.getElementById("home_ssid").innerHTML = obj.network
            document.getElementById("home_status").innerHTML = obj.status;
            document.getElementById("home_commands").innerHTML = obj.commands;
        })
}

function LoadHomePageBots() {
    var data = "";
    for (var i = 0; i < Robots.length; i++) {
        var bot = Robots[i];
        var ip = "OFFLINE";
        var botCtrlLink = "#"
        var offlineTools = "";
        if (bot.vector_settings==null) {
            // Bot offline
            eyeColor = "#aaaaaa";
            offlineTools =
                '                    <div class="input-group input-group-sm mt-2">\n' +
                '                        <input id="robot-ip-'+bot.esn+'" type="text" class="form-control" placeholder="New IP address" value="'+bot.ip_address+'">\n' +
                '                        <span class="input-group-append">\n' +
                '                            <button type="button" class="btn btn-primary btn-flat" onclick="UpdateRobotIP(&quot;'+bot.esn+'&quot;)">Update IP</button>\n' +
                '                        </span>\n' +
                '                    </div>\n' +
                '                    <small id="robot-ip-status-'+bot.esn+'" class="form-text text-muted">Use this if Vector changed IP after rebooting.</small>\n';
        }
        else {
            eyeColor = GetRobotEyeColorRGB(bot);
            ip = bot.ip_address;
            botCtrlLink = "botcontrol.html?esn="+bot.esn;
        }
        var botName = ((bot.custom_settings && bot.custom_settings.RobotName) || "").toUpperCase();
        if (botName.length==0) {
            botName = bot.esn.toUpperCase();
        }
        data += '<div class="col-12 col-sm-6 col-md-3">\n' +
            '            <div class="info-box">\n' +
            '                <span class="info-box-icon elevation-1" style="background-color: '+eyeColor+'">\n'+
            '                    <a id="nav_page_botcontrol_'+bot.esn+'" href="'+botCtrlLink+'" class="bot-link">\n ' +
            '                        <i class="fas fa-robot" aria-hidden="true"></i>' +
            '                    </a>\n'+
            '                </span>\n' +
            '\n' +
            '                <div class="info-box-content">\n' +
            '                    <span class="info-box-text">'+botName+'</span>\n' +
            '                    <span class="info-box-number">\n' +
            '                  '+bot.esn.toUpperCase()+' | '+ip+'\n' +
            '                </span>\n' +
            offlineTools +
            '                </div>\n' +
            '                <!-- /.info-box-content -->\n' +
            '            </div>\n' +
            '            <!-- /.info-box -->\n' +
            '        </div>';
    }
    document.getElementById("row_homepage_bots").insertAdjacentHTML('afterbegin', data);
}

async function UpdateRobotIP(esn) {
    const input = document.getElementById("robot-ip-" + esn);
    const status = document.getElementById("robot-ip-status-" + esn);
    const ip = (input.value || "").trim();
    if (!/^((25[0-5]|2[0-4]\d|1?\d?\d)\.){3}(25[0-5]|2[0-4]\d|1?\d?\d)$/.test(ip)) {
        status.className = "form-text text-danger";
        status.innerHTML = "Enter a valid IPv4 address.";
        return;
    }
    status.className = "form-text text-warning";
    status.innerHTML = "Updating robot IP and restarting robot services...";
    try {
        const body = new URLSearchParams();
        body.set("esn", esn);
        body.set("ip", ip);
        const response = await fetch("/api/update_robot_ip", {
            method: "POST",
            headers: {"Content-Type": "application/x-www-form-urlencoded"},
            body: body.toString()
        });
        const text = await response.text();
        let result = {};
        try { result = JSON.parse(text); } catch {}
        if (response.ok && result.result === "OK") {
            status.className = "form-text text-success";
            status.innerHTML = "IP updated. Reloading robot status...";
            setTimeout(() => window.location.reload(), 2000);
        } else {
            status.className = "form-text text-danger";
            status.innerHTML = result.reason || "Could not update robot IP.";
        }
    } catch (e) {
        status.className = "form-text text-danger";
        status.innerHTML = "Could not update robot IP: " + e.message;
    }
}
