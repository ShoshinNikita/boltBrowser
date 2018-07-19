/* Global variables */
var currentDBPath = "";
// Dictionary. It keeps data like "key of a record": "value of a record"
var currentData = {};

/* Functions for getting html elements */
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