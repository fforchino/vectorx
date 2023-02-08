const WIREPOD_HOME = "http://escapepod:8080"
const QueryParams = new Proxy(new URLSearchParams(window.location.search), {
    get: (searchParams, prop) => searchParams.get(prop),
});

var Robots = {};
var Intents = {}
var Settings = {}

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

function GetCustomIntentByName(name) {
    for (var i=0;i<Intents.length;i++) {
        if (Intents[i].name == name) return Intents[i];
    }
    return null;
}

function goHome() {
    document.location.replace("index.html")
}

