// MujicaUI docs: theme toggle, copy buttons, mobile nav, index search. No network requests.
(function () {
  var root = document.documentElement;

  function effectiveTheme() {
    if (root.dataset.theme) return root.dataset.theme;
    return window.matchMedia("(prefers-color-scheme: dark)").matches ? "dark" : "light";
  }
  var themeBtn = document.querySelector(".theme");
  function labelTheme() {
    if (themeBtn) themeBtn.textContent = effectiveTheme() === "dark" ? "浅色" : "深色";
  }
  if (themeBtn) {
    labelTheme();
    themeBtn.addEventListener("click", function () {
      var next = effectiveTheme() === "dark" ? "light" : "dark";
      root.dataset.theme = next;
      try { localStorage.setItem("mujica-theme", next); } catch (e) {}
      labelTheme();
    });
    window.matchMedia("(prefers-color-scheme: dark)").addEventListener("change", labelTheme);
  }

  document.querySelectorAll(".code .copy").forEach(function (btn) {
    btn.addEventListener("click", function () {
      var code = btn.parentElement.querySelector("code").textContent;
      function done(ok) {
        btn.textContent = ok ? "已复制" : "复制失败";
        setTimeout(function () { btn.textContent = "复制"; }, 1500);
      }
      if (navigator.clipboard && navigator.clipboard.writeText) {
        navigator.clipboard.writeText(code).then(function () { done(true); }, function () { done(false); });
      } else {
        done(false);
      }
    });
  });

  var toggle = document.querySelector(".navtoggle");
  var side = document.getElementById("sidenav");
  if (toggle && side) {
    toggle.addEventListener("click", function () {
      var open = side.classList.toggle("open");
      toggle.setAttribute("aria-expanded", open ? "true" : "false");
    });
    var cur = side.querySelector('[aria-current="page"]');
    if (cur && window.matchMedia("(min-width: 821px)").matches) cur.scrollIntoView({ block: "center" });
  }

  var dataEl = document.getElementById("search-data");
  var q = document.getElementById("q");
  if (!dataEl || !q) return;
  var items = JSON.parse(dataEl.textContent);
  var list = document.getElementById("results");
  var count = document.getElementById("qcount");
  function render() {
    var term = q.value.trim().toLowerCase();
    list.textContent = "";
    if (!term) { count.textContent = ""; return; }
    var hits = items.filter(function (it) {
      return it.n.toLowerCase().indexOf(term) >= 0 || it.z.indexOf(term) >= 0 ||
        it.id.indexOf(term) >= 0 || it.s.indexOf(term) >= 0 || it.c.toLowerCase().indexOf(term) >= 0;
    });
    count.textContent = hits.length + " 个结果";
    hits.slice(0, 60).forEach(function (it) {
      var li = document.createElement("li");
      var a = document.createElement("a");
      a.href = "/c/" + it.s + "/";
      [["num", it.id], ["n", it.n], ["zh", it.z], ["c", it.c]].forEach(function (p) {
        var span = document.createElement("span");
        span.className = p[0];
        span.textContent = p[1];
        a.appendChild(span);
      });
      li.appendChild(a);
      list.appendChild(li);
    });
  }
  q.addEventListener("input", render);
  render();
})();
