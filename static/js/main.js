function getDbButton(dbPath, dbName) {
	var $input = $("<input>", {type: "button", class:"db_button", title: "Choose", value: dbName}).
		click({dbPath: dbPath}, function(event) {
			ChooseDB(event.data.dbPath);
		});

	var attr = {class: "material-icons btn",
				style: "float: right; margin-right: 10px; font-size: 30px !important;",
				title: "Close"};
	var $closeBtn = $("<i>", attr).text("close").
		click({dbPath: dbPath}, function(event){
			CloseDB(event.data.dbPath);
		});

	return $("<div>").append($input).append($closeBtn);
}

function getRecord(key, value) {
	var $icon = $("<i>", {class: "material-icons"}).text("assignment");
	var $key = $("<span>", {class: "record", id: "key", style: "font-weight: bold;"}).html(key).
		click({key: key}, function(event) {
			ShowFullRecord(event.data.key);
		});
	var $value = $("<span>", {id: "value"}).html(" – " + value);

	return $("<div>", {style: "display: table;"}).append($icon).append($key).append($value);
}

function getBucket(key) {
	var $icon = $("<i>", {class: "material-icons"}).text("folder");
	var $key = $("<span>", {class: "bucket", style: "font-weight: bold;"}).html(key).
		click({key: key}, function(event) {
			Next(event.data.key);
		});
	return $("<div>", {style: "display: table;"}).append($icon).append($key);
}

function getBackButton() {
	var $icon = $("<i>", {class: "material-icons btn", title: "Back"}).text("more_horiz").
		click(function(){
			Back();
		});

	return $("<div>", {style: "display: table;"}).append($icon);
}

function getNextRecordsButton() {
	var $icon = $("<i>", {class: "material-icons"}).text("arrow_forward_ios");
	var $btn = $("<span>", {style: "cursor: pointer; font-weight: bold;"}).text("Next page");
		click(function() {
			NextRecords();
		});

	return $("<div>", {style: "display: table;"}).append($icon).append($btn);
}

function getPrevRecordsButton() {
	var $icon = $("<i>", {class: "material-icons"}).text("arrow_back_ios");
	var $btn = $("<span>", {style: "cursor: pointer; font-weight: bold;"}).text("Previous page");
		click(function() {
			NextRecords();
		});

	return $("<div>", {style: "display: table;"}).append($icon).append($btn);
}

function getPathForDeleting(path) {
	var $path = $("<span>").text(path);
	var $btn = $("<i>", {class: "material-icons btn", style: "float: right; font-size: 22px !important; vertical-align: middle;", title: "Delete"}).text("close").
		click({path: path}, function(event) {
			DeletePath(event.data.path);
		});

	return $("<div>", {style: "margin-bottom: 10px; text-align: left;"}).append($path).append($btn);
}


/* Global variables */
var currentDBPath = "";
// Dictionary. It keeps data like "key of a record": "value of a record"
var currentData = {};


/* Local Storage */
function PrepareLS() {
	if (localStorage.getItem("paths") === null) {
		var paths = {};
		localStorage.setItem("paths", JSON.stringify(paths));
	}
}

function putIntoLS(dbPath) {
	var paths = JSON.parse(localStorage.getItem("paths"));
	if (paths[dbPath] == null) {
		paths[dbPath] = {"uses": 1};
	} else {
		paths[dbPath].uses += 1;
	}

	localStorage.setItem("paths", JSON.stringify(paths));
}

function getPaths() {
	var paths = JSON.parse(localStorage.getItem("paths"));

	// Sorting. Return only keys;
	var sortedPaths = Object.keys(paths).sort(function(a, b){
		if (paths[a].uses < paths[b].uses) {
			return 1;
		}
		if (paths[a].uses > paths[b].uses) {
			return -1;
		}
		return 0;
	});

	return sortedPaths;
}

function DeletePath(path) {
	var paths = JSON.parse(localStorage.getItem("paths"));
	delete paths[path];
	localStorage.setItem("paths", JSON.stringify(paths));

	showPathsForDelete();
}


/* API */
function OpenDB() {
	var dbPath = $("#DBPath").val();
	if (dbPath == "" ) {
		ShowErrorPopup("Error: path is empty");
		return;
	}

	$("#DBPath").val("");
	$.ajax({
		url: "/api/databases",
		type: "POST",
		data: {
			"dbPath": dbPath
		},
		success: function(result){
			result= JSON.parse(result);
			putIntoLS(result.dbPath);
			HideOpenDbWindow();
			ShowDBList();
		},
		error: function(result) {
			ShowErrorPopup(result.responseText);
		}
	});
}

function CreateDB() {
	var path = $("#DBPathForCreating").val();
	if (path == "" ) {
		ShowErrorPopup("Error: path is empty");
		return;
	}

	$("#DBPathForCreating").val("");
	$.ajax({
		url: "/api/databases/new",
		type: "POST",
		data: {
			"path": path
		},
		success: function(result){
			result= JSON.parse(result);
			putIntoLS(result.dbPath);
			HideOpenDbWindow();
			ShowDBList();
		},
		error: function(result) {
			ShowErrorPopup(result.responseText);
		}
	});
}

