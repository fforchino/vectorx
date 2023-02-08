function LoadBotControlPage() {
    let botSerialNo = QueryParams.esn;
    let bot = GetRobotInfo(botSerialNo)
    if (bot!=null) {
        let botName = bot.custom_settings.RobotName.toUpperCase();
        document.getElementById("nav_page_robots").click();
        document.getElementById("nav_page_botcontrol_"+botSerialNo).classList.add("active");
        document.getElementById("bc_serial_no").innerHTML = botName
        document.getElementById("bc_title").innerHTML = "Play with "+botName;
    }
}

function BotControlReveal(element, elementIdToShow) {
    element.style.display = "none";
    document.getElementById(elementIdToShow).style.display = "block";
}

async function BotControlSendIntent(intentName, esn) {
    if (intentName=="roll-a-die") {
        var data = "name=" + intentName+"&esn="+esn;
        fetch(WIREPOD_HOME + "/api/send_intent?" + data)
            .then(response => response.text())
            .then((response) => {
                alert(response)
            })
    }
}
