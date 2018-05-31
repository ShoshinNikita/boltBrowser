/* Templates */
// Reset margin for last element
const addMenuHTML = `<input type="button" class="popup_button" onclick="ShowAddModal('bucket');" value="Add bucket">
<input type="button" class="popup_button" style="margin: auto;" onclick="ShowAddModal('record');" value="Add record">`;

const bucketMenuHTML = `<input type="button" class="popup_button" style="margin: auto;" onclick="DeleteBucket('{0}');" value="Delete">`;

const recordMenuHTML = `<input type="button" class="popup_button" onclick="ShowEditModal('record', '{0}');" value="Edit">
<input type="button" class="popup_button" style="margin: auto;" onclick="DeleteRecord('{0}');" value="Delete">`;

const addBucketTemplate = `
<input id="newBucketName" required type="text" placeholder="Bucket" style="margin-bottom: 5px; width: 100%;">
<br>
<input type="submit" class="button" onclick="AddBucket();" value="Add">`;

const addRecordTemplate = `
<input id="newRecordKey" type="text" placeholder="Key" style="margin-bottom: 5px; width: 100%;">
<br>
<input id="newRecordValue" type="text" placeholder="Value" style="margin-bottom: 5px; width: 100%;">
<br>
<input type="submit" class="button" onclick="AddRecord();" value="Add">`;

const editRecordTemplate = `
<div style="margin-bottom: 10px;">Editing record "{0}"</div>
<input id="newRecordKey" type="text" placeholder="Key (leave empty if don't want to edit key)" style="margin-bottom: 5px; width: 100%; box-sizing: border-box;">
<textarea id="newRecordValue" type="text" placeholder="Value" style="resize: none; margin-bottom: 5px; width: 100%; height: 150px; box-sizing: border-box;"></textarea>
<input type="submit" class="button" onclick="EditRecord('{0}');" value="Edit">
`


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
function showPopupMenu(x, y, html) {
	$("#popupMenu").css("top", y + 10 + "px");
	$("#popupMenu").css("left", x + 10 + "px");
	$("#popupMenu").html(html);
	$("#popupMenu").css("visibility", "visible");
}

function showBucketMenu(event) {
	// Reset margin for last element
	var html =  bucketMenuHTML.format(event.target.innerHTML);
	showPopupMenu(event.clientX, event.clientY, html);
	return false;
}

function showRecordMenu(event) {
	// Reset margin for last element
	var html = recordMenuHTML.format(event.target.innerHTML)
	showPopupMenu(event.clientX, event.clientY, html);
	return false;
}

function showAddMenu(event) {
	// Magic. Program shows addMenu only if target is dbTreeWrapper and isn't anything else
	if ((event.target == dbTreeWrapper || event.target == dbTree) && event.which == 3) {
		// Show menu only if db was chosen
		if (currentDBPath != "") {
			showPopupMenu(event.clientX, event.clientY, addMenuHTML);
			return false;
		}
	}
}

// AddModal
function ShowAddModal(type) {
	if (type == "bucket") {
		$("#addItemWindow").html(addBucketTemplate);
		$("#addItemWindowBackground").css("display", "block");
	} else if (type == "record") {
		$("#addItemWindow").html(addRecordTemplate);
		$("#addItemWindowBackground").css("display", "block");
	}
}

function ShowEditModal(type, target) {
	if (type == "bucket") {
		// Nothing
	} else if (type == "record") {
		var html = editRecordTemplate.format(target)
		$("#addItemWindow").html(html);
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