function CloseDB(dbPath) {
	$.ajax({
		url: "/api/databases" + "?" + $.param({"dbPath": dbPath}),
		type: "DELETE",
		success: function(result){
			if (dbPath == currentDBPath) {
				$("#dbName").html("<i>Name:<\/i> ?");
				$("#dbPath").html("<i>Path:<\/i> ?");
				$("#dbSize").html("<i>Size:<\/i> ?");
				$("#dbTree").html("");
				$("#currentPath").html("");
				$("#recordsAmount").html("");
				$("#recordPath").html("?");
				$("#recordValue").html("?");
				$("#searchBox").css("visibility", "hidden");
				currentDBPath = "";
			}
			ShowDBList();
		},
		error: function(result) {
			ShowErrorPopup(result.responseText);
		}
	});
}

function ShowDBList() {
	$.ajax({
		url: "/api/databases",
		type: "GET",
		success: function(result){
			allDB = JSON.parse(result);

			$("#list").empty();
			for (i in allDB) {
				$("#list").append(getDbButton(allDB[i].dbPath, allDB[i].name));
			}
		},
		error: function(result) {
			ShowErrorPopup(result.responseText);
		}
	});
}

function ChooseDB(dbPath) {
	currentDBPath = dbPath;
	$.ajax({
		url: "/api/buckets/current",
		type: "GET",
		data: {
			"dbPath": dbPath,
		},
		success: function(result){
			result = JSON.parse(result);

			$("#dbName").html("<i>Name:<\/i> " + result.db.name);
			$("#dbPath").html("<i>Path:<\/i> " + result.db.dbPath);
			$("#dbSize").html("<i>Size:<\/i> " + result.db.size / 1024 + " Kb");
			$("#recordPath").html("?");
			$("#recordValue").html("?");
			$("#searchBox").css("visibility", "visible");

			ShowTree(result);
		},
		error: function(result) {
			ShowErrorPopup(result.responseText);
		}
	});
}

function Next(bucket) {
	$.ajax({
		url: "/api/buckets/next",
		type: "GET",
		data: {
			"dbPath": currentDBPath,
			"bucket": bucket
		},
		success: function(result){
			result = JSON.parse(result);
			ShowTree(result);
		},
		error: function(result) {
			ShowErrorPopup(result.responseText);
		}
	});
}

function Back() {
	$.ajax({
		url: "/api/buckets/back",
		type: "GET",
		data: {
			"dbPath": currentDBPath,
		},
		success: function(result){
			result = JSON.parse(result);
			ShowTree(result);
		},
		error: function(result) {
			ShowErrorPopup(result.responseText);
		}
	});
}

function NextRecords() {
	$.ajax({
		url: "/api/records/next",
		type: "GET",
		data: {
			"dbPath": currentDBPath,
		},
		success: function(result){
			result = JSON.parse(result);

			ShowTree(result);
		},
		error: function(result) {
			ShowErrorPopup(result.responseText);
		}
	});
}

function PrevRecords() {
	$.ajax({
		url: "/api/records/prev",
		type: "GET",
		data: {
			"dbPath": currentDBPath,
		},
		success: function(result){
			result = JSON.parse(result);

			ShowTree(result);
		},
		error: function(result) {
			ShowErrorPopup(result.responseText);
		}
	});
}

function Search() {
	var text = $("#searchText").val();
	if (text == "") {
		ChooseDB(currentDBPath);
		return;
	}

	var mode = "plain";
	if ($("#regexMode").prop("checked")) {
		mode = "regex";
	}
	$.ajax({
		url: "/api/search",
		type: "GET",
		data: {
			"dbPath": currentDBPath,
			"text": text,
			"mode": mode
		},
		success: function(result){
			result = JSON.parse(result);

			ShowTree(result);
		},
		error: function(result) {
			ShowErrorPopup(result.responseText);
		}
	});
}


/* Animation */
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

function SwitchPathsForDelete() {
	if ($("#dbPathsList").css("display") == "none") {
		showPathsForDelete();
	} else {
		$("#dbPathsList").css("display", "none");
	}
}

function showPathsForDelete() {
	var paths = getPaths();
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

	var sortedPaths = getPaths();

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


/* Secondary functions */
window.onclick = function(event) {
    if (event.target == openDbWindow) {
		HideOpenDbWindow();
	}
	if (event.target == dbListBackground) {
		$("#dbListBackground").css("display", "none");
		$("#dbList").removeClass("db_list_animation");
	}

	// From write_mode.js
	// Hiding AddMenu
	if (event.target == addItemWindowBackground) {
		$("#addItemWindowBackground").css("display", "none");
	}
	// Hiding PopupMenu
	if ($("#popupMenu").css("visibility") == "visible" && event.target != popupMenu) {
		$("#popupMenu").css("visibility", "hidden");
	}
}

window.onkeydown = function(event) {
	if (event.target == searchText) {
		// Enter
		if (event.keyCode == 13 || event.which == 13) {
			Search();
		}
		// Esc
		if (event.keyCode == 27 || event.which == 27) {
			ChooseDB(currentDBPath);
		}
	}
}

String.prototype.format = function () {
	var a = this;
	for (var k in arguments) {
		a = a.replace(new RegExp("\\{" + k + "\\}", 'g'), arguments[k]);
	}
	return a;
}

String.prototype.replaceAll = function(search, replacement) {
    var target = this;
    return target.replace(new RegExp(search, 'g'), replacement);
};