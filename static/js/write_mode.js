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


/* API */
function AddBucket() {
	var bucketName = $("#newBucketName").val();
	if (bucketName == "") {
		ShowErrorPopup("Error: bucketName is empty")
		return
	}

	$.ajax({
		url: "/api/buckets",
		type: "POST",
		data: {
			"dbPath": currentDBPath,
			"bucket": bucketName,
		},
		success: function(result){
			HideAddModal();
			ShowDonePopup();
			Next(bucketName);
		},
		error: function(result) {
			ShowErrorPopup(result.responseText);
		}
	});
}

function EditBucketName(oldName) {
	var newName = $("#newName").val();

	$.ajax({
		url: "/api/buckets",
		type: "PUT",
		data: {
			"dbPath": currentDBPath,
			"oldName": oldName,
			"newName": newName,
		},
		success: function(result){
			HideAddModal();
			ShowDonePopup();
			ChooseDB(currentDBPath);
		},
		error: function(result) {
			ShowErrorPopup(result.responseText);
		}
	});
}

function DeleteBucket(bucket) {
	var query = $.param({"dbPath": currentDBPath, "bucket": bucket})
	$.ajax({
		url: "/api/buckets" + "?" + query,
		type: "DELETE",
		success: function(result){
			ShowDonePopup();
			ChooseDB(currentDBPath)
		},
		error: function(result) {
			ShowErrorPopup(result.responseText);
		}
	});
}

function AddRecord() {
	var key = $("#newRecordKey").val();
	var value = $("#newRecordValue").val();
	if (key == "" || value == "") {
		ShowErrorPopup("Error: key or value is empty")
		return
	}

	$.ajax({
		url: "/api/records",
		type: "POST",
		data: {
			"dbPath": currentDBPath,
			"key": key,
			"value": value,
		},
		success: function(result){
			HideAddModal();
			ShowDonePopup();
			ChooseDB(currentDBPath)
		},
		error: function(result) {
			ShowErrorPopup(result.responseText);
		}
	});
}

function EditRecord(oldKey) {
	var newKey = $("#newRecordKey").val();
	if (newKey == "") {
		newKey = oldKey
	}

	var newValue = $("#newRecordValue").val();

	$.ajax({
		url: "/api/records",
		type: "PUT",
		data: {
			"dbPath": currentDBPath,
			"oldKey": oldKey,
			"newKey": newKey,
			"newValue": newValue,
		},
		success: function(result){
			HideAddModal();
			ShowDonePopup();
			ChooseDB(currentDBPath)
		},
		error: function(result) {
			ShowErrorPopup(result.responseText);
		}
	});
}

function DeleteRecord(key) {
	var query = $.param({"dbPath": currentDBPath, "key": key})
	$.ajax({
		url: "/api/records" + "?" + query,
		type: "DELETE",
		success: function(result){
			ShowDonePopup();
			ChooseDB(currentDBPath)
		},
		error: function(result) {
			ShowErrorPopup(result.responseText);
		}
	});
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
	showPopupMenu(event.clientX, event.clientY, getBucketMenu(event.target.innerText));
	return false;
}

function showRecordMenu(event) {
	showPopupMenu(event.clientX, event.clientY, getRecordMenu(event.target.innerText));
	return false;
}

function showAddMenu(event) {
	// Magic. Program shows addMenu only if target is dbTreeWrapper and isn't anything else
	if ((event.target == dbTreeWrapper || event.target == dbTree) && event.which == 3) {
		// Show menu only if db was chosen
		if (currentDBPath != "") {
			showPopupMenu(event.clientX, event.clientY, getAddMenu());
			return false;
		}
	}
}

// AddModal
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


/* Secondary */
$(document).ready(function() {
	// For hiding default context menu
	$("#dbTreeWrapper").attr("oncontextmenu", "return false;")

	$("#dbTreeWrapper").on("mousedown", function(event) {
		showAddMenu(event);
	});

	$("body").on("contextmenu", ".record", function(event) {
		showRecordMenu(event);
		return false;
	});
	
	$("body").on("contextmenu", ".bucket", function(event) {
		showBucketMenu(event);
		return false;
	});
});