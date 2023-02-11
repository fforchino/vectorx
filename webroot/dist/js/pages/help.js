function LoadHelpPage() {
    let selectTag = document.getElementById("help-vc")
    VectorxIntents.map( (vxIntent, i) => {
        let opt = document.createElement("option");
        opt.value = vxIntent.intentName;
        opt.innerHTML = vxIntent.intentName;
        selectTag.append(opt);
    });
    ChangeIntent();
}

function ChangeIntent() {
    let selectTag = document.getElementById("help-vc")
    VectorxIntents.map( (vxIntent, i) => {
        if (vxIntent.intentName==selectTag.value) {
            loadUtterances(vxIntent, "en");
            loadUtterances(vxIntent, "it");
            loadUtterances(vxIntent, "es");
            loadUtterances(vxIntent, "fr");
            loadUtterances(vxIntent, "de");
        }
    });
}

function loadUtterances(vxIntent, lang) {
    var elName = "utt-"+lang;
    document.getElementById(elName).innerHTML = "";
    for (var j=0;j<vxIntent.utterances[lang].length;j++) {
        document.getElementById(elName).innerHTML+=vxIntent.utterances[lang][j]+"</br>";
    }
}