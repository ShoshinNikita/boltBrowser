package web

import (
	"db"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"dbs"
)

// openDB open db. It also adds db.DBApi to allDB
//
// Params: dbPath
// Return: -
//
func openDB(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	dbPath := r.Form.Get("dbPath")

	// From C:\\users\\help (or C:\users\help) to C:/users/help
	reg := regexp.MustCompile(`\\\\|\\`)
	dbPath = reg.ReplaceAllString(dbPath, "/")

	dbName, code, err := dbs.OpenDB(dbPath)
	if err != nil {
		returnError(w, err, "", code)
		return
	}

	fmt.Printf("[INFO] DB \"%s\" (%s) was opened\n", dbName, dbPath)
	w.WriteHeader(code)
}

// Params: dbPath
// Return: -
//
func closeDB(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	dbPath := r.Form.Get("dbPath")

	code, err := dbs.CloseDB(dbPath)
	if err != nil {
		returnError(w, err, "", code)
	}

	w.WriteHeader(code)
}

// next returns records from bucket with according to the name
//
// Params: dbPath, bucket
// Return:
// {
// 	"prevBucket": bool,
//  "prevRecords": bool,
//  "nextRecords": bool,
//  "bucketsPath": string,
// 	"records": [
// 	  {
// 		"type": "",
// 		"key": "",
// 		"value": ""
// 	  },
// 	]
// }
//
func next(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	dbPath := r.Form.Get("dbPath")
	nextBucket := r.Form.Get("bucket")

	data, code, err := dbs.NextBucket(dbPath, nextBucket)
	if err != nil {
		returnError(w, err, "", code)
		return
	}

	response := struct {
		PrevBucket  bool        `json:"prevBucket"`
		PrevRecords bool        `json:"prevRecords"`
		NextRecords bool        `json:"nextRecords"`
		Path        string      `json:"bucketsPath"`
		Records     []db.Record `json:"records"`
	}{
		data.PrevBucket,
		data.PrevRecords,
		data.NextRecords,
		data.Path,
		data.Records,
	}

	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}

// back returns records from previous directory
//
// Params: dbPath
// Return:
// {
// 	"prevBucket": bool,
//  "prevRecords": bool,
//  "nextRecords": bool,
//  "bucketsPath": string,
// 	"records": [
//   {
// 	   "type": "",
// 	   "key": "",
// 	   "value": ""
// 	 },
// 	]
// }
//
func back(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	dbPath := r.Form.Get("dbPath")

	data, code, err := dbs.PrevBucket(dbPath)
	if err != nil {
		returnError(w, err, "", code)
		return
	}

	response := struct {
		PrevBucket  bool        `json:"prevBucket"`
		PrevRecords bool        `json:"prevRecords"`
		NextRecords bool        `json:"nextRecords"`
		Path        string      `json:"bucketsPath"`
		Records     []db.Record `json:"records"`
	}{
		data.PrevBucket,
		data.PrevRecords,
		data.NextRecords,
		data.Path,
		data.Records,
	}

	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)

}

// root returns records from root of db
//
// Params: dbPath
// Return:
// {
// 	"prevBucket": bool,
//  "prevRecords": bool,
//  "nextRecords": bool,
//  "bucketsPath": string,
// 	"records": [
// 	 {
// 	   "type": "",
// 	   "key": "",
// 	   "value": ""
// 	 },
// 	]
// }
//
func root(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	dbPath := r.Form.Get("dbPath")

	data, code, err := dbs.GetRoot(dbPath)
	if err != nil {
		returnError(w, err, "", code)
		return
	}

	response := struct {
		PrevBucket  bool        `json:"prevBucket"`
		PrevRecords bool        `json:"prevRecords"`
		NextRecords bool        `json:"nextRecords"`
		Path        string      `json:"bucketsPath"`
		Records     []db.Record `json:"records"`
	}{
		data.PrevBucket,
		data.PrevRecords,
		data.NextRecords,
		data.Path,
		data.Records,
	}

	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}

// databasesList return list of dbs
//
// Params: -
// Return:
// [
//	{
// 	  "name": "",
//    "path": "",
// 	  "size": 0
// 	},
// ]
//
func databasesList(w http.ResponseWriter, r *http.Request) {
	list := dbs.GetDBsList()
	// TODO check DBInfo

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(list)
}

