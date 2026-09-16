(function () {
  "use strict";
  const esn = new URLSearchParams(window.location.search).get("esn");
  const status = document.getElementById("settings-status");
  const photoURLs = [];
  let stimTimer = null;
  let stimValues = [];

  document.getElementById("robot-esn").textContent = esn ? "· " + esn.toUpperCase() : "";
  document.getElementById("vector-control").href = "vector-control.html?esn=" + encodeURIComponent(esn || "");
  LoadSite("nav_page_robots");

  function message(text, error) {
    status.className = "alert " + (error ? "alert-danger" : "alert-success");
    status.textContent = text;
  }

  async function request(action, values, type) {
    if (!esn) throw new Error("No robot selected");
    const body = new URLSearchParams(Object.assign({serial: esn}, values || {}));
    const response = await fetch("/api/robot-settings/" + action, {
      method: "POST", headers: {"Content-Type": "application/x-www-form-urlencoded;charset=UTF-8"}, body
    });
    if (!response.ok) throw new Error((await response.text()) || "The robot did not accept the request");
    if (type === "blob") return response.blob();
    return response.text();
  }

  async function change(action, values, success) {
    try { await request(action, values); message(success || "Setting saved on the robot.", false); return true; }
    catch (error) { message("Could not complete the request: " + error.message, true); return false; }
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

  const colorInput = document.getElementById("custom-eye-color");
  colorInput.addEventListener("input", () => document.getElementById("eye-preview").style.background = colorInput.value);
  document.getElementById("save-eye-color").addEventListener("click", async () => {
    const hsv = hexToHSV(colorInput.value);
    if (await change("custom-eye-color", {hue: (hsv.h / 360).toFixed(3), sat: (hsv.s / 100).toFixed(3)}, "Custom eye color applied.")) {
      document.getElementById("eye-color").selectedIndex = -1;
    }
  });

  document.querySelectorAll(".quick-action").forEach(button => button.addEventListener("click", async () => {
    button.classList.add("busy");
    await change("quick-action", {intent: button.dataset.intent}, "Action sent to Vector.");
    button.classList.remove("busy");
  }));

  document.getElementById("alexa-sign-in").addEventListener("click", () => change("alexa-sign-in", {}, "Alexa sign-in started. Follow the instructions on Vector, then enter the code at amazon.com/code."));
  document.getElementById("alexa-sign-out").addEventListener("click", () => change("alexa-sign-out", {}, "Vector has been signed out of Alexa."));

  async function loadFaces() {
    const select = document.getElementById("face-list");
    select.innerHTML = "<option>Loading…</option>";
    try {
      const raw = await request("faces");
      const faces = raw.trim() && raw.trim() !== "null" ? JSON.parse(raw) : [];
      select.innerHTML = "";
      faces.forEach(face => {
        const option = document.createElement("option"); option.value = String(face.face_id); option.textContent = face.name; option.dataset.name = face.name; select.appendChild(option);
      });
      if (!faces.length) select.innerHTML = "<option value=''>No learned faces</option>";
      document.getElementById("rename-face").disabled = document.getElementById("delete-face").disabled = !faces.length;
    } catch (error) { select.innerHTML = "<option value=''>Could not load faces</option>"; message("Could not load faces: " + error.message, true); }
  }
  document.getElementById("refresh-faces").addEventListener("click", loadFaces);
  document.getElementById("rename-face").addEventListener("click", async () => {
    const select = document.getElementById("face-list"), option = select.selectedOptions[0]; if (!option || !option.value) return;
    const newName = window.prompt("Enter the new name:", option.dataset.name); if (newName === null) return;
    if (!newName.trim()) { message("The face name cannot be empty.", true); return; }
    if (await change("face-rename", {id: option.value, oldname: option.dataset.name, newname: newName.trim()}, "Face renamed.")) loadFaces();
  });
  document.getElementById("delete-face").addEventListener("click", async () => {
    const select = document.getElementById("face-list"), option = select.selectedOptions[0]; if (!option || !option.value) return;
    if (!window.confirm("Remove " + option.dataset.name + " from Vector's learned faces?")) return;
    if (await change("face-delete", {id: option.value}, "Face removed.")) loadFaces();
  });

  function drawStim() {
    const canvas = document.getElementById("stim-chart"), ctx = canvas.getContext("2d"), w = canvas.width, h = canvas.height;
    ctx.clearRect(0, 0, w, h); ctx.strokeStyle = "#495057"; ctx.lineWidth = 1;
    for (let i = 1; i < 4; i++) { ctx.beginPath(); ctx.moveTo(0, h * i / 4); ctx.lineTo(w, h * i / 4); ctx.stroke(); }
    if (stimValues.length < 2) return;
    ctx.strokeStyle = "#33ed6d"; ctx.lineWidth = 3; ctx.beginPath();
    stimValues.forEach((value, i) => { const x = i * w / 29, y = h - Math.max(0, Math.min(1, value)) * h; i ? ctx.lineTo(x, y) : ctx.moveTo(x, y); }); ctx.stroke();
  }
  async function stopStim(silent) {
    if (stimTimer) { clearInterval(stimTimer); stimTimer = null; }
    document.getElementById("start-stim").disabled = false; document.getElementById("stop-stim").disabled = true;
    try { await request("stim-stop"); if (!silent) message("Live stimulation view stopped.", false); } catch (error) { if (!silent) message("Could not stop the live view: " + error.message, true); }
  }
  document.getElementById("start-stim").addEventListener("click", async () => {
    try {
      await request("stim-start"); stimValues = []; drawStim();
      document.getElementById("start-stim").disabled = true; document.getElementById("stop-stim").disabled = false;
      stimTimer = setInterval(async () => { try { const value = Number(await request("stim-status")); if (!Number.isFinite(value)) return; stimValues.push(value); if (stimValues.length > 30) stimValues.shift(); document.getElementById("stim-value").textContent = value.toFixed(2); drawStim(); } catch (_) {} }, 500);
      message("Live stimulation view started.", false);
    } catch (error) { message("Could not start the live view: " + error.message, true); }
  });
  document.getElementById("stop-stim").addEventListener("click", () => stopStim(false));

  function clearPhotoURLs() { photoURLs.splice(0).forEach(URL.revokeObjectURL); }
  async function loadPhotos() {
    const grid = document.getElementById("photo-grid"); grid.innerHTML = "<div class='empty-state'>Loading photos…</div>"; clearPhotoURLs();
    try {
      const raw = await request("photo-ids"), ids = raw.trim() && raw.trim() !== "null" ? JSON.parse(raw) : [];
      grid.innerHTML = "";
      if (!ids.length) { grid.innerHTML = "<div class='empty-state'>No photos found. Ask Vector to take a photo, then refresh.</div>"; return; }
      await Promise.all(ids.map(async id => {
        const blob = await request("photo-thumb", {id}, "blob"), url = URL.createObjectURL(blob); photoURLs.push(url);
        const card = document.createElement("div"); card.className = "photo-card";
        const img = document.createElement("img"); img.src = url; img.alt = "Vector photo " + id; img.title = "Open full-size photo";
        img.addEventListener("click", async () => { try { const full = await request("photo", {id}, "blob"), fullURL = URL.createObjectURL(full); photoURLs.push(fullURL); window.open(fullURL, "_blank", "noopener"); } catch (error) { message("Could not open the photo: " + error.message, true); } });
        const remove = document.createElement("button"); remove.className = "btn btn-sm btn-danger btn-block"; remove.textContent = "Delete";
        remove.addEventListener("click", async () => { if (window.confirm("Delete this photo from Vector?") && await change("photo-delete", {id}, "Photo deleted.")) loadPhotos(); });
        card.append(img, remove); grid.appendChild(card);
      }));
    } catch (error) { grid.innerHTML = "<div class='empty-state'>Could not load photos.</div>"; message("Could not load photos: " + error.message, true); }
  }
  document.getElementById("refresh-photos").addEventListener("click", loadPhotos);

  function hexToHSV(hex) {
    const r = parseInt(hex.slice(1, 3), 16) / 255, g = parseInt(hex.slice(3, 5), 16) / 255, b = parseInt(hex.slice(5, 7), 16) / 255;
    const max = Math.max(r, g, b), min = Math.min(r, g, b), d = max - min; let h = 0;
    if (d) { if (max === r) h = 60 * (((g - b) / d) % 6); else if (max === g) h = 60 * ((b - r) / d + 2); else h = 60 * ((r - g) / d + 4); }
    if (h < 0) h += 360; return {h, s: max ? d / max * 100 : 0};
  }
  function hsvToHex(h, s, v) {
    s /= 100; v /= 100; const c = v * s, x = c * (1 - Math.abs((h / 60) % 2 - 1)), m = v - c; let rgb;
    if (h < 60) rgb = [c,x,0]; else if (h < 120) rgb = [x,c,0]; else if (h < 180) rgb = [0,c,x]; else if (h < 240) rgb = [0,x,c]; else if (h < 300) rgb = [x,0,c]; else rgb = [c,0,x];
    return "#" + rgb.map(n => Math.round((n + m) * 255).toString(16).padStart(2,"0")).join("");
  }

  request("get").then(raw => {
    const s = JSON.parse(raw); volume.value = s.master_volume == null ? 3 : s.master_volume; document.getElementById("volume-value").textContent = volume.value;
    if (s.locale) document.getElementById("locale").value = s.locale; if (s.eye_color != null) document.getElementById("eye-color").value = String(s.eye_color);
    if (s.custom_eye_color && s.custom_eye_color.enabled) { colorInput.value = hsvToHex(Number(s.custom_eye_color.hue) * 360, Number(s.custom_eye_color.saturation) * 100, 100); document.getElementById("eye-color").selectedIndex = -1; }
    document.getElementById("eye-preview").style.background = colorInput.value; document.getElementById("location").value = s.default_location || ""; document.getElementById("timezone").value = s.time_zone || "";
    document.getElementById("time-format").value = s.clock_24_hour ? "time-24" : "time-12"; document.getElementById("temperature").value = s.temp_is_fahrenheit ? "temp-f" : "temp-c";
    document.getElementById("button-action").value = Number(s.button_wakeword) === 1 || s.button_wakeword === "alexa" ? "button-alexa" : "button-vector";
    message("Settings loaded. Changes are applied immediately.", false);
  }).catch(error => message("Could not read the robot settings: " + error.message, true));
  loadFaces();
  window.addEventListener("beforeunload", () => { clearPhotoURLs(); if (stimTimer) request("stim-stop").catch(() => {}); });
})();
