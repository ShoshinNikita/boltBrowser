const buttonTemplate = `<div>
	<input type="button" class="db_button" value="{0}" onclick="ChooseDB('{1}')" title="Choose">
	<i class="material-icons btn" style="float: right; margin-right: 1vw; font-size: 30px !important;" title="Close" onclick="CloseDB('{1}')">close<\/i>
	<br>
	<br>
<\/div>`

const recordTemplate = `<div>
	<i class="material-icons" icon>assignment<\/i>
	<span class="record" onclick="ShowFullRecord({0});"><b>{1}<\/b>:<\/span> {2}
<\/div>`

const bucketTemplate = `<div>
	<i class="material-icons" icon>folder<\/i>
	<span class="bucket" onclick="Next('{0}');"><b>{0}<\/b><\/span>
<\/div>`

const backButton = `<div>
	<i class="material-icons btn" icon onclick="Back();" title="Back">more_horiz<\/i>
<\/div>`

const fullRecordTemplate = `<div>
	<b>Key:<\/b> {0}
<\/div>
<br>
<div>
	<b>Value:<\/b> {1}
<\/div>
<br>`

const nextRecordersButtonTemplate = `<div>
<i class="material-icons" icon>arrow_forward_ios<\/i>
<span style="cursor: pointer;" onclick="NextRecords();"><b>Next page<\/b><\/span>
<\/div>`

const prevRecordersButtonTemplate = `<div>
<i class="material-icons" icon>arrow_back_ios<\/i>
<span style="cursor: pointer;" onclick="PrevRecords();"><b>Previous page<\/b><\/span>
<\/div>`

var currentDBPath = ""
var currentData = null


function getError(result) {
	return errorMessageTemplate.format(result.status, result.responseText);
}

function ShowPopup(message) {
	$("#popupMessage").html(message);	
	$("#popup").addClass("popup_animation")
}

function HidePopup() {
	$("#popup").removeClass("popup_animation")
}

function ShowTree(data) {
	var result = ""
	if (data.prevRecords) {
		result += prevRecordersButtonTemplate
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
				value += "..."
			}
			result += recordTemplate.format(i, records[i].key, value);
		}
	}

	if (data.nextRecords) {
		result += nextRecordersButtonTemplate
	}
	$("#db_tree").html(result);

	document.getElementById("db_tree_wrapper").scrollTop = 0;
}

function OpenDB() {
	var dbPath = prompt("Please, enter the path to the database")
	$.ajax({
		url: "/api/databases",
		type: "POST",
		data: {
			"dbPath": dbPath
		},
		success: function(result){
			ShowDBList()
		},
		error: function(result) {
			ShowPopup(result.responseText)
		}
	})
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
				$("#dbName").html("<i>Name:<\/i> ?")
				$("#dbPath").html("<i>Path:<\/i> ?")
				$("#dbSize").html("<i>Size:<\/i> ?")
				$("#db_tree").html("")
				$("#currentPath").html("")
				$("#record_data").html("")
				$("#searchBox").css("visibility", "hidden")
				currentDBPath = ""
			}
			ShowDBList()
		},
		error: function(result) {
			ShowPopup(result.responseText)
		}
	})
}

function ShowDBList() {
	$.ajax({
		url: "/api/databases",
		type: "GET",
		success: function(result){
			allDB = JSON.parse(result)
			var result = ""
			for (i in allDB) {
				result += buttonTemplate.format(allDB[i].name, allDB[i].dbPath)
			}
			$("#list").html(result)
		},
		error: function(result) {
			ShowPopup(result.responseText)
		}
	})
}

function ChooseDB(dbPath) {
	currentDBPath = dbPath
	$.ajax({
		url: "/api/current",
		type: "GET",
		data: {
			"dbPath": dbPath,
		},
		success: function(result){
			result = JSON.parse(result)
			var path = result.db.dbPath
			if (path.length > 40) {
				path = path.substring(0, 20) + "..." + path.substring(path.length - 20, path.length)
			}

			$("#dbName").html("<i>Name:<\/i> " + result.db.name)
			$("#dbPath").html("<i>Path:<\/i> " + path)
			$('#dbPath').prop("title", result.db.dbPath);
			$("#dbSize").html("<i>Size:<\/i> " + result.db.size / 1024 + " Kb")
			$("#currentPath").html("<i>" + result.bucketsPath + "<\/i> ")
			$("#record_data").html("")
			$("#searchBox").css("visibility", "visible")
		
			ShowTree(result)
		},
		error: function(result) {
			ShowPopup(result.responseText)
		}
	})
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
			result = JSON.parse(result)
			$("#currentPath").html("<i>" + result.bucketsPath + "<\/i> ")
			$("#record_data").html("")
			ShowTree(result)
		},
		error: function(result) {
			ShowPopup(result.responseText)
		}
	})
}

function Back() {
	$.ajax({
		url: "/api/back",
		type: "GET",
		data: {
			"dbPath": currentDBPath,
		},
		success: function(result){
			result = JSON.parse(result)
			$("#currentPath").html("<i>" + result.bucketsPath + "<\/i> ")
			$("#record_data").html("")
			ShowTree(result)
		},
		error: function(result) {
			ShowPopup(result.responseText)
		}
	})
}

function NextRecords() {
	$.ajax({
		url: "/api/nextRecords",
		type: "GET",
		data: {
			"dbPath": currentDBPath,
		},
		success: function(result){
			result = JSON.parse(result)
			$("#record_data").html("")

			ShowTree(result)
		},
		error: function(result) {
			ShowPopup(result.responseText)
		}
	})
}

function PrevRecords() {
	$.ajax({
		url: "/api/prevRecords",
		type: "GET",
		data: {
			"dbPath": currentDBPath,
		},
		success: function(result){
			result = JSON.parse(result)
			$("#record_data").html("")

			ShowTree(result)
		},
		error: function(result) {
			ShowPopup(result.responseText)
		}
	})
}

function ShowFullRecord(number) {
	var result = fullRecordTemplate.format(currentData[number].key, currentData[number].value)
	$("#record_data").html(result)
}

function Search() {
	var text = $("#searchText").val()

	if (text == "") {
		ShowPopup("Error: empty input")
		return
	}

	if ($("#regexMode").prop("checked")) {
		$.ajax({
			url: "/api/searchRegex",
			type: "GET",
			data: {
				"dbPath": currentDBPath,
				"expr": text,
			},
			success: function(result){
				result = JSON.parse(result)
				$("#record_data").html("")
				$("#currentPath").html("<i>" + result.bucketsPath + "<\/i> ")

				ShowTree(result)
			},
			error: function(result) {
				ShowPopup(result.responseText)
			}
		})
	} else {
		$.ajax({
			url: "/api/search",
			type: "GET",
			data: {
				"dbPath": currentDBPath,
				"needle": text,
			},
			success: function(result){
				result = JSON.parse(result)
				$("#record_data").html("")
				$("#currentPath").html("<i>" + result.bucketsPath + "<\/i> ")

				ShowTree(result)
			},
			error: function(result) {
				ShowPopup(result.responseText)
			}
		})
	}
}


String.prototype.format = function () {
	var a = this;
	for (var k in arguments) {
		a = a.replace(new RegExp("\\{" + k + "\\}", 'g'), arguments[k]);
	}
	return a;
}