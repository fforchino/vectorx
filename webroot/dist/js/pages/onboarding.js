(() => {
  "use strict";

  const steps = ["kind", "oskr", "discover", "pin", "wifi", "verify", "done"];
  let kind = "factory";
  let oskrPreparationSkipped = false;
  let selectedWifi = null;
  let robotsBefore = new Map();

  const el = id => document.getElementById(id);
  const sleep = ms => new Promise(resolve => setTimeout(resolve, ms));

  function showStep(name) {
    steps.forEach(step => { el(`step-${step}`).hidden = step !== name; });
    const labels = kind === "oskr" ? ["Type", oskrPreparationSkipped ? "Prepare (skipped)" : "Prepare", "Find", "PIN", "Wi-Fi", "Verify"] : ["Type", "Find", "PIN", "Wi-Fi", "Verify"];
    const order = kind === "oskr" ? ["kind", "oskr", "discover", "pin", "wifi", "verify"] : ["kind", "discover", "pin", "wifi", "verify"];
    const current = name === "done" ? labels.length : Math.max(0, order.indexOf(name));
    el("progress").innerHTML = labels.map((label, i) => `<span class="step-dot ${i <= current ? "active" : ""}">${i + 1}</span><small class="mr-2">${label}</small>`).join("");
  }

  function status(message, type = "info") {
    const box = el("status");
    box.hidden = false;
    box.className = `status-box mt-4 ${type}`;
    box.textContent = message;
  }

  function clearStatus() { el("status").hidden = true; }

  function responseArray(raw, field) {
    const parsed = JSON.parse(raw);
    if (Array.isArray(parsed)) return parsed;
    if (parsed && Array.isArray(parsed[field])) return parsed[field];
    return [];
  }

  async function api(action, fields, method = "POST") {
    const options = { method, headers: { "Accept": "text/plain" } };
    if (fields instanceof FormData) options.body = fields;
    else if (fields) {
      options.headers["Content-Type"] = "application/x-www-form-urlencoded;charset=UTF-8";
      options.body = new URLSearchParams(fields);
    }
    const response = await fetch(`/api/onboarding/${action}`, options);
    const text = (await response.text()).trim();
    if (!response.ok || !text || text.toLowerCase().startsWith("error")) throw new Error(friendlyError(text, response.status));
    return text;
  }

  function friendlyError(raw, code) {
    const value = raw.toLowerCase();
    if (value.includes("cannot find read channel")) return "Vector was found, but its setup Bluetooth channel is unavailable. Place it on the charger, double-press the back button until the PIN screen appears, then search again.";
    if (value.includes("robot status: connection error")) return "The PIN was accepted, but Vector could not reach escapepod.local. Check that the robot and VectorX are on the same network, then try again.";
    if (value.includes("certificate")) return "The server certificate does not match its address. VectorX must regenerate it before continuing.";
    if (value.includes("incorrect pin") || value.includes("eof")) return "The PIN is incorrect. Return Vector to the code screen and try again.";
    if (value.includes("unavailable") || code === 502) return "The onboarding service is not responding. Check that VectorX was installed completely.";
    if (value.includes("ssh") || value.includes("publickey")) return "VectorX cannot access the robot with this SSH key.";
    return raw || "The operation failed. You can safely try again.";
  }

  async function snapshotRobots() {
    try {
      const response = await fetch("/api/get_robots", { cache: "no-store" });
      const robots = await response.json();
      robotsBefore = new Map((Array.isArray(robots) ? robots : []).map(robot => [robot.esn, Boolean(robot.vector_settings)]));
    } catch (_) { robotsBefore = new Map(); }
  }

  document.querySelectorAll("[data-kind]").forEach(button => button.addEventListener("click", async () => {
    kind = button.dataset.kind;
    oskrPreparationSkipped = false;
    clearStatus();
    showStep(kind === "oskr" ? "oskr" : "discover");
    // Robot inventory can be slow while Wire-Pod is starting; never block
    // the first onboarding transition on that background snapshot.
    snapshotRobots();
  }));

  el("skip-oskr").addEventListener("click", () => {
    oskrPreparationSkipped = true;
    clearStatus();
    showStep("discover");
  });

  el("prepare-oskr").addEventListener("click", async () => {
    const ip = el("oskr-ip").value.trim();
    const key = el("oskr-key").files[0];
    if (!/^\d{1,3}(\.\d{1,3}){3}$/.test(ip) || !key) return status("Enter a valid IP address and select the robot's private key.", "error");
    const data = new FormData();
    data.append("ip", ip); data.append("key", key);
    try {
      status("OSKR preparation started. The display may turn off for a few minutes.");
      await api("oskr-setup", data);
      while (true) {
        await sleep(1000);
        const state = await api("oskr-status", null, "GET");
        if (state.includes("done")) break;
        if (state.includes("error")) throw new Error(friendlyError(state, 500));
        status(state.replace(/\.\.\.$/, "…"));
      }
      el("oskr-key").value = "";
      clearStatus(); showStep("discover");
    } catch (error) { status(error.message, "error"); }
  });

  el("scan").addEventListener("click", async () => {
    try {
      status("Initializing Bluetooth…");
      await api("init");
      status("Looking for nearby robots…");
      const devices = responseArray(await api("scan"), "devices");
      el("robots").innerHTML = "";
      if (!devices.length) return status("No Vector found. Check the robot's display and search again.", "error");
      devices.forEach(device => {
        const button = document.createElement("button");
        button.className = "btn btn-outline-light robot-result";
        button.textContent = device.name || `Vector ${device.address}`;
        button.addEventListener("click", () => connectRobot(device.id));
        el("robots").appendChild(button);
      });
      status("Select your Vector from the list.", "success");
    } catch (error) { status(error.message, "error"); }
  });

  async function connectRobot(id) {
    try {
      status("Connecting securely to Vector…");
      await api("connect", { id });
      clearStatus(); showStep("pin"); el("pin").focus();
    } catch (error) {
      await api("disconnect").catch(() => "");
      el("robots").innerHTML = "";
      status(error.message, "error");
    }
  }

  el("pin").addEventListener("input", event => { event.target.value = event.target.value.replace(/\D/g, "").slice(0, 6); });
  el("send-pin").addEventListener("click", async () => {
    const pin = el("pin").value;
    if (!/^\d{6}$/.test(pin)) return status("The PIN must contain exactly six digits.", "error");
    try {
      status("Verifying PIN…");
      await api("pin", { pin });
      const wifiState = await api("wifi-status");
      if (wifiState === "1") return authenticate();
      clearStatus(); showStep("wifi");
    } catch (error) { status(error.message, "error"); }
  });

  el("scan-wifi").addEventListener("click", async () => {
    try {
      status("Looking for Wi-Fi networks…");
      const networks = responseArray(await api("wifi-scan"), "networks");
      el("networks").innerHTML = "";
      networks.filter(network => network.ssid).forEach(network => {
        const button = document.createElement("button");
        button.className = "btn btn-outline-light robot-result"; button.textContent = network.ssid;
        button.addEventListener("click", () => { selectedWifi = network; el("selected-ssid").textContent = network.ssid; el("wifi-password").hidden = false; el("password").focus(); });
        el("networks").appendChild(button);
      });
      clearStatus();
    } catch (error) { status(error.message, "error"); }
  });

  el("connect-wifi").addEventListener("click", async () => {
    if (!selectedWifi || !el("password").value) return status("Select a network and enter its password.", "error");
    try {
      status("Connecting Vector to the Wi-Fi network…");
      const result = await api("wifi-connect", { ssid: selectedWifi.ssid, password: el("password").value, authType: selectedWifi.authtype });
      el("password").value = "";
      if (!result.includes("255")) throw new Error("Vector could not connect to the network. Check the password and try again.");
      await authenticate();
    } catch (error) { status(error.message, "error"); }
  });

  async function authenticate() {
    showStep("verify"); status("Authenticating the robot…");
    try {
      const result = await api("authenticate");
      if (result !== "true") throw new Error("Vector did not complete authentication.");
      await api("disconnect").catch(() => "");
      status("Credentials received. Verifying the VectorX connection…");
      await verifyOnline();
      clearStatus(); showStep("done");
    } catch (error) { status(error.message, "error"); }
  }

  async function verifyOnline() {
    for (let attempt = 0; attempt < 6; attempt++) {
      await sleep(5000);
      try {
        const response = await fetch("/api/get_robots", { cache: "no-store" });
        const robots = await response.json();
        // The robot may already be registered and online (for example after a
        // retry). In that case waiting for an offline -> online transition
        // would leave the wizard stuck forever.
        const online = (Array.isArray(robots) ? robots : []).some(robot => Boolean(robot.vector_settings));
        if (online) return;
      } catch (_) {}
    }
    throw new Error("The PIN was accepted, but VectorX still cannot authenticate the robot. You can retry verification without repeating the entire procedure.");
  }

  showStep("kind");
})();
