class ChatMessage {
    constructor(msgFrom, msgFromId, msgTo, msgToId, msgLang, msgText) {
        this.from = msgFrom;
        this.to = msgTo;
        this.fromId = msgFromId;
        this.toId = msgToId;
        this.lang = msgLang;
        this.msg = msgText;
    }
}

var targetSerialNo = "";
var targetName = "";
var fromSerialNo = "";
var fromName = "";
var chatLanguage = "en-US";

var chatDebug = false;

function LoadChatPage() {
    var conn;
    var msg = document.getElementById("msg");
    var log = document.getElementById("log");

    LoadSources();
    LoadTargets();

    function appendLog(item) {
        var doScroll = log.scrollTop > log.scrollHeight - log.clientHeight - 1;
        log.appendChild(item);
        if (doScroll) {
            log.scrollTop = log.scrollHeight - log.clientHeight;
        }
    }

    document.getElementById("form").onsubmit = function () {
        if (!conn) {
            return false;
        }
        if (msg.value=="" || fromSerialNo=="" || targetSerialNo=="") {
            return false;
        }
        var realTargetName = targetName;
        var humanName = document.getElementById('chat-human-name').value
        if (realTargetName.toLowerCase()=="human" && humanName!="") {
            realTargetName = humnName;
        }
        var chatMsg = new ChatMessage(fromName, fromSerialNo, realTargetName, targetSerialNo, chatLanguage, msg.value);
        conn.send(JSON.stringify(chatMsg));
        msg.value = "";
        return false;
    };

    if (window["WebSocket"]) {
        conn = new WebSocket("ws://" + document.location.host + "/ws");
        conn.onclose = function (evt) {
            var item = document.createElement("div");
            item.innerHTML = "<b>Connection closed.</b>";
            appendLog(item);
        };
        conn.onmessage = function (evt) {
            var messages = evt.data.split('\n');
            for (var i = 0; i < messages.length; i++) {
                var item = document.createElement("div");
                if (chatDebug) {
                    item.innerText = messages[i];
                } else {
                    var chatMsg = JSON.parse(messages[i]);
                    item.innerText = chatMsg.from+" -> "+chatMsg.to+" > "+chatMsg.msg;
                }
                appendLog(item);
            }
        };
    } else {
        var item = document.createElement("div");
        item.innerHTML = "<b>Your browser does not support WebSockets.</b>";
        appendLog(item);
    }
}

function loadDropDownHelper(name, id) {
    let dropdown = $(id);
    dropdown.empty();
    dropdown.append('<option selected="true" disabled>'+name+'</option>');
    dropdown.prop('selectedIndex', 0);

    const url = 'http://'+document.location.host+'/api/vim_list_targets';

    // Populate dropdown
    $.getJSON(url, function (data) {
        $.each(data, function (key, entry) {
            dropdown.append($('<option></option>').attr('value', entry.user_id).text(entry.display_name));
        })
    });
}

function LoadSources() {
    loadDropDownHelper("Choose Source", "#chat-from");
}

function LoadTargets() {
    loadDropDownHelper("Choose Target", "#chat-to");
}

function ChangeChatTarget() {
    var el = document.getElementById('chat-to');
    targetSerialNo = el.value;
    targetName = el.options[el.selectedIndex].innerHTML;
    document.getElementById('chat-send').disabled = (targetSerialNo=="" || fromSerialNo=="");
}

function ChangeChatSource() {
    var el = document.getElementById('chat-from');
    fromSerialNo = el.value;
    fromName = el.options[el.selectedIndex].innerHTML;
    document.getElementById('chat-send').disabled = (targetSerialNo=="" || fromSerialNo=="");
}

function ChangeChatLanguage() {
    chatLanguage = document.getElementById('chat-lang').value;
}