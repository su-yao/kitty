const state = {
  category: "resource",
  page: 1,
  pageSize: 40,
  version: "4.4.6",
  totalPages: 1,
  cols: 8,
  rows: 5,
  profile: "27-2k",
  scale: 100,
  sending: false,
};

const $ = (id) => document.getElementById(id);

function setStatus(text, err) {
  const el = $("status");
  el.textContent = text;
  el.className = "status" + (err ? " err" : "");
}

function detectProfile() {
  const dpr = window.devicePixelRatio || 1;
  const cssW = window.screen.width;
  const cssH = window.screen.height;
  const physW = Math.round(cssW * dpr);
  const scale = dpr >= 1.35 ? 150 : 100;
  let profile = "27-2k";
  let label = "27寸 2K";
  if (physW >= 3000) {
    profile = "31.5-4k";
    label = "31.5寸 4K";
  }
  return { profile, label, scale, dpr, cssW, cssH, physW };
}

function tileMin(profile, scale) {
  if (profile === "31.5-4k" && scale === 100) return { w: 148, h: 58 };
  if (profile === "31.5-4k" && scale === 150) return { w: 158, h: 64 };
  if (profile === "27-2k" && scale === 150) return { w: 168, h: 70 };
  return { w: 162, h: 66 };
}

function fitGrid() {
  const info = detectProfile();
  const grid = $("grid");
  const min = tileMin(info.profile, info.scale);
  const cs = getComputedStyle(document.documentElement);
  const gap = parseFloat(cs.getPropertyValue("--gap")) || 8;
  const w = Math.max(grid.clientWidth, 640);
  const h = Math.max(grid.clientHeight, 280);
  let cols = Math.max(4, Math.floor((w + gap) / (min.w + gap)));
  let rows = Math.max(3, Math.floor((h + gap) / (min.h + gap)));
  // Keep the client page size identical to the server cap, otherwise the two
  // sides disagree on page boundaries and items can repeat across pages.
  const maxPage = 800;
  while (cols * rows > maxPage && rows > 3) rows -= 1;
  while (cols * rows > maxPage && cols > 4) cols -= 1;
  const pageSize = cols * rows;
  const changed = state.cols !== cols || state.rows !== rows || state.pageSize !== pageSize;
  state.cols = cols;
  state.rows = rows;
  state.pageSize = pageSize;
  state.profile = info.profile;
  state.scale = info.scale;
  document.documentElement.dataset.profile = info.profile;
  document.documentElement.dataset.scale = String(info.scale);
  grid.style.setProperty("--cols", String(cols));
  grid.style.setProperty("--rows", String(rows));
  return { changed, info, cols, rows, pageSize };
}

async function api(path, opts) {
  const res = await fetch(path, opts);
  const data = await res.json().catch(() => ({}));
  if (!res.ok) throw new Error(data.error || res.statusText);
  return data;
}

async function loadMeta() {
  const data = await api("/api/meta");
  const sel = $("version");
  sel.innerHTML = "";
  (data.versions || []).forEach((v) => {
    const o = document.createElement("option");
    o.value = v;
    o.textContent = v;
    sel.appendChild(o);
  });
  sel.value = data.default_version || "4.4.6";
  state.version = sel.value;
  const cats = $("cats");
  cats.innerHTML = "";
  const all = document.createElement("button");
  all.textContent = "全部";
  all.dataset.id = "all";
  cats.appendChild(all);
  (data.categories || []).forEach((c) => {
    const b = document.createElement("button");
    b.textContent = c.name;
    b.dataset.id = c.id;
    if (c.id === state.category) b.classList.add("active");
    cats.appendChild(b);
  });
  if (data.version) {
    const foot = $("appVersion");
    if (foot) foot.textContent = "v" + data.version;
  }
  if (data.settings) {
    if (data.settings.game_root) $("gameRoot").value = data.settings.game_root;
    if (data.settings.workshop) $("workshop").value = data.settings.workshop;
    if (data.settings.extra_mods) $("extraMods").value = data.settings.extra_mods;
    if (data.settings.console_key) $("consoleKey").value = data.settings.console_key;
  }
}

async function loadList() {
  fitGrid();
  const q = encodeURIComponent($("q").value.trim());
  const url = `/api/commands?category=${encodeURIComponent(state.category)}&q=${q}&version=${encodeURIComponent(state.version)}&page=${state.page}&page_size=${state.pageSize}`;
  const data = await api(url);
  state.page = data.page;
  state.totalPages = data.total_pages;
  const info = detectProfile();
  $("pageInfo").textContent = `${data.page} / ${data.total_pages}  共 ${data.total} 条`;
  $("meta").textContent = `${info.label} ${info.scale}% · ${state.cols}×${state.rows} 平铺 ${state.pageSize} 键 · 版本 ${state.version} · 扫描 ${data.scanned || 0}`;
  const grid = $("grid");
  grid.innerHTML = "";
  if (!(data.items || []).length) {
    const empty = document.createElement("div");
    empty.className = "empty";
    empty.textContent = "没有匹配的命令，换个分类或关键词试试。";
    grid.appendChild(empty);
  }
  (data.items || []).forEach((item) => {
    const b = document.createElement("button");
    b.className = "cmd";
    b.title = `${item.name_zh}\n${item.command}\n${item.name_en || ""} ${item.note || ""}`.trim();
    b.innerHTML = `<b>${escapeHtml(item.name_zh)}</b><code>${escapeHtml(item.command)}</code>`;
    b.addEventListener("click", () => sendCmd(item.command, item.name_zh));
    grid.appendChild(b);
  });
  document.querySelectorAll(".cats button").forEach((b) => {
    b.classList.toggle("active", b.dataset.id === state.category);
  });
}

