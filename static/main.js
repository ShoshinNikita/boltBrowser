/* Constants */
const buttonTemplate = `<div>
	<input type="button" class="db_button" value="{0}" onclick="ChooseDB('{1}')" title="Choose">
	<i class="material-icons btn" style="float: right; margin-right: 1vw; font-size: 30px !important;" title="Close" onclick="CloseDB('{1}');">close<\/i>
<\/div>`;

const recordTemplate = `<div style="display: table;">
	<i class="material-icons" icon>assignment<\/i>
	<span class="record" onclick="ShowFullRecord({0});"><b>{1}<\/b><\/span>: {2}
<\/div>`;

const bucketTemplate = `<div style="display: table;">
	<i class="material-icons" icon>folder<\/i>
	<span class="bucket" onclick="Next('{0}');"><b>{0}<\/b><\/span>
<\/div>`;

const backButton = `<div style="display: table;">
	<i class="material-icons btn" icon onclick="Back();" title="Back">more_horiz<\/i>
<\/div>`;

const nextRecordsButtonTemplate = `<div style="display: table;">
	<i class="material-icons" icon>arrow_forward_ios<\/i>
	<span style="cursor: pointer;" onclick="NextRecords();"><b>Next page<\/b><\/span>
<\/div>`;

const prevRecordsButtonTemplate = `<div style="display: table;">
	<i class="material-icons" icon>arrow_back_ios<\/i>
	<span style="cursor: pointer;" onclick="PrevRecords();"><b>Previous page<\/b><\/span>
<\/div>`;


/* Global variables */
var currentDBPath = "";
var currentData = null;


/* Local Storage */
function PrepareLS() {
	if (localStorage.getItem("paths") === null) {
		var paths = {}
		localStorage.setItem("paths", JSON.stringify(paths));
	}
}

function putIntoLS(dbPath) {
	var paths = SafeParse(localStorage.getItem("paths"));
	if (paths[dbPath] == null) {
		paths[dbPath] = {
			"uses": 1
		}
	} else {
		paths[dbPath].uses += 1;
	}

	localStorage.setItem("paths", JSON.stringify(paths));
}

function getPaths() {
	var paths = SafeParse(localStorage.getItem("paths"));

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
	var paths = SafeParse(localStorage.getItem("paths"));
	delete paths[path];
	localStorage.setItem("paths", JSON.stringify(paths));

	ShowPathsForDelete();
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
			result= SafeParse(result)
			putIntoLS(result.dbPath);
			ShowDBList();
		},
		error: function(result) {
			ShowErrorPopup(result.responseText);
		}
	});
	;
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
			allDB = SafeParse(result);
			var result = "";
			for (i in allDB) {
				result += buttonTemplate.format(allDB[i].name, allDB[i].dbPath);
			}
			$("#list").html(result);
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
			result = SafeParse(result);

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
			result = SafeParse(result);
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
			result = SafeParse(result);
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
			result = SafeParse(result);

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
			result = SafeParse(result);

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
			result = SafeParse(result);
			
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

function ShowFullRecord(number) {
	var currentPath = $("#currentPath").text();
	$("#recordPath").html(currentData[number].key + " – <i>" + currentPath + "<\/i>");
	$("#recordValue").html(currentData[number].value);
}

function ShowTree(data) {
	$("#currentPath").html(data.bucketsPath);
	if (data.recordsAmount == 0) {
		$("#recordsAmount").html("(empty)")
	} else if (data.recordsAmount == 1) {
		$("#recordsAmount").html("(" + data.recordsAmount + " item)")
	} else {
		$("#recordsAmount").html("(" + data.recordsAmount + " items)")
	}
	
	var result = "";
	if (data.prevRecords) {
		result += prevRecordsButtonTemplate;
	} else if (data.prevBucket) {
		result = backButton;
	}

	var records = data.records;
	currentData = records;

	for (i in records) {
		if (records[i].type == "bucket") {
			result += bucketTemplate.format(records[i].key);
		} else if (records[i].type == "record") {
			var value = records[i].value;
			if (value.length > 40) {
				value = value.substring(0, 60);
				value += "...";
			}
			result += recordTemplate.format(i, records[i].key, value);
		}
	}

	if (data.nextRecords) {
		result += nextRecordsButtonTemplate;
	}
	$("#dbTree").html(result);

	document.getElementById("dbTreeWrapper").scrollTop = 0;
}

function ShowPathsForDelete() {
	const button = `<div style="margin-bottom: 10px; text-align: left;"><span>{0}</span>
	<i class="material-icons btn" style="float: right; margin-right: 1vw; font-size: 25px !important;" title="Delete" onclick="DeletePath('{0}');">close<\/i><\/div>`;

	var paths = getPaths();

	var res = ""
	for (var i = 0; i < paths.length; i++) {
		res += button.format(paths[i]);
	}

	if (res == "") {
		res = "Empty"
	}

	$("#dbPathsList").html(res);
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

// Return parsed object with "good" symbols
function SafeParse(text) {
	var object = JSON.parse(text);
	makeSafe(object);
	return object
}

// Erase all "bad" symbols from object like '<', '>', '\'', '"'
// Works recursively
function makeSafe(object) {
	if (typeof object === 'object') {
		for (i in object) {
			object[i] = makeSafe(object[i]);
		}
	} else if (typeof object === 'string') {
		object = object.replaceAll("<", "❮").replaceAll(">", "❯").replaceAll("\"", "＂").replaceAll("'", "ߴ")
	}
	return object
}

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