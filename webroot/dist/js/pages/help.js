function LoadHelpPage() {
    let selectTag = document.getElementById("help-vc")
    VectorxIntents.map( (vxIntent, i) => {
        let opt = document.createElement("option");
        opt.value = vxIntent.intentName;
        opt.innerHTML = vxIntent.intentName;
        selectTag.append(opt);
    });
}

function ChangeIntent() {
    let selectTag = document.getElementById("help-vc")
    VectorxIntents.map( (vxIntent, i) => {
        if (vxIntent.intentName==selectTag.value) {
            document.getElementById("utt-en").innerHTML = "";
            for (var j=0;j<vxIntent.keyPhrases["en"].length;j++) {
                document.getElementById("utt-en").innerHTML+=vxIntent.keyPhrases["en"][i];
            }
        }
    });
}