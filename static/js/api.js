function OpenDB() {
	var dbPath = $("#DBPath").val();
	if (dbPath == "" ) {
		ShowErrorPopup("Error: path is empty");
		return;
	}

	var readOnly = false;
	if ($("#readOnlyCheckbox").prop("checked")) {
		readOnly = true;
	}

	$("#DBPath").val("");
	$.ajax({
		url: "/api/databases",
		type: "POST",
		data: {
			"dbPath": dbPath,
			"readOnly": readOnly
		},
		success: function(result){
			// Uncheck the checkbox
			$("#readOnlyCheckbox").prop("checked", false);

			result= JSON.parse(result);
			PutIntoLS(result.dbPath);
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
			PutIntoLS(result.dbPath);
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
			if (dbPath == currentDB.dbPath) {
				$("#dbName").html("<i>Name:<\/i> ?");
				$("#dbPath").html("<i>Path:<\/i> ?");
				$("#dbSize").html("<i>Size:<\/i> ?");
				$("#dbTree").html("");
				$("#currentPath").html("");
				$("#recordsAmount").html("");
				$("#recordPath").html("?");
				$("#recordValue").html("?");
				$("#searchBox").css("visibility", "hidden");
				currentDB = {};
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
	$.ajax({
		url: "/api/buckets/current",
		type: "GET",
		data: {
			"dbPath": dbPath,
		},
		success: function(result){
			result = JSON.parse(result);
			currentDB = result.db;

			var getDiv = function(field, text) {
				var $field = $("<i>").text(field + ": ");
				return $("<div>").append($field).append(text);
			}
			
			$("#dbName").empty();
			$("#dbName").append(getDiv("Name", result.db.name));

			$("#dbPath").empty();
			$("#dbPath").html(getDiv("Path", result.db.dbPath));

			$("#dbSize").empty();
			$("#dbSize").html(getDiv("Size", result.db.size / 1024 + " Kb"));

			$("#dbMode").empty();
			var mode = "Read & Write"
			if (result.db.readOnly) {
				mode = "Read-Only"
			}
			$("#dbMode").append(getDiv("Mode", mode));

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
			"dbPath": currentDB.dbPath,
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
			"dbPath": currentDB.dbPath,
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
			"dbPath": currentDB.dbPath,
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
			"dbPath": currentDB.dbPath,
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
		ChooseDB(currentDB.dbPath);
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
			"dbPath": currentDB.dbPath,
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

// Write mode only
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
			"dbPath": currentDB.dbPath,
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
			"dbPath": currentDB.dbPath,
			"oldName": oldName,
			"newName": newName,
		},
		success: function(result){
			HideAddModal();
			ShowDonePopup();
			ChooseDB(currentDB.dbPath);
		},
		error: function(result) {
			ShowErrorPopup(result.responseText);
		}
	});
}

function DeleteBucket(bucket) {
	var query = $.param({"dbPath": currentDB.dbPath, "bucket": bucket})
	$.ajax({
		url: "/api/buckets" + "?" + query,
		type: "DELETE",
		success: function(result){
			ShowDonePopup();
			ChooseDB(currentDB.dbPath)
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
			"dbPath": currentDB.dbPath,
			"key": key,
			"value": value,
		},
		success: function(result){
			HideAddModal();
			ShowDonePopup();
			ChooseDB(currentDB.dbPath)
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
			"dbPath": currentDB.dbPath,
			"oldKey": oldKey,
			"newKey": newKey,
			"newValue": newValue,
		},
		success: function(result){
			HideAddModal();
			ShowDonePopup();
			ChooseDB(currentDB.dbPath)
		},
		error: function(result) {
			ShowErrorPopup(result.responseText);
		}
	});
}

function DeleteRecord(key) {
	var query = $.param({"dbPath": currentDB.dbPath, "key": key})
	$.ajax({
		url: "/api/records" + "?" + query,
		type: "DELETE",
		success: function(result){
			ShowDonePopup();
			ChooseDB(currentDB.dbPath)
		},
		error: function(result) {
			ShowErrorPopup(result.responseText);
		}
	});
}