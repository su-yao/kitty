(function () {
  var t = "light";
  try { t = localStorage.getItem("stellaris-panel-theme") || "light"; } catch (e) {}
  if (t !== "dark") t = "light";
  document.documentElement.setAttribute("data-theme", t);
})();
