(function () {
  "use strict";
  const esn = new URLSearchParams(window.location.search).get("esn");
  const status = document.getElementById("control-status");
  let controlHeld = false, cameraRunning = false, keyboardEnabled = false;
  const pressed = new Set();
  let lastWheels = "", lastLift = null, lastHead = null;
  document.getElementById("robot-esn").textContent = esn ? esn.toUpperCase() : "No robot selected";
  document.getElementById("back-settings").href = "robot-settings.html?esn=" + encodeURIComponent(esn || "");
  LoadSite("nav_page_robots");

  function show(text, error) { status.className = "alert " + (error ? "alert-danger" : "alert-success"); status.textContent = text; }
  async function command(action, values, keepalive) {
    if (!esn) throw new Error("No robot selected");
    const response = await fetch("/api/robot-settings/" + action, {method:"POST", headers:{"Content-Type":"application/x-www-form-urlencoded;charset=UTF-8"}, body:new URLSearchParams(Object.assign({serial:esn}, values || {})), keepalive:Boolean(keepalive)});
    const text = await response.text(); if (!response.ok || /^error\b/i.test(text.trim())) throw new Error(text || "Vector rejected the command"); return text;
  }
  function setControl(enabled) {
    controlHeld = enabled; document.getElementById("control-state").textContent = enabled ? "Active" : "Released"; document.getElementById("control-dot").classList.toggle("active", enabled);
    document.getElementById("assume-control").disabled = enabled; document.getElementById("release-control").disabled = !enabled; document.getElementById("motor-controls").classList.toggle("control-disabled", !enabled);
    document.querySelectorAll(".controlled").forEach(el => el.disabled = !enabled);
  }
  document.getElementById("assume-control").addEventListener("click", async () => { try { await command("behavior-assume", {priority:"high"}); setControl(true); show("Behavior control granted.", false); } catch(e) { show("Could not assume behavior control: " + e.message, true); } });
  document.getElementById("release-control").addEventListener("click", async () => { await stopAll(); try { await command("behavior-release"); setControl(false); show("Behavior control released.", false); } catch(e) { show("Could not release behavior control: " + e.message, true); } });

  async function wheels(lw, rw) { const key = lw+":"+rw; if (key === lastWheels) return; lastWheels = key; try { await command("move-wheels", {lw, rw}); } catch(e) { show("Drive command failed: "+e.message, true); } }
  async function axis(name, speed) { if (name === "lift") { if (lastLift === speed) return; lastLift = speed; } else { if (lastHead === speed) return; lastHead = speed; } try { await command("move-"+name, {speed}); } catch(e) { show("Motor command failed: "+e.message, true); } }
  async function stopAll() { lastWheels = lastLift = lastHead = null; await Promise.allSettled([command("move-wheels",{lw:0,rw:0}), command("move-lift",{speed:0}), command("move-head",{speed:0})]); }
  document.querySelectorAll(".motor").forEach(btn => { const start=e=>{e.preventDefault();if(controlHeld)wheels(btn.dataset.lw,btn.dataset.rw)}; const stop=e=>{e.preventDefault();if(controlHeld){lastWheels="";wheels(0,0)}}; btn.addEventListener("pointerdown",start); btn.addEventListener("pointerup",stop); btn.addEventListener("pointercancel",stop); btn.addEventListener("pointerleave",stop); });
  document.getElementById("drive-stop").addEventListener("click",()=>{lastWheels="";wheels(0,0)});
  document.querySelectorAll(".motor-axis").forEach(btn => { const start=e=>{e.preventDefault();if(controlHeld)axis(btn.dataset.axis,btn.dataset.speed)}; const stop=e=>{e.preventDefault();if(controlHeld){if(btn.dataset.axis==="lift")lastLift=null;else lastHead=null;axis(btn.dataset.axis,0)}}; btn.addEventListener("pointerdown",start); btn.addEventListener("pointerup",stop); btn.addEventListener("pointercancel",stop); btn.addEventListener("pointerleave",stop); });

  document.getElementById("keyboard-toggle").addEventListener("click", () => { keyboardEnabled = !keyboardEnabled; pressed.clear(); stopAll(); document.getElementById("keyboard-toggle").textContent = keyboardEnabled ? "Disable keyboard controls" : "Enable keyboard controls"; document.getElementById("keyboard-hint").classList.toggle("d-none", !keyboardEnabled); });
  function updateKeys() {
    if (!keyboardEnabled || !controlHeld) return;
    const w=pressed.has("w"),a=pressed.has("a"),s=pressed.has("s"),d=pressed.has("d"); let lw=0,rw=0;
    if(w){lw=140;rw=140} if(s){lw=-150;rw=-150} if(a){lw-=150;rw+=150} if(d){lw+=150;rw-=150} wheels(Math.max(-250,Math.min(250,lw)),Math.max(-250,Math.min(250,rw)));
    axis("lift",pressed.has("r")?2:pressed.has("f")?-2:0); axis("head",pressed.has("t")?2:pressed.has("g")?-2:0);
  }
  document.addEventListener("keydown",e=>{if(!keyboardEnabled||e.repeat||!["w","a","s","d","r","f","t","g"].includes(e.key.toLowerCase()))return;e.preventDefault();pressed.add(e.key.toLowerCase());updateKeys()});
  document.addEventListener("keyup",e=>{if(!keyboardEnabled)return;pressed.delete(e.key.toLowerCase());updateKeys()});
  window.addEventListener("blur",()=>{pressed.clear();if(controlHeld)stopAll()});

  document.getElementById("say-button").addEventListener("click", async () => { const text=document.getElementById("say-text").value.trim(); if(!text){show("Enter some text first.",true);return} try{await command("say-text",{text});show("Text sent to Vector.",false)}catch(e){show("Could not make Vector speak: "+e.message,true)} });
  document.getElementById("say-text").addEventListener("keydown",e=>{if(e.key==="Enter")document.getElementById("say-button").click()});
  document.getElementById("mirror-on").addEventListener("click",()=>command("mirror",{enable:"true"}).then(()=>show("Mirror mode enabled.",false)).catch(e=>show("Could not enable mirror mode: "+e.message,true)));
  document.getElementById("mirror-off").addEventListener("click",()=>command("mirror",{enable:"false"}).then(()=>show("Mirror mode disabled.",false)).catch(e=>show("Could not disable mirror mode: "+e.message,true)));

  document.getElementById("camera-start").addEventListener("click",()=>{if(cameraRunning||!esn)return;const img=document.createElement("img");img.id="camera-image";img.alt="Live feed from Vector";img.onload=()=>show("Live camera feed started.",false);img.onerror=()=>show("The live camera feed was interrupted. Stop it and try again.",true);cameraRunning=true;document.getElementById("camera-start").disabled=true;document.getElementById("camera-stop").disabled=false;document.getElementById("camera-frame").replaceChildren(img);img.src="/api/robot-camera?serial="+encodeURIComponent(esn)+"&t="+Date.now();show("Starting live camera feed…",false)});
  async function stopCamera(silent) { cameraRunning=false;const img=document.getElementById("camera-image");if(img){img.removeAttribute("src");img.remove()}document.getElementById("camera-frame").innerHTML='<span id="camera-placeholder" class="text-muted">Camera is off</span>';document.getElementById("camera-start").disabled=false;document.getElementById("camera-stop").disabled=true;try{await command("camera-stop",{},silent);if(!silent)show("Live camera feed stopped.",false)}catch(e){if(!silent)show("Could not stop the camera: "+e.message,true)} }
  document.getElementById("camera-stop").addEventListener("click",()=>stopCamera(false));
  window.addEventListener("beforeunload",()=>{if(cameraRunning)command("camera-stop",{},true).catch(()=>{});if(controlHeld){command("move-wheels",{lw:0,rw:0},true).catch(()=>{});command("behavior-release",{},true).catch(()=>{})}});
  setControl(false);
})();
