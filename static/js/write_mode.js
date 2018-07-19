// Binding of events
// Write mode only
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