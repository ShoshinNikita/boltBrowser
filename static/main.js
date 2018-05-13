const buttonTemplate = `
<div>
	<input type="button" class="db_button" value="{0}" onclick="ChooseDB('{1}')" title="Choose">
	<i class="material-icons btn" style="float: right; margin-right: 1vw; font-size: 30px !important;" title="Close" onclick="CloseDB('{1}')">close<\/i>
	<br>
	<br>
<\/div>`

const recordTemplate = `
<div>
	<i class="material-icons" icon>assignment<\/i>
	<span class="record" onclick="ShowFullRecord({0});"><b>{1}<\/b>:<\/span> {2}
<\/div>`

const bucketTemplate = `
<div>
	<i class="material-icons" icon>folder<\/i>
	<span class="bucket" onclick="Next(currentDBPath, '{0}');"><b>{0}<\/b><\/span>
<\/div>`

const backButton = `
<div>
	<i class="material-icons btn" icon onclick="Back(currentDBPath);" title="Back">more_horiz<\/i>
<\/div>`

const fullRecordTemplate = `
<div>
	<b>Key:<\/b> {0}
<\/div>
<br>
<div>
	<b>Value:<\/b> {1}
<\/div>
<br>`

const errorMessageTemplate = `<b>Code:<\/b> {0}<br><b>Message:<\/b> {1}`

var currentDBPath = ""
var currentData = null

function getError(result) {
	return errorMessageTemplate.format(result.status, result.responseText);
}


function OpenDB() {
	var path = prompt("Please, enter the path to the database")
	$.ajax({
		url: "/api/openDB",
		type: "POST",
		data: {
			"path": path
		},
		success: function(result){
			result = JSON.parse(result);
			ShowDBList()
		},
		error: function(result) {
			ShowPopup(getError(result))
		}
	})
}

function CloseDB(path) {
	$.ajax({
		url: "/api/closeDB",
		type: "POST",
		data: {
			"path": path,
		},
		success: function(result){
			if (path == currentDBPath) {
				$("#dbName").html("<i>Name:<\/i> ?")
				$("#dbPath").html("<i>Path:<\/i> ?")
				$("#dbSize").html("<i>Size:<\/i> ?")
				$("#db_tree").html("")
				$("#record_data").html("")
				currentDBPath = ""
			}
			ShowDBList()
		},
		error: function(result) {
			ShowPopup(getError(result))
		}
	})
}

function ShowDBList() {
	$.ajax({
		url: "/api/databases",
		type: "GET",
		success: function(result){
			allDB = JSON.parse(result)
			console.log(allDB)
			var result = ""
			for (i in allDB) {
				result += buttonTemplate.format(allDB[i].name, allDB[i].path)
			}
			$("#list").html(result)
		},
		error: function(result) {
			ShowPopup(getError(result))
		}
	})
}

function ChooseDB(path) {
	currentDBPath = path
	$.ajax({
		url: "/api/current",
		type: "GET",
		data: {
			"path": path,
		},
		success: function(result){
			result = JSON.parse(result)
			$("#dbName").html("<i>Name:<\/i> " + result.name)
			$("#dbPath").html("<i>Path:<\/i> " + result.path)
			$("#dbSize").html("<i>Size:<\/i> " + result.size + " Kb")
			$("#record_data").html("")
			ShowTree(result)
		},
		error: function(result) {
			ShowPopup(getError(result))
		}
	})
}

function ShowTree(data) {
	var result = ""
	if (data.canBack) {
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
	$("#db_tree").html(result);
}

function Next(path, bucket) {
	$.ajax({
		url: "/api/next",
		type: "GET",
		data: {
			"path": path,
			"bucket": bucket
		},
		success: function(result){
			result = JSON.parse(result)
			$("#record_data").html("")
			ShowTree(result)
		},
		error: function(result) {
			ShowPopup(getError(result))
		}
	})
}

function Back(path) {
	$.ajax({
		url: "/api/back",
		type: "GET",
		data: {
			"path": path,
		},
		success: function(result){
			result = JSON.parse(result)
			$("#record_data").html("")
			ShowTree(result)
		},
		error: function(result) {
			ShowPopup(getError(result))
		}
	})
}

function ShowFullRecord(number) {
	console.log(currentData[number])
	var result = fullRecordTemplate.format(currentData[number].key, currentData[number].value)
	$("#record_data").html(result)
}


String.prototype.format = function () {
	var a = this;
	for (var k in arguments) {
		a = a.replace(new RegExp("\\{" + k + "\\}", 'g'), arguments[k]);
	}
	return a;
}