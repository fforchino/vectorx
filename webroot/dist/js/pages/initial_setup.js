function checkForm() {
    alert("TODO!");
}

function goToTutorial(n) {
    if (checkTutorial(n)) {
        for (i=0;i<10;i++) {
            var el = document.getElementById("tut0"+i);
            if (el!=null) {
                el.style.display = i==n ? "block" : "none";
            }
        }
    }
}

function checkTutorial(n) {
    switch(n) {
        case 2:
            if (document.getElementById("weatherapi").value=="") {
                return confirm("You didn't enter a Weather API KEY. You won't be able to benefit of weather forecasts. Are you sure?");
            }
            return true;
        case 3:
            if (document.getElementById("kgapi").value=="") {
                return confirm("You didn't enter a Knowledge Graph API KEY. You won't be able to benefit of Vector's Q/A. Are you sure?");
            }
            return true;
    }
    return true;
}

function loadData() {
    fetch("/api/consistency_check")
        .then(response => response.text())
        .then((response) => {
            try {
                obj = JSON.parse(response);
                document.getElementById("weatherapi").value = obj.WEATHERAPI_KEY;
                document.getElementById("kgapi").value = obj.KNOWLEDGE_KEY;
                document.getElementById("kgprovider").value = obj.KNOWLEDGE_PROVIDER;
                document.getElementById("language").value = obj.STT_LANGUAGE;
                document.getElementById("weatherunits").value = obj.WEATHERAPI_UNIT;
            } catch {}
        })
}
loadData();

function saveData() {
    var weatherapi = document.getElementById("weatherapi").value;
    var kgapi = document.getElementById("kgapi").value;
    var kgprovider = document.getElementById("kgprovider").value;
    var language = document.getElementById("language").value;
    var weatherunits = document.getElementById("weatherunits").value;
    var data = "language=" + language +
               "&weatherapi="+weatherapi +
               "&kgapi="+kgapi +
               "&weatherunits="+weatherunits +
               "&kgprovider="+kgprovider;
    theUrl = "/api/initial_setup?" + data;
    //alert(theUrl);
    fetch(theUrl)
        .then(response => response.text())
        .then((response) => {
            var failure = false;
            try {
                obj = JSON.parse(response);
                if (obj.result!="OK") {
                    failure = true;
                }
            } catch { failure = true; }
            if (failure) {
                alert("There was a problem saving data!");
            } else {
                alert("Data saved successfully");
                goToTutorial(5);
            }
        })
}

function goHome() {
    document.location.replace("index.html")
}

function goOnboarding() {
    document.location.replace("onboarding.html")
}
