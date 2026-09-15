function LoadSidebar(selectedPage) {
  var data = '      <!-- Sidebar Menu -->\n' +
    '      <nav class="mt-2">\n' +
    '        <ul class="nav nav-pills nav-sidebar flex-column" data-widget="treeview" role="menu" data-accordion="false">\n' +
    '          <!-- Add icons to the links using the .nav-icon class\n' +
    '               with font-awesome or any other icon font library -->\n' +
    '          <li class="nav-item menu-open">\n' +
    '            <a href="#" class="nav-link">\n' +
    '              <i class="nav-icon fas fa-tachometer-alt"></i>\n' +
    '              <p>\n' +
    '                VectorX\n' +
    '                <i class="right fas fa-angle-left"></i>\n' +
    '              </p>\n' +
    '            </a>\n' +
    '            <ul class="nav nav-treeview">\n' +
    '              <li class="nav-item">\n' +
    '                <a id="nav_page_home" href="index.html" class="nav-link">\n' +
    '                  <i class="fas fa-home nav-icon"></i>\n' +
    '                  <p>Home</p>\n' +
    '                </a>\n' +
    '              </li>\n' +
    '              <li class="nav-item">\n' +
    '                <a href="#" class="nav-link">\n' +
    '                  <i class="nav-icon fas fa-gear""></i>\n' +
    '                  <p>\n' +
    '                    Setup\n' +
    '                    <i class="right fas fa-angle-left"></i>\n' +
    '                  </p>\n' +
    '                </a>\n' +
    '                <ul class="nav nav-treeview">\n' +
    '                  <li class="nav-item">\n' +
    '                    <a href="initial_setup.html" class="nav-link">\n' +
    '                      <i class="fas fa-hat-wizard nav-icon text-sm"></i>\n' +
    '                      <p>Run Setup Wizard</p>\n' +
    '                    </a>\n' +
    '                  </li>\n' +
    '                </ul>\n' +
    '              </li>\n' +
    '              <li class="nav-item" style="display:none">\n' +
    '                <a id="nav_group_custom_intents" href="#" class="nav-link">\n' +
    '                  <i class="nav-icon fas fa-microphone-lines""></i>\n' +
    '                  <p>\n' +
    '                    Custom Intents\n' +
    '                    <i class="right fas fa-angle-left"></i>\n' +
    '                  </p>\n' +
    '                </a>\n' +
    '                <ul id="sidebar_intent_list" class="nav nav-treeview">\n' +
    '                </ul>\n' +
    '              </li>\n' +
    '              <li class="nav-item">\n' +
    '                <a href="#" class="nav-link" id="nav_page_robots">\n' +
    '                  <i class="nav-icon fas fa-robot"></i>\n' +
    '                  <p>\n' +
    '                    Robots\n' +
    '                    <i class="right fas fa-angle-left"></i>\n' +
    '                  </p>\n' +
    '                </a>\n' +
    '                <ul id="sidebar_robot_list" class="nav nav-treeview">\n' +
    '                </ul>\n' +
    '                <ul class="nav nav-treeview">\n' +
    '                  <li class="nav-item">\n' +
    '                    <a id="nav_page_add_robot" href="onboarding.html" class="nav-link">\n' +
    '                      <i class="fas fa-plus-circle nav-icon text-sm"></i>\n' +
    '                      <p>Add robot</p>\n' +
    '                    </a>\n' +
    '                  </li>\n' +
    '                </ul>\n' +
    '              </li>\n' +
    '              <li class="nav-item">\n' +
    '                <a id="nav_page_update" href="chat.html" class="nav-link">\n' +
    '                  <i class="fas fa-comments nav-icon"></i>\n' +
    '                  <p>Chat</p>\n' +
    '                </a>\n' +
    '              </li>\n' +
    '              <li class="nav-item">\n' +
    '                <a id="nav_page_update" href="update.html" class="nav-link">\n' +
    '                  <i class="fas fa-download nav-icon"></i>\n' +
    '                  <p>Check for updates</p>\n' +
    '                </a>\n' +
    '              </li>\n' +
    '              <li class="nav-item">\n' +
    '                <a id="nav_page_help" href="voicecommands.html" class="nav-link">\n' +
    '                  <i class="fas fa-life-ring nav-icon"></i>\n' +
    '                  <p>Voice Command Help</p>\n' +
    '                </a>\n' +
    '              </li>\n' +
    '<!--'+
    '              <li class="nav-item">\n' +
    '                <a href="./index.html" class="nav-link">\n' +
    '                  <i class="fas fa-file-lines nav-icon"></i>\n' +
    '                  <p>Log</p>\n' +
    '                </a>\n' +
    '              </li>\n' +
    '-->'+
    '            </ul>\n' +
    '          </li>\n' +
    '        </ul>\n' +
    '      </nav>\n' +
    '      <!-- /.sidebar-menu -->\n';
  document.getElementById("sidebar").innerHTML = data;
}

