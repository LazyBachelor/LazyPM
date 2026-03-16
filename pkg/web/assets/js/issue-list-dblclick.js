/**
 * Open edit-issue modal when double-clicking a row in the issue list (list view).
 * Uses event delegation and preventDefault so double-click doesn't select text.
 */
(function () {
  document.addEventListener("dblclick", function (e) {
    var row = e.target.closest("tr[data-issue-id]");
    if (!row) return;
    e.preventDefault();
    var id = row.getAttribute("data-issue-id");
    if (!id) return;
    if (typeof htmx !== "undefined") {
      htmx.ajax("GET", "/issues/" + id + "/edit", {
        target: "#modal-container",
        swap: "innerHTML",
      });
    }
  });
})();
