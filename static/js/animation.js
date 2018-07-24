function ShowTree(data) {
  $("#currentPath").html(data.bucketsPath);
  if (data.recordsAmount == 0) {
    $("#recordsAmount").text("(empty)");
  } else if (data.recordsAmount == 1) {
    $("#recordsAmount").text("(" + data.recordsAmount + " item)");
  } else {
    $("#recordsAmount").text("(" + data.recordsAmount + " items)");
  }

  $("#dbTree").empty();

  if (data.prevRecords) {
    $("#dbTree").append(getPrevRecordsButton());
  } else if (data.prevBucket) {
    $("#dbTree").append(getBackButton());
  }

  // Update currentData
  var records = data.records;
  currentData = {};
  for (i in records) {
    currentData[records[i].key] = records[i].value;
  }

  for (i in records) {
    if (records[i].type == "bucket") {
      $("#dbTree").append(getBucketButton(records[i].key));
    } else if (records[i].type == "record") {
      $("#dbTree").append(getRecordButton(records[i].key, records[i].value));
    }
  }

  if (data.nextRecords) {
    $("#dbTree").append(getNextRecordsButton());
  }

  $("#dbTreeWrapper").scrollTop(0);
}

function ShowDBsList() {
  $("#dbListBackground").css("display", "block");
  $("#dbList").addClass("db_list_animation");
}

function ShowFullRecord(key) {
  // This 2 lines fix the bug, when #recordData disappeared after changing of #recordPath and #currentPath
  $("#recordPath").html("");
  $("#recordValue").html("");

  $("#recordData").scrollTop(0);

  $("#recordPath").html(key);
  $("#recordValue").html(currentData[key]);
}

// For offering of already opened before dbs
function ShowOpenedDBsPaths() {
  var sortedPaths = GetPaths();

  $("#paths").empty();
  for (var i = 0; i < sortedPaths.length; i++) {
    $("#paths").append($("<option>", { value: sortedPaths[i] }));
  }

  $("#openDbWindow").css("display", "block");
  $("#DBPath").focus();
}

function HideOpenDbWindow() {
  $("#openDbWindow").css("display", "none");
  $("#dbPathsList").css("display", "none");
}

function SwitchPathsForDeleting() {
  if ($("#dbPathsList").css("display") == "none") {
    ShowPathsForDeleting();
  } else {
    $("#dbPathsList").css("display", "none");
  }
}

function ShowPathsForDeleting() {
  var paths = GetPaths();
  $("#dbPathsList").empty();

  if (paths.length == 0) {
    $("#dbPathsList").append($("<span>").text("Empty"));
  } else {
    for (var i = 0; i < paths.length; i++) {
      $("#dbPathsList").append(getPathForDeleting(paths[i]));
    }
  }

  $("#dbPathsList").css("display", "block");
}

/* Popups */
function ShowErrorPopup(message) {
  $("#popupMessage").html(message);
  $("#errorPopup").addClass("popup_animation");
}

function HideErrorPopup() {
  $("#errorPopup").removeClass("popup_animation");
}

function ShowDonePopup() {
  $("#donePopup").css("display", "block");
  $("#donePopup").addClass("done_popup_animation");

  setTimeout(function() {
    $("#donePopup").css("display", "none");
    $("#donePopup").removeClass("done_popup_animation");
  }, 2000);
}

/* Menus */
function showPopupMenu(x, y, htmlElem) {
  $("#popupMenu").css("top", y + 10 + "px");
  $("#popupMenu").css("left", x + 10 + "px");
  $("#popupMenu").empty();
  $("#popupMenu").append(htmlElem);
  $("#popupMenu").css("visibility", "visible");
}

function showBucketMenu(event) {
  showPopupMenu(
    event.clientX,
    event.clientY,
    getBucketMenu(event.target.innerText)
  );
  return false;
}

function showRecordMenu(event) {
  showPopupMenu(
    event.clientX,
    event.clientY,
    getRecordMenu(event.target.innerText)
  );
  return false;
}

function showAddMenu(event) {
  // Magic. Program shows addMenu only if target is dbTreeWrapper and isn't anything else
  if (
    (event.target == dbTreeWrapper || event.target == dbTree) &&
    event.which == 3
  ) {
    // Show menu only if db was chosen
    if (currentDB.dbPath != "") {
      showPopupMenu(event.clientX, event.clientY, getAddMenu());
      return false;
    }
  }
}

/* Windows */
// Write mode only
function ShowAddModal(type) {
  $("#addItemWindow").empty();
  if (type == "bucket") {
    $("#addItemWindow").append(getAddBucketWindow());
    $("#addItemWindowBackground").css("display", "block");
  } else if (type == "record") {
    $("#addItemWindow").append(getAddRecordWindow());
    $("#addItemWindowBackground").css("display", "block");
  }
}

function ShowEditModal(type, target) {
  // target == key of the record or bucket
  $("#addItemWindow").empty();
  if (type == "bucket") {
    $("#addItemWindow").append(getEditBucketWindow(target));
    $("#addItemWindowBackground").css("display", "block");
  } else if (type == "record") {
    $("#addItemWindow").append(getEditRecordWindow(target));
    $("#addItemWindowBackground").css("display", "block");
  }
}

function HideAddModal() {
  $("#addItemWindowBackground").css("display", "none");
}