function escapeHtml(s) {
  return String(s || "").replace(/[&<>"']/g, (c) => ({ "&": "&amp;", "<": "&lt;", ">": "&gt;", '"': "&quot;", "'": "&#39;" }[c]));
}

async function sendCmd(command, name) {
  if (state.sending) {
    setStatus("上一条命令还在发送，请稍候…", true);
    return;
  }
  state.sending = true;
  document.body.classList.add("busy");
  setStatus(`发送中：${name}`);
  try {
    const data = await api("/api/send", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        command,
        console_key: $("consoleKey").value,
        skip_open: $("consoleReady").value === "ready",
      }),
    });
    setStatus(data.message || `已发送 ${command}`);
  } catch (err) {
    setStatus(err.message, true);
  } finally {
    state.sending = false;
    document.body.classList.remove("busy");
  }
}

async function saveSettings() {
  await api("/api/settings", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      game_root: $("gameRoot").value.trim(),
      workshop: $("workshop").value.trim(),
      extra_mods: $("extraMods").value.trim(),
      console_key: $("consoleKey").value,
      version: state.version,
    }),
  });
}

$("cats").addEventListener("click", (ev) => {
  const id = ev.target.dataset && ev.target.dataset.id;
  if (!id) return;
  state.category = id;
  state.page = 1;
  loadList().catch((e) => setStatus(e.message, true));
});
$("version").addEventListener("change", () => {
  state.version = $("version").value;
  state.page = 1;
  scheduleSave();
  loadList().catch((e) => setStatus(e.message, true));
});
let searchTimer = 0;
$("q").addEventListener("input", () => {
  clearTimeout(searchTimer);
  searchTimer = setTimeout(() => {
    state.page = 1;
    loadList().catch((e) => setStatus(e.message, true));
  }, 220);
});

let saveTimer = 0;
function scheduleSave() {
  clearTimeout(saveTimer);
  saveTimer = setTimeout(() => {
    saveSettings().catch((e) => setStatus("设置保存失败: " + e.message, true));
  }, 400);
}
["gameRoot", "workshop", "extraMods"].forEach((id) => {
  $(id).addEventListener("input", scheduleSave);
});
$("consoleKey").addEventListener("change", scheduleSave);
$("prev").addEventListener("click", () => {
  if (state.page > 1) {
    state.page -= 1;
    loadList().catch((e) => setStatus(e.message, true));
  }
});
$("next").addEventListener("click", () => {
  if (state.page < state.totalPages) {
    state.page += 1;
    loadList().catch((e) => setStatus(e.message, true));
  }
});
$("btnOpenConsole").addEventListener("click", async () => {
  setStatus("正在打开群星控制台…");
  try {
    const data = await api("/api/open_console", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ console_key: $("consoleKey").value }),
    });
    setStatus(data.message || "已打开控制台");
    $("consoleReady").value = "ready";
  } catch (e) {
    setStatus(e.message, true);
  }
});
$("btnDetect").addEventListener("click", async () => {
  try {
    const data = await api("/api/detect");
    setStatus(data.message || "已找到群星");
  } catch (e) {
    setStatus(e.message, true);
  }
});
$("btnScan").addEventListener("click", async () => {
  setStatus("正在扫描游戏与 MOD…");
  try {
    await saveSettings();
    const data = await api("/api/scan", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: "{}",
    });
    setStatus(data.message || "扫描完成");
    state.page = 1;
    await loadList();
  } catch (e) {
    setStatus(e.message, true);
  }
});
$("btnAdminHint").addEventListener("click", () => {
  setStatus("请以管理员运行。先启动群星（非铁人），点「打开控制台」或手动按 ~，确认控制台出现后再点命令按钮。");
});

let resizeTimer = 0;
window.addEventListener("resize", () => {
  clearTimeout(resizeTimer);
  resizeTimer = setTimeout(() => {
    const before = state.pageSize;
    const { changed } = fitGrid();
    if (changed || before !== state.pageSize) {
      state.page = 1;
      loadList().catch((e) => setStatus(e.message, true));
    }
  }, 120);
});

loadMeta().then(loadList).catch((e) => setStatus(e.message, true));
