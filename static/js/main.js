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

	var $closeBtn = $("<i>", {class: "material-icons btn", style: "float: right; margin-right: 10px; font-size: 30px !important;", title: "Close"}).text("close").
		click({dbPath: dbPath}, function(event){
			CloseDB(event.data.dbPath);
		});

	return $("<div>").append($input).append($closeBtn);
}

function getRecordButton(key, value) {
	var $icon = $("<i>", {class: "material-icons"}).text("assignment");
	var $key = $("<span>", {class: "record", id: "key", style: "font-weight: bold;"}).html(key).
		click({key: key}, function(event) {
			ShowFullRecord(event.data.key);
		});
	var $value = $("<span>", {id: "value"}).html(" – " + value);

	return $("<div>", {style: "display: table;"}).append($icon).append($key).append($value);
}

function getBucketButton(key) {
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

// For creating list of paths for deleting
function getPathForDeleting(path) {
	var $path = $("<span>").text(path);
	var $btn = $("<i>", {class: "material-icons btn", style: "float: right; font-size: 22px !important; vertical-align: middle;", title: "Delete"}).text("close").
		click({path: path}, function(event) {
			DeletePath(event.data.path);
		});

	return $("<div>", {style: "margin-bottom: 10px; text-align: left;"}).append($path).append($btn);
}

// Write mode only
function getAddMenu() {
	var $bucket = $("<input>", {type: "button", class: "popup_button", value: "Add bucket"}).
		click(function() {
			ShowAddModal("bucket");
		});
	var $record =  $("<input>", {type: "button", class: "popup_button", value: "Add record"}).
	click(function() {
		ShowAddModal("record");
	});

	return $("<div>").append($bucket).append($record);
}

function getBucketMenu(bucketKey) {
	var $editBtn = $("<input>", {type: "button", class: "popup_button", value: "Edit name"}).
		click({key: bucketKey}, function(event) {
			ShowEditModal("bucket", event.data.key);
		});
	var $delBtn =  $("<input>", {type: "button", class: "popup_button", value: "Delete", style: "margin: auto;"}).
	click({key: bucketKey}, function(event) {
		DeleteBucket(event.data.key);
	});

	return $("<div>").append($editBtn).append($delBtn);
}

function getRecordMenu(recordKey) {
	var $editBtn = $("<input>", {type: "button", class: "popup_button", value: "Edit"}).
		click({key: recordKey}, function(event) {
			ShowEditModal("record", event.data.key);
		});
	var $delBtn =  $("<input>", {type: "button", class: "popup_button", value: "Delete", style: "margin: auto;"}).
	click({key: recordKey}, function(event) {
		DeleteRecord(event.data.key);
	});

	return $("<div>").append($editBtn).append($delBtn);
}

function getAddBucketWindow() {
	var $nameInput = $("<input>", {id: "newBucketName", "type": "text", placeholder: "Bucket", style: "margin-bottom: 5px; width: 100%;"}).prop("required", true);
	var $btn = $("<input>", {type: "submit", "class": "button", value: "Add"}).
		click(function() {
			AddBucket();
		});

	return $("<div>").append($nameInput).append($btn);
}

function getEditBucketWindow(bucketName) {
	var $title = $("<div>", {style: "margin-bottom: 10px;"}).text("The old name: " + bucketName);
	var $nameInput = $("<input>", {id: "newName", type: "text", placeholder: "New name", style: "margin-bottom: 5px; width: 100%;"}).prop("required", true);
	var $btn = $("<input>", {type: "submit", class: "button", value: "Edit"}).
		click({key: bucketName}, function(event) {
			EditBucketName(event.data.key);
		});

	return $("<div>").append($title).append($nameInput).append($("<br>")).append($btn);
}

function getAddRecordWindow() {
	var $key = $("<input>", {id: "newRecordKey", type: "text", placeholder: "Key", style: "margin-bottom: 5px; width: 100%;"});
	var $br = $("<br>");
	var $value = $("<input>", {id: "newRecordValue", type: "text", placeholder: "Value", style: "margin-bottom: 5px; width: 100%;"});
	var $btn = $("<input>", {type: "submit", class: "button", value: "Add"}).
		click({key: bucketName}, function(event) {
			AddRecord();
		});

	return $("<div>").append($key).append($br).append($value).append($br).append($btn);
}

function getEditRecordWindow(recordKey, recordValue) {
	var $title = $("<div>", {style: "margin-bottom: 10px;"}).text("Editing of record \"" + recordKey + "\"");
	var $newKey = $("<input>", {id: "newRecordKey", type: "text", placeholder: "Key (leave empty if don't want to edit key)", style: "margin-bottom: 5px; width: 100%; box-sizing: border-box;", value: recordKey});
	var $newValue = $("<textarea>", {id: "newRecordValue", placeholder: "Value", style: "resize: none; margin-bottom: 5px; width: 100%; height: 150px; box-sizing: border-box;"}).val(recordValue);
	var $btn = $("<input>", {type: "submit", class: "button", value: "Edit"}).
		click({key: recordKey}, function(event) {
			EditRecord(event.data.key);
		});
	
	return $("<div>").append($title).append($newKey).append($newValue).append($btn);
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
