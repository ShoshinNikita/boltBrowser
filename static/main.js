const buttonTemplate = `
<div>
	<input type="button" class="db_button" value="{0}" onclick="ChooseDB('{1}')" title="Choose">
	<input type="button" value="X" style="float: right; margin-right: 1vw; height: 4vh; width: 15%;" onclick="CloseDB('{1}')" title="Close">
</div>`

const recordTemplate = `
<div>
	<i class="material-icons" icon>assignment</i>
	{0}: {1}
</div>
`
const bucketTemplate = `
<div>
	<i class="material-icons" icon>folder</i>
	<span onclick="Next(currentDBPath, '{0}');">{0}</span>
</div>
`

const backButton = `
<div>
	<i class="material-icons" icon onclick="Back(currentDBPath);">more_horiz</i>
</div>
`

var currentDBPath = ""

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
			showDBList()
		},
		error: function(result) {
			console.log(result.responseText)
		}
	})
}

function showDBList() {
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
			console.log(result.responseText)
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
			console.log(result.responseText)
			showDBList()
		},
		error: function(result) {
			console.log(result.responseText)
		}
	})
}

function ShowTree(records) {
	var result = backButton
	for (i in records) {
		if (records[i].type == "bucket") {
			result += bucketTemplate.format(records[i].key)
		} else if (records[i].type == "record") {
			result += recordTemplate.format(records[i].key, records[i].value)
		}
	}
	$("#db_tree").html(result)
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
			ShowTree(result)
		},
		error: function(result) {
			console.log(result.responseText)
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
			ShowTree(result)
		},
		error: function(result) {
			console.log(result.responseText)
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
			$("#dbName").html("Name: " + result.name)
			$("#dbPath").html("Path: " + result.path)
			$("#dbSize").html("Size: " + result.size + " Kb")
			ShowTree(result.records)
		},
		error: function(result) {
			console.log(result.responseText)
		}
	})
}

String.prototype.format = function () {
	var a = this;
	for (var k in arguments) {
		a = a.replace(new RegExp("\\{" + k + "\\}", 'g'), arguments[k]);
	}
	return a;
}