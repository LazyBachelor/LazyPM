/**
 * Drag and drop for board view - allows moving issue cards between status columns
 */
(function () {
  document.addEventListener("dragstart", function (e) {
    if (e.target.closest("button") || e.target.closest("a")) return;
    const card = e.target.closest(".board-card");
    if (!card) return;
    e.dataTransfer.setData("text/plain", card.dataset.issueId);
    e.dataTransfer.effectAllowed = "move";
    card.classList.add("opacity-50");
  });

  document.addEventListener("dragend", function (e) {
    const card = e.target.closest(".board-card");
    if (card) card.classList.remove("opacity-50");
  });

  document.addEventListener("dragover", function (e) {
    const zone = e.target.closest(".column-drop-zone");
    if (!zone) return;
    e.preventDefault();
    e.dataTransfer.dropEffect = "move";
    zone.classList.add("ring-2", "ring-primary", "ring-inset");
  });

  document.addEventListener("dragleave", function (e) {
    const zone = e.target.closest(".column-drop-zone");
    if (!zone || zone.contains(e.relatedTarget)) return;
    zone.classList.remove("ring-2", "ring-primary", "ring-inset");
  });

  document.addEventListener("drop", function (e) {
    const zone = e.target.closest(".column-drop-zone");
    if (!zone) return;
    e.preventDefault();
    zone.classList.remove("ring-2", "ring-primary", "ring-inset");
    const issueId = e.dataTransfer.getData("text/plain");
    const newStatus = zone.dataset.status;
    if (!issueId || !newStatus) return;

    const formData = new URLSearchParams();
    formData.append("status", newStatus);

    fetch("/issues/" + issueId + "?from=board", {
      method: "PATCH",
      body: formData,
      headers: {
        "Content-Type": "application/x-www-form-urlencoded",
        "HX-Request": "true",
      },
    }).then(function (r) {
      if (r.ok) {
        var redirect = r.headers.get("HX-Redirect");
        if (redirect) {
          window.location.href = redirect;
        } else if (typeof htmx !== "undefined") {
          htmx.ajax("GET", "/?board=true", {
            target: "main",
            swap: "innerHTML",
          });
        } else {
          window.location.href = "/?board=true";
        }
      }
    });
  });
})();
