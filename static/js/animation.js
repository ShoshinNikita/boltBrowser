function ShowTree(data) {
	$("#currentPath").html(data.bucketsPath);
	if (data.recordsAmount == 0) {
		$("#recordsAmount").text("(empty)")
	} else if (data.recordsAmount == 1) {
		$("#recordsAmount").text("(" + data.recordsAmount + " item)")
	} else {
		$("#recordsAmount").text("(" + data.recordsAmount + " items)")
	}

	$("#dbTree").empty();

	if (data.prevRecords) {
		result += prevRecordsButtonTemplate;
	} else if (data.prevBucket) {
		$("#dbTree").append(getBackButton());
	}

	// Update currentData
	var records = data.records;
	currentData = {}
	for (i in records) {
		currentData[records[i].key] = records[i].value
	}

	for (i in records) {
		if (records[i].type == "bucket") {
			$("#dbTree").append(getBucket(records[i].key));
		} else if (records[i].type == "record") {
			$("#dbTree").append(getRecord(records[i].key, records[i].value));
		}
	}

	if (data.nextRecords) {
		$("#dbTree").append(getNextRecordsButton());
	}

	document.getElementById("dbTreeWrapper").scrollTop = 0;
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

function SwitchPathsForDelete() {
	if ($("#dbPathsList").css("display") == "none") {
		ShowPathsForDelete();
	} else {
		$("#dbPathsList").css("display", "none");
	}
}

function ShowPathsForDelete() {
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

// ErrorPopup
function ShowErrorPopup(message) {
	$("#popupMessage").html(message);
	$("#errorPopup").addClass("popup_animation");
}

function HideErrorPopup() {
	$("#errorPopup").removeClass("popup_animation");
}

// OpenDbWindow
function ShowOpenDbWindow() {
	const template = `<option value="{0}">`;

	var sortedPaths = GetPaths();

	var options = "";
	for (var i = 0; i < sortedPaths.length && i < 5; i++) {
		options += template.format(sortedPaths[i]);
	}

	$("#paths").html(options);
	$("#openDbWindow").css("display", "block");
	$("#DBPath").focus();
}

function HideOpenDbWindow() {
	$("#openDbWindow").css("display", "none");
	$("#dbPathsList").css("display", "none");
}

// DonePopup
function ShowDonePopup() {
	$("#donePopup").css("display", "block")
	$("#donePopup").addClass("done_popup_animation")

	setTimeout(function() {
		$("#donePopup").css("display", "none")
		$("#donePopup").removeClass("done_popup_animation")
	}, 2000)
}