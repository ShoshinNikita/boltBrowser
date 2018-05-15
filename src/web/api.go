package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"db"
)

// allDB keeps all opened databases. string – the path to the db
var allDB map[string]*db.BoltAPI

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

	// Check if db was opened
	if _, ok := allDB[dbPath]; ok {
		returnError(w, nil, "This DB was already opened", http.StatusBadRequest)
		return
	}

	newDB, err := db.Open(dbPath)
	if err != nil {
		returnError(w, err, "", http.StatusInternalServerError)
		return
	}

	allDB[dbPath] = newDB
	fmt.Printf("[INFO] DB \"%s\" (%s) was opened\n", newDB.Name, newDB.DBPath)
	w.WriteHeader(http.StatusOK)
}

// Params: dbPath
// Return: -
//
func closeDB(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	dbPath := r.Form.Get("dbPath")
	if _, ok := allDB[dbPath]; ok {
		dbName := allDB[dbPath].Name
		allDB[dbPath].Close()
		delete(allDB, dbPath)
		fmt.Printf("[INFO] DB \"%s\" (%s) was closed\n", dbName, dbPath)
	}
	w.WriteHeader(http.StatusOK)
}

// next returns records from bucket with according to the name
//
// Params: dbPath, bucket
// Return:
// {
// 	"prevBucket": bool,
//  "prevRecords": bool,
//  "nextRecords": bool,
//  "bucketsPath": [],
// 	"records": [
// 		{
// 			"type": "",
// 			"key": "",
// 			"value": ""
// 		},
// 	]
// }
//
func next(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	dbPath := r.Form.Get("dbPath")
	nextBucket := r.Form.Get("bucket")

	if _, ok := allDB[dbPath]; ok {
		records, bucketsPath, err := allDB[dbPath].Next(nextBucket)
		if err != nil {
			returnError(w, err, "", http.StatusInternalServerError)
			return
		}
		response := struct {
			PrevBucket bool        `json:"prevBucket"`
			Path       []string    `json:"bucketsPath"`
			Records    []db.Record `json:"records"`
		}{
			true,
			bucketsPath,
			records,
		}
		json.NewEncoder(w).Encode(response)
	} else {
		returnError(w, nil, "Bad path of db "+dbPath, http.StatusBadRequest)
	}
}

// back returns records from previous directory
//
// Params: dbPath
// Return:
// {
// 	"prevBucket": bool,
//  "prevRecords": bool,
//  "nextRecords": bool,
//  "bucketsPath": [],
// 	"records": [
// 		{
// 			"type": "",
// 			"key": "",
// 			"value": ""
// 		},
// 	]
// }
//
func back(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	dbPath := r.Form.Get("dbPath")

	if _, ok := allDB[dbPath]; ok {
		records, bucketsPath, err := allDB[dbPath].Back()
		if err != nil {
			returnError(w, err, "", http.StatusInternalServerError)
			return
		}
		response := struct {
			PrevBucket bool        `json:"prevBucket"`
			Path       []string    `json:"bucketsPath"`
			Records    []db.Record `json:"records"`
		}{
			func() bool { return len(bucketsPath) != 0 }(),
			bucketsPath,
			records,
		}
		json.NewEncoder(w).Encode(response)
	} else {
		returnError(w, nil, "Bad path of db "+dbPath, http.StatusBadRequest)
	}
}

// root returns records from root of db
//
// Params: dbPath
// Return:
// {
// 	"prevBucket": bool,
//  "prevRecords": bool,
//  "nextRecords": bool,
//  "bucketsPath": [],
// 	"records": [
// 		{
// 			"type": "",
// 			"key": "",
// 			"value": ""
// 		},
// 	]
// }
//
func root(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	dbPath := r.Form.Get("dbPath")

	if _, ok := allDB[dbPath]; ok {
		records, _, err := allDB[dbPath].GetRoot()
		if err != nil {
			returnError(w, err, "", http.StatusInternalServerError)
			return
		}
		response := struct {
			PrevBucket bool        `json:"prevBucket"`
			Path       []string    `json:"bucketsPath"`
			Records    []db.Record `json:"records"`
		}{
			false,
			[]string{},
			records,
		}
		json.NewEncoder(w).Encode(response)
	} else {
		returnError(w, nil, "Bad path of db "+dbPath, http.StatusBadRequest)
	}
}

// Params: -
// Return:
// [
//	{
// 		"name": "",
// 		"path": "",
// 		"size": 0
// 	},
// ]
//
func databasesList(w http.ResponseWriter, r *http.Request) {
	var list []db.BoltAPI
	for _, v := range allDB {
		list = append(list, *v)
	}
	json.NewEncoder(w).Encode(list)
}

// current returns records in current bucket
//
// Params: dbPath
// Return:
// {
// 	"name": "",
// 	"dbPath": "",
// 	"size": 0,
//  "prevBucket": bool,
//  "prevRecords": bool,
//  "nextRecords": bool,
//  "bucketsPath": [],
// 	"records": [
// 		{
// 			"type": "",
// 			"key": "",
// 			"value": ""
// 		},
// 	]
// }
//
func current(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	dbPath := r.Form.Get("dbPath")

	if _, ok := allDB[dbPath]; ok {
		records, bucketsPath, err := allDB[dbPath].GetCurrent()
		if err != nil {
			returnError(w, err, "", http.StatusInternalServerError)
			return
		}
		response := struct {
			*db.BoltAPI
			CanBack bool        `json:"canBack"`
			Path    []string    `json:"bucketsPath"`
			Records []db.Record `json:"records"`
		}{
			allDB[dbPath],
			func() bool { return len(bucketsPath) != 0 }(),
			bucketsPath,
			records,
		}
		json.NewEncoder(w).Encode(response)
	} else {
		returnError(w, nil, "Bad path of db "+dbPath, http.StatusBadRequest)
	}
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

func nextRecords(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	dbPath := r.Form.Get("dbPath")
	if _, ok := allDB[dbPath]; ok {
		records, isNextRecords, err := allDB[dbPath].NextRecords()
		if err != nil {
			returnError(w, err, "", http.StatusInternalServerError)
			return
		}

		response := struct {
			IsNextRecords bool        `json:"nextRecords"`
			Records       []db.Record `json:"records"`
		}{isNextRecords, records}
		json.NewEncoder(w).Encode(response)
	} else {
		returnError(w, nil, "Bad path of db "+dbPath, http.StatusBadRequest)
	}
}

func prevRecords(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	dbPath := r.Form.Get("dbPath")
	if _, ok := allDB[dbPath]; ok {
		records, isPrevRecords, err := allDB[dbPath].PrevRecords()
		if err != nil {
			returnError(w, err, "", http.StatusInternalServerError)
			return
		}

		response := struct {
			isPrevRecords bool        `json:"prevRecords"`
			Records       []db.Record `json:"records"`
		}{isPrevRecords, records}
		json.NewEncoder(w).Encode(response)
	} else {
		returnError(w, nil, "Bad path of db "+dbPath, http.StatusBadRequest)
	}
}
