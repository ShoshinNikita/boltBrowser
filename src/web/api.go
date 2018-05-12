package web

import (
	"net/http"
	"db"
	"encoding/json"
	"regexp"
)

// allDB keeps all opened databases. string – the path to the db
var allDB map[string]*db.DBApi

// openDB return json with information about a database
// It also adds db.DBApi to allDB
//
// Params: path
// Return:
// {
// 		"name": "",
// 		"path": "",
// 		"size": 0
// }
//
func openDB(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	path := r.Form.Get("path")

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
	json.NewEncoder(w).Encode(newDB)
}

// Params: path
// Return: -
// 
func closeDB(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	dbPath := r.Form.Get("path")
	if _, ok := allDB[dbPath]; ok {
		allDB[dbPath].Close()
		delete(allDB, dbPath)
	}
	w.WriteHeader(http.StatusOK)
}

// next returns records from bucket with according to the name
//
// Params: path, bucket
// Return:
// {
// 	"canBack": bool,
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
	dbPath := r.Form.Get("path")
	nextBucket := r.Form.Get("bucket")
	if _, ok := allDB[dbPath]; ok {
		elements, err := allDB[dbPath].Next(nextBucket)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		response := struct {
			CanBack	bool `json:"canBack"`
			Records []db.Element `json:"records"`
		} {
			true,
			elements,
		}
		json.NewEncoder(w).Encode(response)
	} else {
		http.Error(w, "bad path " + dbPath, http.StatusBadRequest)
	}
}

// back returns records from previous directory
//
// Params: path
// Return:
// {
// 	"canBack": bool,
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
	dbPath := r.Form.Get("path")
	if _, ok := allDB[dbPath]; ok {
		elements, canBack, err := allDB[dbPath].Back()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		response := struct {
			CanBack	bool `json:"canBack"`
			Records []db.Element `json:"records"`
		} {
			canBack,
			elements,
		}
		json.NewEncoder(w).Encode(response)
	} else {
		http.Error(w, "bad path " + dbPath, http.StatusBadRequest)
	}
}

// cmd returns records from root of db
//
// Params: path
// {
// 	"canBack": bool,
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
	dbPath := r.Form.Get("path")
	if _, ok := allDB[dbPath]; ok {
		elements, canBack, err := allDB[dbPath].GetCMD()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		response := struct {
			CanBack	bool `json:"canBack"`
			Records []db.Element `json:"records"`
		} {
			canBack,
			elements,
		}
		json.NewEncoder(w).Encode(response)
	} else {
		http.Error(w, "bad path " + dbPath, http.StatusBadRequest)
	}
}

// Params: 
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
// Params: path
// Return:
// {
// 	"name": "",
// 	"path": "",
// 	"size": 0,
//	"canBack": bool,
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
	dbPath := r.Form.Get("path")
	if _, ok := allDB[dbPath]; ok {
		elements, canBack, err := allDB[dbPath].GetCurrent()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		response := struct {
			*db.DBApi
			CanBack bool `json:"canBack"`
			Elements []db.Element `json:"records"`
		} {
			allDB[dbPath],
			canBack,
			elements,
		}
		json.NewEncoder(w).Encode(response)
	} else {
		http.Error(w, "bad path " + dbPath, http.StatusBadRequest)
	}
}