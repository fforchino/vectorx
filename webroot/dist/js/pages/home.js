function doConsistencyCheck() {
    try {
        if (Settings.STT_SERVICE == "vosk") {
            document.getElementById("span_cc1").classList.add("text-success");
            document.getElementById("i_cc1").classList.add("fa-check");
        } else {
            document.getElementById("span_cc1").classList.add("text-danger");
            document.getElementById("i_cc1").classList.add("fa-xmark");
        }
        if (Settings.VOSK_OK == "true") {
            document.getElementById("span_cc2").classList.add("text-success");
            document.getElementById("i_cc2").classList.add("fa-check");
        } else {
            document.getElementById("span_cc2").classList.add("text-danger");
            document.getElementById("i_cc2").classList.add("fa-xmark");
        }
        if (Settings.WEATHERAPI_PROVIDER == "openweathermap.org") {
            document.getElementById("span_cc3").classList.add("text-success");
            document.getElementById("i_cc3").classList.add("fa-check");
        } else {
            document.getElementById("span_cc3").classList.add("text-danger");
            document.getElementById("i_cc3").classList.add("fa-xmark");
        }
        if ((Settings.KNOWLEDGE_PROVIDER != "") && (obj.KNOWLEDGE_KEY!="")) {
            document.getElementById("span_cc4").classList.add("text-success");
            document.getElementById("i_cc4").classList.add("fa-check");
        } else {
            document.getElementById("span_cc4").classList.add("text-danger");
            document.getElementById("i_cc4").classList.add("fa-xmark");
        }
        if (Settings.WEBSERVER_PORT == "8080") {
            document.getElementById("span_cc5").classList.add("text-success");
            document.getElementById("i_cc5").classList.add("fa-check");
        } else {
            document.getElementById("span_cc5").classList.add("text-danger");
            document.getElementById("i_cc5").classList.add("fa-xmark");
        }
        document.getElementById("home_locale").innerHTML = obj.STT_LANGUAGE;
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

function LoadHomePageBots() {
    var data = "";
    for (var i = 0; i < Robots.length; i++) {
        var bot = Robots[i];
        data += '<div class="col-12 col-sm-6 col-md-3">\n' +
            '            <div class="info-box">\n' +
            '                <span class="info-box-icon bg-info elevation-1"><i class="fas fa-robot" aria-hidden="true"></i></span>\n' +
            '\n' +
            '                <div class="info-box-content">\n' +
            '                    <span class="info-box-text">'+bot.custom_settings.RobotName.toUpperCase()+'</span>\n' +
            '                    <span class="info-box-number">\n' +
            '                  '+bot.esn.toUpperCase()+' | '+bot.ip_address+'\n' +
            '                </span>\n' +
            '                </div>\n' +
            '                <!-- /.info-box-content -->\n' +
            '            </div>\n' +
            '            <!-- /.info-box -->\n' +
            '        </div>';
    }
    document.getElementById("row_homepage_bots").insertAdjacentHTML('afterbegin', data);
}
