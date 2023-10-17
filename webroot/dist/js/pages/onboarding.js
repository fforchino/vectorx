LoadOnboardingPage();

async function LoadOnboardingPage() {
    LoadSettings().then(() => {
        if (document.getElementById("wirepod_console_url")!=null) {
            var setupUrl = Settings["WIREPOD_CONSOLE"];
            var hackUrl = Settings["WIREPOD_CONSOLE"]+"/api-chipper/use_ip?port=443";
            document.getElementById("wirepod_console_url").href = setupUrl;
            document.getElementById("wirepod_console_url").innerHTML = setupUrl;

            document.getElementById("wirepod_console_hack_url").href = hackUrl;
            document.getElementById("wirepod_console_hack_url").innerHTML = hackUrl;
        }
    });
}