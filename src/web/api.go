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
// Return:
// {
// 	"dbPath": str
// }
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
	response := struct {
		DBPath string `json:"dbPath"`
	}{dbPath}
	json.NewEncoder(w).Encode(response)
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
//	"recordsAmount": int,
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
		PrevBucket    bool        `json:"prevBucket"`
		PrevRecords   bool        `json:"prevRecords"`
		NextRecords   bool        `json:"nextRecords"`
		Path          string      `json:"bucketsPath"`
		RecordsAmount int         `json:"recordsAmount"`
		Records       []db.Record `json:"records"`
	}{
		data.PrevBucket,
		data.PrevRecords,
		data.NextRecords,
		data.Path,
		data.RecordsAmount,
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
//	"recordsAmount": int,
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
		PrevBucket    bool        `json:"prevBucket"`
		PrevRecords   bool        `json:"prevRecords"`
		NextRecords   bool        `json:"nextRecords"`
		Path          string      `json:"bucketsPath"`
		RecordsAmount int         `json:"recordsAmount"`
		Records       []db.Record `json:"records"`
	}{
		data.PrevBucket,
		data.PrevRecords,
		data.NextRecords,
		data.Path,
		data.RecordsAmount,
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
//	"recordsAmount": int,
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
		PrevBucket    bool        `json:"prevBucket"`
		PrevRecords   bool        `json:"prevRecords"`
		NextRecords   bool        `json:"nextRecords"`
		Path          string      `json:"bucketsPath"`
		RecordsAmount int         `json:"recordsAmount"`
		Records       []db.Record `json:"records"`
	}{
		data.PrevBucket,
		data.PrevRecords,
		data.NextRecords,
		data.Path,
		data.RecordsAmount,
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
//    "dbPath": "",
// 	  "size": 0
// 	},
// ]
//
func databasesList(w http.ResponseWriter, r *http.Request) {
	list := dbs.GetDBsList()

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
//	"recordsAmount": int,
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
		DB            dbs.DBInfo  `json:"db"`
		PrevBucket    bool        `json:"prevBucket"`
		PrevRecords   bool        `json:"prevRecords"`
		NextRecords   bool        `json:"nextRecords"`
		Path          string      `json:"bucketsPath"`
		RecordsAmount int         `json:"recordsAmount"`
		Records       []db.Record `json:"records"`
	}{
		dbs.DBInfo{
			Name:   info.Name,
			DBPath: info.DBPath,
			Size:   info.Size,
		},
		data.PrevBucket,
		data.PrevRecords,
		data.NextRecords,
		data.Path,
		data.RecordsAmount,
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
//	"recordsAmount": int,
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
		PrevBucket    bool        `json:"prevBucket"`
		PrevRecords   bool        `json:"prevRecords"`
		NextRecords   bool        `json:"nextRecords"`
		RecordsAmount int         `json:"recordsAmount"`
		Records       []db.Record `json:"records"`
	}{
		data.PrevBucket,
		data.PrevRecords,
		data.NextRecords,
		data.RecordsAmount,
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
//	"recordsAmount": int,
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
		PrevBucket    bool        `json:"prevBucket"`
		PrevRecords   bool        `json:"prevRecords"`
		NextRecords   bool        `json:"nextRecords"`
		RecordsAmount int         `json:"recordsAmount"`
		Records       []db.Record `json:"records"`
	}{
		data.PrevBucket,
		data.PrevRecords,
		data.NextRecords,
		data.RecordsAmount,
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
//	"recordsAmount": int,
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
	r.ParseForm()
	dbPath := r.Form.Get("dbPath")
	text := r.Form.Get("text")
	mode := r.Form.Get("mode")

	records, path, recordsAmount, code, err := dbs.Search(dbPath, mode, text)
	if err != nil {
		returnError(w, err, "", code)
		return
	}

	response := struct {
		PrevBucket    bool        `json:"prevBucket"`
		PrevRecords   bool        `json:"prevRecords"`
		NextRecords   bool        `json:"nextRecords"`
		Path          string      `json:"bucketsPath"`
		RecordsAmount int         `json:"recordsAmount"`
		Records       []db.Record `json:"records"`
	}{
		false,
		false,
		false,
		path + " (Search \"" + text + "\")",
		recordsAmount,
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

// addBucket
//
// Params: dbPath, bucket
// Return: -
//
func addBucket(w http.ResponseWriter, r *http.Request) {
	dbPath := r.Form.Get("dbPath")
	bucket := r.Form.Get("bucket")

	code, err := dbs.AddBucket(dbPath, bucket)
	if err != nil {
		returnError(w, err, "", code)
		return
	}

	w.WriteHeader(code)
}

// deleteBucket
//
// Params: dbPath, bucket (int URI)
// Return: -
//
func deleteBucket(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	dbPath := r.Form.Get("dbPath")
	bucket := r.Form.Get("bucket")

	code, err := dbs.DeleteBucket(dbPath, bucket)
	if err != nil {
		returnError(w, err, "", code)
		return
	}

	w.WriteHeader(code)
}

// addRecord
//
// Params: dbPath, key, value
// Return: -
//
func addRecord(w http.ResponseWriter, r *http.Request) {
	dbPath := r.Form.Get("dbPath")
	key := r.Form.Get("key")
	value := r.Form.Get("value")

	code, err := dbs.AddRecord(dbPath, key, value)
	if err != nil {
		returnError(w, err, "", code)
		return
	}

	w.WriteHeader(code)
}

// editRecord
//
// Params: dbPath, oldKey, newKey, newValue
// Return: -
//
func editRecord(w http.ResponseWriter, r *http.Request) {
	dbPath := r.Form.Get("dbPath")
	oldKey := r.Form.Get("oldKey")
	newKey := r.Form.Get("newKey")
	newValue := r.Form.Get("newValue")

	code, err := dbs.EditRecord(dbPath, oldKey, newKey, newValue)
	if err != nil {
		returnError(w, err, "", code)
		return
	}

	w.WriteHeader(code)
}

// deleteRecord
//
// Params: dbPath, key (int URI)
// Return: -
//
func deleteRecord(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	dbPath := r.Form.Get("dbPath")
	key := r.Form.Get("key")

	code, err := dbs.DeleteRecord(dbPath, key)
	if err != nil {
		returnError(w, err, "", code)
		return
	}

	w.WriteHeader(code)
}
