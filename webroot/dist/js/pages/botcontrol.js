var CurrentRobot = null;

async function LoadVoices(lang) {
    let dropdown = $('#language-ls-voice');
    dropdown.empty();
    dropdown.append('<option selected="true" disabled>Choose Voice</option>');
    dropdown.prop('selectedIndex', 0);
    const url = 'https://www.wondergarden.app/voiceserver/index.php/getVoices?lang='+lang;

    // Populate dropdown
    $.getJSON(url, function (data) {
        $.each(data, function (key, entry) {
            dropdown.append($('<option></option>').attr('value', entry.id).text(entry.name));
        })
        document.getElementById("language-ls-ttslanguage").value = lang;
        LoadTTSSettings();
    });
}

function LoadBotControlPage() {
    CurrentRobot = GetRobotInfo(QueryParams.esn)
    if (CurrentRobot!=null) {
        let botName = CurrentRobot.custom_settings.RobotName.toUpperCase();
        document.getElementById("nav_page_robots").click();
        document.getElementById("nav_page_botcontrol_"+CurrentRobot.esn).classList.add("active");
        document.getElementById("bc_serial_no").innerHTML = botName
        document.getElementById("bc_title").innerHTML = "Play with "+botName;
    }
}

function LoadTTSSettings() {
    document.getElementById("language-ls-engine").value = ""+CurrentRobot.custom_settings.TTSEngine;
    var hasVoice = !! $('#language-ls-voice > option[value="'+CurrentRobot.custom_settings.TTSVoice+'"]').length;
    if (hasVoice) {
        document.getElementById("language-ls-voice").value = CurrentRobot.custom_settings.TTSVoice;
    }
    BotControlHandleTTSEngineChange();
}

function BotControlReveal(element, elementIdToShow) {
    // Show all the "try it"...
    var elementsToShow = document.getElementsByClassName("small-box-footer");
    for (var i=0;i<elementsToShow.length;i++) {
        elementsToShow[i].style.display = "block";
    }
    // ...Except the selcted one
    element.style.display = "none";
    var elementsToHide = document.getElementsByClassName("bot-control-reveal");
    for (var i=0;i<elementsToHide.length;i++) {
        elementsToHide[i].style.display = "none";
    }
    document.getElementById(elementIdToShow).style.display = "block";
}

async function BotControlSendIntent(intentName, resultElement) {
    var extraData = "";
    // Handle parametric intents
    if (intentName=="how-do-you-say") {
        var word = document.getElementById("language-hds-word").value;
        var language = document.getElementById("language-hds-language").value;
        word = word.replaceAll('&', '');
        word = word.replaceAll('?', '');
        if (word!="" && language!="") {
            extraData += "&p1="+word+"&p2="+language;
        } else {
            document.getElementById(resultElement).innerHTML="Input error. Check parameters.";
            return
        }
    }
    else if (intentName=="lets-speak") {
        var language = document.getElementById("language-ls-language").value;
        if (language!="") {
            extraData += "&p1="+language;
        } else {
            document.getElementById(resultElement).innerHTML="Input error. Check parameters.";
            return
        }
    }
    else if (intentName=="weather") {
        var location = document.getElementById("weather-now-location").value;
        location = location.replaceAll('&', '');
        location = location.replaceAll('?', '');
        if (location!="") {
            extraData += "&p1="+location;
        }
    }
    else if (intentName=="weather-forecast") {
        var location = document.getElementById("weather-fc-location").value;
        var dt = document.getElementById("weather-fc-dt").value;
        location = location.replaceAll('&', '');
        location = location.replaceAll('?', '');
        if (dt!="") {
            extraData += "&p1="+dt;
        }
        else {
            document.getElementById(resultElement).innerHTML="Input error. Check parameters.";
            return
        }
        if (location!="") {
            extraData += "&p2="+location;
        }
    }
    else if (intentName=="set-name") {
        var name = document.getElementById("chat-name").value;
        name = name.replaceAll('&', '');
        name = name.replaceAll('?', '');
        if (name!="") {
            extraData += "&p1="+name;
        } else {
            document.getElementById(resultElement).innerHTML="Input error. Check parameters.";
            return
        }
    }
    else if (intentName=="tts-test") {
        var sentence = document.getElementById("language-ttslanguage-word").value;
        sentence = encodeURIComponent(sentence);
        if (sentence!="") {
            extraData += "&p1="+sentence;
        } else {
            document.getElementById(resultElement).innerHTML="Input error. Check parameters.";
            return
        }
    }
    else if (intentName=="tts-configure") {
        var engine = document.getElementById("language-ls-engine").value;
        var language = document.getElementById("language-ls-ttslanguage").value;
        var voice = document.getElementById("language-ls-voice").value;

        if (engine != "" && language != "") {
            extraData += "&p1=" + engine + "&p2=" + language;
            if (voice !="") {
                extraData += "&p3=" + voice;
            }
        } else {
            document.getElementById(resultElement).innerHTML = "Input error. Check parameters.";
            return
        }
    }
    if (intentName=="roll-a-die" ||
        intentName=="bingo" ||
        intentName=="pong" ||
        intentName=="rps" ||
        intentName=="how-do-you-say" ||
        intentName=="lets-speak" ||
        intentName=="weather" ||
        intentName=="weather-forecast" ||
        intentName=="set-name" ||
        intentName=="pills-of-wisdom" ||
        intentName=="tts-configure" ||
        intentName=="tts-test" ||
        intentName=="oskr-trivia") {
        var data = "name=" + intentName+"&esn="+CurrentRobot.esn+extraData;
        await fetch("/api/send_intent?" + data)
            .then(response => response.text())
            .then((response) => {
                var res = "Error sending command.";
                try {
                    //alert(response);
                    obj = JSON.parse(response);
                    if (obj.result=="OK") {
                        res="Command sent to Vector.";
                    }
                } catch {}
                document.getElementById(resultElement).innerHTML=res;
            })
        if (intentName=="set-name") {
            await new Promise(r => setTimeout(r, 5000));
            ReloadSite();
        }
    }
}

function BotControlHandleTTSEngineChange() {
    var engine = document.getElementById("language-ls-engine").value;
    if (engine!="2") {
        document.getElementById("language-ls-voice").style.display = "none"
        document.getElementById("language-ls-voice-label").style.display = "none"
    }
    else {
        document.getElementById("language-ls-voice").style.display = "block"
        document.getElementById("language-ls-voice-label").style.display = "block"
    }
}
