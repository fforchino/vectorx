var CurrentPage = "";

function ReloadSite() {
  LoadSite(CurrentPage);
}

async function LoadSite(selectedPage) {
  CurrentPage = selectedPage;
    LoadNavBar();
    LoadBrandLogo();
    LoadSidebar(selectedPage);
    LoadSettings().then(() => {
      checkSetupMissing().then(() => {
        LoadRobots().then(() => {
          SidebarGetRobotList();
          if (selectedPage == "nav_page_home") {
            LoadHomePageBots();
          } else if (selectedPage == "nav_page_help") {
            LoadVectorXCustomIntents().then(() => {
              LoadHelpPage();
            })
          } else if (selectedPage.startsWith("nav_page_botcontrol")) {
            LoadVoices(Settings.STT_LANGUAGE).then(() => {
              LoadBotControlPage();
            });
          } else {
            // Handle selection for normal pages
            document.getElementById(selectedPage).classList.add("active");
          }
          LoadFooter();
          if (document.getElementById("wirepod_console_url")!=null) {
            document.getElementById("wirepod_console_url").href = getWirePodConsoleUrl();
          }

          /*
          LoadIntents().then(() => {
            SidebarGetRobotList();
            if (selectedPage == "nav_page_home") {
              LoadHomePageBots();
            }
            SidebarGetIntentList();
            if (selectedPage.startsWith("nav_page_intent_edit")) {
              LoadEditIntentPage();
            }
            // Handle selection
            document.getElementById(selectedPage).classList.add("active");
            if (selectedPage=="nav_page_intent_add" || selectedPage.startsWith("nav_page_intent_edit")) {
              document.getElementById("nav_group_custom_intents").classList.add("active");
            }
          });
           */
      });
    });
  });
}

function getWirePodConsoleUrl() {
  return Settings["WIREPOD_CONSOLE"];
}

async function checkSetupMissing() {
  await fetch("/api/is_setup_done")
      .then(response => response.text())
      .then((response) => {
        try {
          obj = JSON.parse(response);
          if (obj.result=="OK") {
            if (CurrentPage == "nav_page_home") {
              doConsistencyCheck()
            }
          }
          else {
            goInitialSetup();
          }
        } catch { goInitialSetup(); }
      })
}

function LoadNavBar() {
  var data = '    <!-- Left navbar links -->\n' +
    '    <ul class="navbar-nav">\n' +
    '      <li class="nav-item">\n' +
    '        <a class="nav-link" data-widget="pushmenu" href="#" role="button"><i class="fas fa-bars"></i></a>\n' +
    '      </li>\n' +
    '      <li class="nav-item d-none d-sm-inline-block">\n' +
    '        <a href="index.html" class="nav-link">Home</a>\n' +
    '      </li>\n' +
    '      <li class="nav-item d-none d-sm-inline-block">\n' +
    '        <a href="https://github.com/fforchino/vectorx" target="_blank" class="nav-link">Contact</a>\n' +
    '      </li>\n' +
    '    </ul>\n' +
    '\n' +
    '    <!-- Right navbar links -->\n' +
    '    <ul class="navbar-nav ml-auto">\n' +
    '      <!-- Navbar Search -->\n' +
    '      <li class="nav-item">\n' +
    '      </li>\n' +
    '\n' +
    '      <!-- Messages Dropdown Menu -->\n' +
    '      <li class="nav-item">\n' +
    '        <a class="nav-link" data-widget="fullscreen" href="#" role="button">\n' +
    '          <i class="fas fa-expand-arrows-alt"></i>\n' +
    '        </a>\n' +
    '      </li>\n' +
    '      <li class="nav-item">\n' +
    '        <a class="nav-link" data-widget="control-sidebar" data-slide="true" href="#" role="button">\n' +
    '          <i class="fas fa-th-large"></i>\n' +
    '        </a>\n' +
    '      </li>\n' +
    '    </ul>\n';
  document.getElementById("navbar").innerHTML = data;
}

function LoadBrandLogo() {
  var data = '    <a href="index.html" class="brand-link">\n' +
    '      <img src="dist/img/VectorXLogo.png" alt="VectorX Logo" class="brand-image img-circle elevation-3" style="opacity: .8">\n' +
    '      <span class="brand-text font-weight-light">VectorX</span>\n' +
    '    </a>\n';
  document.getElementById("brand-logo").innerHTML = data;
}
