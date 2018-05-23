/* Constants */
const buttonTemplate = `<div>
	<input type="button" class="db_button" value="{0}" onclick="ChooseDB('{1}')" title="Choose">
	<i class="material-icons btn" style="float: right; margin-right: 1vw; font-size: 30px !important;" title="Close" onclick="CloseDB('{1}');">close<\/i>
<\/div>`;

const recordTemplate = `<div>
	<i class="material-icons" icon>assignment<\/i>
	<span class="record" onclick="ShowFullRecord({0});"><b>{1}<\/b>:<\/span> {2}
<\/div>`;

const bucketTemplate = `<div>
	<i class="material-icons" icon>folder<\/i>
	<span class="bucket" onclick="Next('{0}');"><b>{0}<\/b><\/span>
<\/div>`;

const backButton = `<div>
	<i class="material-icons btn" icon onclick="Back();" title="Back">more_horiz<\/i>
<\/div>`;

const nextRecordersButtonTemplate = `<div>
	<i class="material-icons" icon>arrow_forward_ios<\/i>
	<span style="cursor: pointer;" onclick="NextRecords();"><b>Next page<\/b><\/span>
<\/div>`;

const prevRecordersButtonTemplate = `<div>
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
	var paths = JSON.parse(localStorage.getItem("paths"));
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

	ShowPathsForDelete();
}


/* API */
function OpenDB() {
	var dbPath = $("#DBPath").val();
	if (dbPath == "" ) {
		ShowPopup("Error: path is empty");
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
			result= JSON.parse(result)
			putIntoLS(result.dbPath);
			ShowDBList();
		},
		error: function(result) {
			ShowPopup(result.responseText);
		}
	});
	;
}

function CloseDB(dbPath) {
	$.ajax({
		url: "/api/closeDB",
		type: "POST",
		data: {
			"dbPath": dbPath,
		},
		success: function(result){
			if (dbPath == currentDBPath) {
				$("#dbName").html("<i>Name:<\/i> ?");
				$("#dbPath").html("<i>Path:<\/i> ?");
				$("#dbSize").html("<i>Size:<\/i> ?");
				$("#dbTree").html("");
				$("#currentPath").html("");
				$("#recordPath").html("?");
				$("#recordValue").html("?");
				$("#searchBox").css("visibility", "hidden");
				currentDBPath = "";
			}
			ShowDBList();
		},
		error: function(result) {
			ShowPopup(result.responseText);
		}
	});
}

function ShowDBList() {
	$.ajax({
		url: "/api/databases",
		type: "GET",
		success: function(result){
			allDB = JSON.parse(result);
			var result = "";
			for (i in allDB) {
				result += buttonTemplate.format(allDB[i].name, allDB[i].dbPath);
			}
			$("#list").html(result);
		},
		error: function(result) {
			ShowPopup(result.responseText);
		}
	});
}

function ChooseDB(dbPath) {
	currentDBPath = dbPath;
	$.ajax({
		url: "/api/current",
		type: "GET",
		data: {
			"dbPath": dbPath,
		},
		success: function(result){
			result = JSON.parse(result);

			$("#dbName").html("<i>Name:<\/i> " + result.db.name);
			$("#dbPath").html("<i>Path:<\/i> " + result.db.dbPath);
			$("#dbSize").html("<i>Size:<\/i> " + result.db.size / 1024 + " Kb");
			$("#currentPath").html("<i>" + result.bucketsPath + "<\/i> ");
			$("#recordPath").html("?");
			$("#recordValue").html("?");
			$("#searchBox").css("visibility", "visible");
		;
			ShowTree(result);
		},
		error: function(result) {
			ShowPopup(result.responseText);
		}
	});
}

function Next(bucket) {
	$.ajax({
		url: "/api/next",
		type: "GET",
		data: {
			"dbPath": currentDBPath,
			"bucket": bucket
		},
		success: function(result){
			result = JSON.parse(result);
			$("#currentPath").html("<i>" + result.bucketsPath + "<\/i> ");
			ShowTree(result);
		},
		error: function(result) {
			ShowPopup(result.responseText);
		}
	});
}

function Back() {
	$.ajax({
		url: "/api/back",
		type: "GET",
		data: {
			"dbPath": currentDBPath,
		},
		success: function(result){
			result = JSON.parse(result);
			$("#currentPath").html("<i>" + result.bucketsPath + "<\/i> ");
			ShowTree(result);
		},
		error: function(result) {
			ShowPopup(result.responseText);
		}
	});
}

function NextRecords() {
	$.ajax({
		url: "/api/nextRecords",
		type: "GET",
		data: {
			"dbPath": currentDBPath,
		},
		success: function(result){
			result = JSON.parse(result);

			ShowTree(result);
		},
		error: function(result) {
			ShowPopup(result.responseText);
		}
	});
}

function PrevRecords() {
	$.ajax({
		url: "/api/prevRecords",
		type: "GET",
		data: {
			"dbPath": currentDBPath,
		},
		success: function(result){
			result = JSON.parse(result);

			ShowTree(result);
		},
		error: function(result) {
			ShowPopup(result.responseText);
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
			$("#currentPath").html("<i>" + result.bucketsPath + "<\/i> ");

			ShowTree(result);
		},
		error: function(result) {
			ShowPopup(result.responseText);
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
	// Removing last space
	currentPath = currentPath.slice(0, currentPath.length - 1);
	$("#recordPath").html(currentData[number].key + " – <i>" + currentPath + "<\/i>");
	$("#recordValue").html(currentData[number].value);
}

function ShowTree(data) {
	var result = "";
	if (data.prevRecords) {
		result += prevRecordersButtonTemplate;
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
		result += nextRecordersButtonTemplate;
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

// Popup
function ShowPopup(message) {
	$("#popupMessage").html(message);
	$("#popup").addClass("popup_animation");
}

function HidePopup() {
	$("#popup").removeClass("popup_animation");
}

// Modal
function ShowModal() {
	const template = `<option value="{0}">`;

	var sortedPaths = getPaths();

	var options = "";
	for (var i = 0; i < sortedPaths.length && i < 5; i++) {
		options += template.format(sortedPaths[i]);
	}

	$("#paths").html(options);
	$("#modal").css("display", "block");
	$("#DBPath").focus();
}

function HideModal() {
	$("#modal").css("display", "none");
	$("#dbPathsList").css("display", "none");
}


/* Secondary functions */
window.onclick = function(event) {
    if (event.target == modal) {
		HideModal();
	}
	if (event.target == dbListBackground) {
		$("#dbListBackground").css("display", "none");
		$("#dbList").removeClass("db_list_animation");
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