// current returns records in current bucket
//
// Params: dbPath
// Return:
// {
//  "db" : {
//    "name": "",
// 	  "dbPath": "",
//    "size": 0,
//  },
//  "prevBucket": bool,
//  "prevRecords": bool,
//  "nextRecords": bool,
//  "bucketsPath": string,
// 	"records": [
// 	  {
// 	    "type": "",
// 		"key": "",
// 		"value": ""
// 	  },
// 	]
// }
//
func current(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	dbPath := r.Form.Get("dbPath")

	info, data, code, err := dbs.GetCurrent(dbPath)
	if err != nil {
		returnError(w, err, "", code)
		return
	}

	response := struct {
		DB          dbs.DBInfo  `json:"db"`
		PrevBucket  bool        `json:"prevBucket"`
		PrevRecords bool        `json:"prevRecords"`
		NextRecords bool        `json:"nextRecords"`
		Path        string      `json:"bucketsPath"`
		Records     []db.Record `json:"records"`
	}{
		dbs.DBInfo{
			info.Name,
			info.DBPath,
			info.Size,
		},
		data.PrevBucket,
		data.PrevRecords,
		data.NextRecords,
		data.Path,
		data.Records,
	}

	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}

// nextRecords
//
// Params: dbPath
// Return:
// {
//  "prevBucket": bool,
//  "prevRecords": bool,
//  "nextRecords": bool,
// 	"records": [
// 	  {
// 	    "type": "",
// 		"key": "",
// 		"value": ""
// 	  },
// 	]
// }
//
func nextRecords(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	dbPath := r.Form.Get("dbPath")

	data, code, err := dbs.GetNextRecords(dbPath)
	if err != nil {
		returnError(w, err, "", code)
		return
	}

	response := struct {
		PrevBucket  bool        `json:"prevBucket"`
		PrevRecords bool        `json:"prevRecords"`
		NextRecords bool        `json:"nextRecords"`
		Records     []db.Record `json:"records"`
	}{
		data.PrevBucket,
		data.PrevRecords,
		data.NextRecords,
		data.Records,
	}

	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}

// prevRecords
//
// Params: dbPath
// Return:
// {
//  "prevBucket": bool,
//  "prevRecords": bool,
//  "nextRecords": bool,
//  "bucketsPath": string,
// 	"records": [
// 	  {
// 	    "type": "",
// 		"key": "",
// 		"value": ""
// 	  },
// 	]
// }
//
func prevRecords(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	dbPath := r.Form.Get("dbPath")

	data, code, err := dbs.GetPrevRecrods(dbPath)
	if err != nil {
		returnError(w, err, "", code)
		return
	}

	response := struct {
		PrevBucket  bool        `json:"prevBucket"`
		PrevRecords bool        `json:"prevRecords"`
		NextRecords bool        `json:"nextRecords"`
		Records     []db.Record `json:"records"`
	}{
		data.PrevBucket,
		data.PrevRecords,
		data.NextRecords,
		data.Records,
	}

	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}

// search
//
// Params: dbPath, text, mode ("regex" or "plain")
// Return:
// {
//  "prevBucket": bool,
//  "prevRecords": bool,
//  "nextRecords": bool,
//  "bucketsPath": string,
// 	"records": [
// 	  {
// 	    "type": "",
// 		"key": "",
// 		"value": ""
// 	  },
// 	]
// }
//
func search(w http.ResponseWriter, r *http.Request) {
	const (
		regex = "regex"
		plain = "plain"
	)

	r.ParseForm()
	dbPath := r.Form.Get("dbPath")
	text := r.Form.Get("text")
	mode := r.Form.Get("mode")

	records, path, code, err := dbs.Search(dbPath, mode, text)
	if err != nil {
		returnError(w, err, "", code)
		return
	}

	response := struct {
		PrevBucket  bool        `json:"prevBucket"`
		PrevRecords bool        `json:"prevRecords"`
		NextRecords bool        `json:"nextRecords"`
		Path        string      `json:"bucketsPath"`
		Records     []db.Record `json:"records"`
	}{
		false,
		false,
		false,
		path + " (Search \"" + text + "\")",
		records,
	}

	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}

// returnError writes error into http.ResponseWriter and into terminal
func returnError(w http.ResponseWriter, err error, message string, code int) {
	var text string
	if message != "" && err != nil {
		text = fmt.Sprintf("Error: %s Message: %s", err.Error(), message)
	} else if message != "" {
		text = fmt.Sprintf("Message: %s", message)
	} else if err != nil {
		text = fmt.Sprintf("Error: %s", err.Error())
	} else {
		text = "Nothing"
	}

	fmt.Printf("[ERR] %s\n", text)

	http.Error(w, text, code)
}
