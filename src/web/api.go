package web

import (
	"db"
	"encoding/json"
	"net/http"
	"regexp"
)

// allDB keeps all opened databases. string – the path to the db
var allDB map[string]*db.DBApi

// openDB return json with information about a database
// It also adds db.DBApi to allDB
//
func openDB(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	path := r.Form.Get("filePath")

	// From C:\\users\\help (or C:\users\help) to C:/users/help
	reg := regexp.MustCompile(`\\\\|\\`)
	path = reg.ReplaceAllString(path, "/")

	// Check if db was opened
	if _, ok := allDB[path]; ok {
		http.Error(w, "This DB was already opened", http.StatusBadRequest)
		return
	}

	newDB, err := db.Open(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	allDB[path] = newDB
	w.WriteHeader(http.StatusOK)
}

// Params: filePath
// Return: -
//
func closeDB(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	dbPath := r.Form.Get("filePath")
	if _, ok := allDB[dbPath]; ok {
		allDB[dbPath].Close()
		delete(allDB, dbPath)
	}
	w.WriteHeader(http.StatusOK)
}

// next returns records from bucket with according to the name
//
// Params: filePath, bucket
// Return:
// {
// 	"canBack": bool,
//  "path": [],
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
	dbPath := r.Form.Get("filePath")
	nextBucket := r.Form.Get("bucket")
	if _, ok := allDB[dbPath]; ok {
		elements, path, err := allDB[dbPath].Next(nextBucket)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		response := struct {
			CanBack bool         `json:"canBack"`
			Path    []string     `json:"path"`
			Records []db.Element `json:"records"`
		}{
			true,
			path,
			elements,
		}
		json.NewEncoder(w).Encode(response)
	} else {
		http.Error(w, "bad path " + dbPath, http.StatusBadRequest)
	}
}

// back returns records from previous directory
//
// Params: filePath
// Return:
// {
// 	"canBack": bool,
//  "path": [],
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
	dbPath := r.Form.Get("filePath")
	if _, ok := allDB[dbPath]; ok {
		elements, path, err := allDB[dbPath].Back()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		response := struct {
			CanBack bool         `json:"canBack"`
			Path    []string     `json:"path"`
			Records []db.Element `json:"records"`
		}{
			func() bool { return len(path) != 0 }(),
			path,
			elements,
		}
		json.NewEncoder(w).Encode(response)
	} else {
		http.Error(w, "bad path " + dbPath, http.StatusBadRequest)
	}
}

// cmd returns records from root of db
//
// Params: filePath
// {
// 	"canBack": bool,
//  "path": [],
// 	"records": [
// 		{
// 			"type": "",
// 			"key": "",
// 			"value": ""
// 		},
// 	]
// }
//
func cmd(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	dbPath := r.Form.Get("filePath")
	if _, ok := allDB[dbPath]; ok {
		elements, _, err := allDB[dbPath].GetCMD()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		response := struct {
			CanBack bool         `json:"canBack"`
			Path    []string     `json:"path"`
			Records []db.Element `json:"records"`
		}{
			false,
			[]string{},
			elements,
		}
		json.NewEncoder(w).Encode(response)
	} else {
		http.Error(w, "bad path " + dbPath, http.StatusBadRequest)
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
	var list []db.DBApi
	for _, v := range allDB {
		list = append(list, *v)
	}
	json.NewEncoder(w).Encode(list)
}

// current returns records in current bucket
//
// Params: filePath
// Return:
// {
// 	"name": "",
// 	"filePath": "",
// 	"size": 0,
//	"canBack": bool,
//  "path": [],
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
	dbPath := r.Form.Get("filePath")

	if _, ok := allDB[dbPath]; ok {
		elements, path, err := allDB[dbPath].GetCurrent()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		response := struct {
			*db.DBApi
			CanBack  bool         `json:"canBack"`
			Path     []string     `json:"path"`
			Elements []db.Element `json:"records"`
		}{
			allDB[dbPath],
			func() bool { return len(path) != 0 }(),
			path,
			elements,
		}
		json.NewEncoder(w).Encode(response)
	} else {
		http.Error(w, "bad path " + dbPath, http.StatusBadRequest)
	}
}