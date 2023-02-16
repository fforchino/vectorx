const WIREPOD_HOME = "http://escapepod:8080"
const QueryParams = new Proxy(new URLSearchParams(window.location.search), {
    get: (searchParams, prop) => searchParams.get(prop),
});

var Robots = {};
var Intents = {};
var Settings = {};
var VectorxIntents = {};

async function LoadRobots() {
    await fetch("/api/get_robots", {cache: "no-store"})
        .then(response => response.text())
        .then((response) => {
            try {
                obj = JSON.parse(response);
                Robots = obj;
            } catch {}
        })
}

async function LoadIntents() {
    await fetch(WIREPOD_HOME+"/api/get_custom_intents_json", {cache: "no-store"})
        .then(response => response.text())
        .then((response) => {
            try {
                obj = JSON.parse(response);
                Intents = obj;
            } catch {}
        })
}

async function LoadSettings() {
    await fetch("/api/consistency_check", {cache: "no-store"})
        .then(response => response.text())
        .then((response) => {
            try {
                //alert(response);
                obj = JSON.parse(response);
                Settings = obj;
            } catch {}
        })
}

async function LoadVectorXCustomIntents() {
    await fetch("/api/get_vectorx_intents", {cache: "no-store"})
        .then(response => response.text())
        .then((response) => {
            try {
                obj = JSON.parse(response);
                VectorxIntents = obj;
            } catch {}
        })
}

function GetCustomIntentByName(name) {
    for (var i=0;i<Intents.length;i++) {
        if (Intents[i].name == name) return Intents[i];
    }
    return null;
}

function goHome() {
    document.location.replace("index.html")
}

function goInitialSetup() {
    document.location.replace("initial_setup.html")
}

function GetRobotInfo(esn) {
    var theBot = null;
    for (var i = 0; i < Robots.length; i++) {
        var bot = Robots[i];
        if (bot.esn==esn) {
            theBot = bot;
            break;
        }
    }
    return theBot;
}

function GetRobotEyeColorRGB(bot) {
    var eyeColor = "#00ff00";

    if (bot.vector_settings.custom_eye_color.enabled) {

        eyeColor = hslToHex(parseFloat(bot.vector_settings.custom_eye_color.hue*360),
                            parseFloat(bot.vector_settings.custom_eye_color.saturation)*100,
                            50);
    }
    else {
        switch (bot.vector_settings.eye_color) {
            case 0: //TIP_OVER_TEAL
                eyeColor = "#29ae70ff";
                break;
            case 1: //OVERFIT_ORANGE
                eyeColor = "#fe7614ff";
                break;
            case 2: //UNCANNY_YELLOW
                eyeColor = "#f7cb04ff";
                break;
            case 3: //NON_LINEAR_LIME
                eyeColor = "#a8d304ff";
                break;
            case 4: //SINGULARITY_SAPPHIRE
                eyeColor = "#0d97f0ff";
                break;
            case 5: //FALSE_POSITIVE_PURPLE
                eyeColor = "#c661fcff";
                break;
            case 6: //CONFUSION_MATRIX_GREEN
                eyeColor = "#00ff00ff";
                break;
        }
    }
    return eyeColor;
}

function hslToHex(h, s, l) {
    l /= 100;
    const a = s * Math.min(l, 1 - l) / 100;
    const f = n => {
        const k = (n + h / 30) % 12;
        const color = l - a * Math.max(Math.min(k - 3, 9 - k, 1), -1);
        return Math.round(255 * color).toString(16).padStart(2, '0');   // convert to Hex and prefix "0" if needed
    };
    return `#${f(0)}${f(8)}${f(4)}`;
}