
async function LoadOnboardingPage() {
    LoadSettings().then(() => {
        checkSetupMissing().then(() => {
            if (document.getElementById("wirepod_console_url")!=null) {
                var setupUrl = Settings["WIREPOD_CONSOLE"];
                document.getElementById("wirepod_console_url").href = setupUrl;
                document.getElementById("wirepod_console_url").innerHTML = setupUrl;
            }
        });
    });
}