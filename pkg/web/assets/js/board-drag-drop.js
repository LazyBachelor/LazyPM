/**
 * Drag and drop for board view - allows moving issue cards between status columns
 */
(function () {
  document.addEventListener("dragstart", function (e) {
    if (e.target.closest("button") || e.target.closest("a")) return;
    const card = e.target.closest(".board-card");
    if (!card) return;
    
    const column = card.closest(".board-column");
    const sourceStatus = column?.dataset?.status || "";
    const sourceSprint = column?.dataset?.sprint || "";
    
    e.dataTransfer.setData("text/plain", card.dataset.issueId);
    e.dataTransfer.setData("source-status", sourceStatus);
    e.dataTransfer.setData("source-sprint", sourceSprint);
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
    const sourceStatus = e.dataTransfer.getData("source-status");
    const sourceSprint = e.dataTransfer.getData("source-sprint");
    const newStatus = zone.dataset.status;
    const targetSprint = zone.dataset.sprint;
    
    if (!issueId || !newStatus) return;

    let sourceSprintNum = 0;
    if (sourceSprint && sourceSprint.startsWith("Sprint ")) {
      sourceSprintNum = parseInt(sourceSprint.replace("Sprint ", ""), 10);
    }
    
    let targetSprintNum = 0;
    if (targetSprint && targetSprint.startsWith("Sprint ")) {
      targetSprintNum = parseInt(targetSprint.replace("Sprint ", ""), 10);
    }

    const formData = new URLSearchParams();
    
    let issueStatus = newStatus;
    if (newStatus === "todo" || newStatus === "backlog") {
      issueStatus = "open";
    } else if (newStatus === "done") {
      issueStatus = "closed";
    }
    
    formData.append("status", issueStatus);

    if (sourceStatus === "backlog" && newStatus !== "backlog" && targetSprintNum > 0) {
      formData.append("add_to_sprint", targetSprintNum.toString());
    }
    else if (sourceStatus !== "backlog" && newStatus === "backlog" && sourceSprintNum > 0) {
      formData.append("remove_from_sprint", sourceSprintNum.toString());
    }
    else if (sourceStatus !== "backlog" && newStatus !== "backlog" && targetSprintNum > 0) {
      formData.append("add_to_sprint", targetSprintNum.toString());
    }

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
