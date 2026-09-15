(function () {
  const esn = new URLSearchParams(window.location.search).get("esn");
  const status = document.getElementById("settings-status");
  document.getElementById("robot-esn").textContent = esn ? "· " + esn.toUpperCase() : "";
  LoadSite("nav_page_robots");

  function show(message, error) {
    status.className = "alert " + (error ? "alert-danger" : "alert-success");
    status.textContent = message;
  }
  async function send(action, values) {
    if (!esn) throw new Error("No robot selected");
    const body = new URLSearchParams(Object.assign({serial: esn}, values || {}));
    const response = await fetch("/api/robot-settings/" + action, {method: "POST", headers: {"Content-Type": "application/x-www-form-urlencoded;charset=UTF-8"}, body: body});
    const text = await response.text();
    if (!response.ok) throw new Error(text || "The robot did not accept the setting");
    return text;
  }
  async function change(action, values) {
    try { await send(action, values); show("Setting saved on the robot.", false); }
    catch (error) { show("Could not save: " + error.message, true); }
  }

  const volume = document.getElementById("volume");
  volume.addEventListener("input", () => document.getElementById("volume-value").textContent = volume.value);
  volume.addEventListener("change", () => change("volume", {volume: volume.value}));
  document.getElementById("locale").addEventListener("change", e => change("locale", {locale: e.target.value}));
  document.getElementById("eye-color").addEventListener("change", e => change("eye-color", {color: e.target.value}));
  document.getElementById("time-format").addEventListener("change", e => change(e.target.value));
  document.getElementById("temperature").addEventListener("change", e => change(e.target.value));
  document.getElementById("button-action").addEventListener("change", e => change(e.target.value));
  document.getElementById("save-location").addEventListener("click", () => change("location", {location: document.getElementById("location").value.trim()}));
  document.getElementById("save-timezone").addEventListener("click", () => change("timezone", {timezone: document.getElementById("timezone").value.trim()}));

  send("get").then(text => {
    const s = JSON.parse(text);
    volume.value = s.master_volume == null ? 3 : s.master_volume;
    document.getElementById("volume-value").textContent = volume.value;
    if (s.locale) document.getElementById("locale").value = s.locale;
    if (s.eye_color != null) document.getElementById("eye-color").value = String(s.eye_color);
    document.getElementById("location").value = s.default_location || "";
    document.getElementById("timezone").value = s.time_zone || "";
    document.getElementById("time-format").value = s.clock_24_hour ? "time-24" : "time-12";
    document.getElementById("temperature").value = s.temp_is_fahrenheit ? "temp-f" : "temp-c";
    document.getElementById("button-action").value = s.button_wakeword === "alexa" ? "button-alexa" : "button-vector";
    show("Settings loaded. Changes are applied immediately.", false);
  }).catch(error => show("Could not read the robot settings: " + error.message, true));
})();
