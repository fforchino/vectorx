var CurrentRobot = null;

function LoadBotControlPage() {
    CurrentRobot = GetRobotInfo(QueryParams.esn)
    if (bot!=null) {
        let botName = CurrentRobot.custom_settings.RobotName.toUpperCase();
        document.getElementById("nav_page_robots").click();
        document.getElementById("nav_page_botcontrol_"+CurrentRobot.esn).classList.add("active");
        document.getElementById("bc_serial_no").innerHTML = botName
        document.getElementById("bc_title").innerHTML = "Play with "+botName;
    }
}

function BotControlReveal(element, elementIdToShow) {
    element.style.display = "none";
    var elementsToHide = getElementsByClassName("bot-control-reveal");
    for (var i=0;i<elementsToHide.length;i++) {
        elementsToHide[i].style.display = "none";
    }
    document.getElementById(elementIdToShow).style.display = "block";
}

async function BotControlSendIntent(intentName, resultElement) {
    if (intentName=="roll-a-die" ||
        intentName=="bingo") {
        var data = "name=" + intentName+"&esn="+CurrentRobot.esn;
        fetch("/api/send_intent?" + data)
            .then(response => response.text())
            .then((response) => {
                try {
                    //alert(response);
                    obj = JSON.parse(response);
                    if (obj.result=="OK") {
                        document.getElementById(resultElement).innerHTML="Command successfully sent to Vector.";
                        return;
                    }
                } catch {}
                document.getElementById(resultElement).innerHTML="Error sending command.";
            })
    }
}