function SidebarGetRobotList() {
  var data = "";
  for (var i = 0; i < Robots.length; i++) {
    var bot = Robots[i];
    if (bot.vector_settings!=null) {
    var botName = bot.custom_settings.RobotName.toUpperCase();
    if (botName.length==0) botName = bot.esn.toUpperCase();
    data +=
          '                  <li class="nav-item">\n' +
          '                    <a id="nav_page_botcontrol_'+bot.esn+'" href="botcontrol.html?esn='+bot.esn+'" class="nav-link">\n' +
          '                      <i class="fas fa-square nav-icon text-sm"></i>\n' +
          '                      <p>'+botName+'</p>\n' +
          '                    </a>\n' +
          '                    <a href="robot-settings.html?esn='+bot.esn+'" class="nav-link pl-4">\n' +
          '                      <i class="fas fa-sliders-h nav-icon text-sm"></i>\n' +
          '                      <p>Settings</p>\n' +
          '                    </a>\n' +
          '                    <a href="#" class="nav-link text-danger pl-4" onclick="RemoveRobot(\''+bot.esn+'\', \''+botName.replaceAll("'", "")+'\'); return false;">\n' +
          '                      <i class="fas fa-unlink nav-icon text-sm"></i>\n' +
          '                      <p>Remove '+botName+'</p>\n' +
          '                    </a>\n' +
          '                  </li>\n';
    }
  }
  document.getElementById("sidebar_robot_list").innerHTML = data;
}

async function RemoveRobot(esn, botName) {
  if (!confirm("Remove " + botName + " from VectorX? The robot can be added again later.")) return;
  try {
    const response = await fetch("/api/remove_robot", {
      method: "POST",
      headers: {"Content-Type": "application/x-www-form-urlencoded;charset=UTF-8"},
      body: new URLSearchParams({esn: esn})
    });
    const result = await response.json();
    if (!response.ok || result.result !== "OK") throw new Error(result.reason || "Remove failed");
    window.location.replace("index.html");
  } catch (error) {
    alert("VectorX could not remove the robot: " + error.message);
  }
}

function SidebarGetIntentList() {
  var data = "";
  for (var i = 0; i < Intents.length; i++) {
    var intent = Intents[i];
    if (!intent.issystem) {
      data +=
          '                  <li class="nav-item">\n' +
          '                    <a id="nav_page_intent_edit_'+intent.name+'" href="intent-edit.html?id='+intent.name+'" class="nav-link">\n' +
          '                      <i class="fas fa-square nav-icon text-sm"></i>\n' +
          '                      <p>'+intent.name+'</p>\n' +
          '                    </a>\n' +
          '                  </li>\n';
    }
  }
  data+=
      '                  <li class="nav-item">\n' +
      '                    <a id="nav_page_intent_add" href="intent-add.html" class="nav-link">\n' +
      '                      <i class="fas fa-plus nav-icon"></i>\n' +
      '                      <p>Add New</p>\n' +
      '                    </a>\n' +
      '                  </li>\n';
  document.getElementById("sidebar_intent_list").innerHTML = data;
}
