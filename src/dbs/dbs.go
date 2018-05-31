// Package dbs is a wrapper for Package db
// Package provides functions for working with map[string]*db.BoltAPI
package dbs

import (
	"errors"
	"fmt"
	"net/http"

	"db"
)

// allDB keeps all opened databases. string – the path to the db
var allDB map[string]*db.BoltAPI

// DBInfo consist main info about db
type DBInfo struct {
	Name   string `json:"name"`
	DBPath string `json:"dbPath"`
	Size   int64  `json:"size"`
}

// Init – initializing allDB
func Init() {
	allDB = make(map[string]*db.BoltAPI)
}

func OpenDB(dbPath string) (dbName string, code int, err error) {
	// Check if db was opened
	if _, ok := allDB[dbPath]; ok {
		return "", http.StatusBadRequest, errors.New("This DB was already opened")
	}

	newDB, err := db.Open(dbPath)
	if err != nil {
		return "", http.StatusInternalServerError, err
	}

	allDB[dbPath] = newDB

	return newDB.Name, http.StatusOK, nil
}

func CloseDB(dbPath string) (code int, err error) {
	if _, ok := allDB[dbPath]; !ok {
		return http.StatusBadRequest, errors.New("There's no any db with such path (" + dbPath + ")")
	}

	dbName := allDB[dbPath].Name
	allDB[dbPath].Close()
	delete(allDB, dbPath)

	fmt.Printf("[INFO] DB \"%s\" (%s) was closed\n", dbName, dbPath)
	return http.StatusOK, nil
}

func NextBucket(dbPath, bucket string) (data db.Data, code int, err error) {
	if _, ok := allDB[dbPath]; !ok {
		return data, http.StatusBadRequest, errors.New("There's no any db with such path (" + dbPath + ")")
	}

	data, err = allDB[dbPath].Next(bucket)
	if err != nil {
		return data, http.StatusInternalServerError, err
	}

	return data, http.StatusOK, nil
}

func PrevBucket(dbPath string) (data db.Data, code int, err error) {
	if _, ok := allDB[dbPath]; !ok {
		return data, http.StatusBadRequest, errors.New("There's no any db with such path (" + dbPath + ")")
	}

	data, err = allDB[dbPath].Back()
	if err != nil {
		return data, http.StatusInternalServerError, err
	}

	return data, http.StatusOK, nil
}

func GetRoot(dbPath string) (data db.Data, code int, err error) {
	if _, ok := allDB[dbPath]; !ok {
		return data, http.StatusBadRequest, errors.New("There's no any db with such path (" + dbPath + ")")
	}

	data, err = allDB[dbPath].GetRoot()
	if err != nil {
		return data, http.StatusInternalServerError, err
	}

	return data, http.StatusOK, nil
}

func GetDBsList() (list []DBInfo) {
	for _, v := range allDB {
		info := DBInfo{Name: v.Name, DBPath: v.DBPath, Size: v.Size}
		list = append(list, info)
	}

	return list
}

func GetCurrent(dbPath string) (info DBInfo, data db.Data, code int, err error) {
	if _, ok := allDB[dbPath]; !ok {
		return info, data, http.StatusBadRequest, errors.New("There's no any db with such path (" + dbPath + ")")
	}

	info.DBPath = dbPath
	info.Name = allDB[dbPath].Name
	info.Size = allDB[dbPath].Size

	data, err = allDB[dbPath].GetCurrent()
	if err != nil {
		return info, data, http.StatusInternalServerError, err
	}

	return info, data, http.StatusOK, nil
}

func GetNextRecords(dbPath string) (data db.Data, code int, err error) {
	if _, ok := allDB[dbPath]; !ok {
		return data, http.StatusBadRequest, errors.New("There's no any db with such path (" + dbPath + ")")
	}

	data, err = allDB[dbPath].NextRecords()
	if err != nil {
		return data, http.StatusInternalServerError, err
	}

	return data, http.StatusOK, nil
}

func GetPrevRecrods(dbPath string) (data db.Data, code int, err error) {
	if _, ok := allDB[dbPath]; !ok {
		return data, http.StatusBadRequest, errors.New("There's no any db with such path (" + dbPath + ")")
	}

	data, err = allDB[dbPath].PrevRecords()
	if err != nil {
		return data, http.StatusInternalServerError, err
	}

	return data, http.StatusOK, nil
}

func Search(dbPath, mode, text string) (records []db.Record, path string, recordsAmount int, code int, err error) {
	if _, ok := allDB[dbPath]; !ok {
		return records, "", 0, http.StatusBadRequest, errors.New("There's no any db with such path (" + dbPath + ")")
	}

	if mode == "regex" {
		records, path, recordsAmount, err = allDB[dbPath].SearchRegexp(text)
	} else {
		records, path, recordsAmount, err = allDB[dbPath].Search(text)
	}

	if err != nil {
		return records, "", 0, http.StatusInternalServerError, err
	}

	return records, path, recordsAmount, http.StatusOK, nil
}

func AddBucket(dbPath, bucketName string) (code int, err error) {
	if _, ok := allDB[dbPath]; !ok {
		return http.StatusBadRequest, errors.New("There's no any db with such path (" + dbPath + ")")
	}

	err = allDB[dbPath].AddBucket(bucketName)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusCreated, nil
}

func DeleteBucket(dbPath, bucketName string) (code int, err error) {
	if _, ok := allDB[dbPath]; !ok {
		return http.StatusBadRequest, errors.New("There's no any db with such path (" + dbPath + ")")
	}

	err = allDB[dbPath].DeleteBucket(bucketName)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func AddRecord(dbPath, key, value string) (code int, err error) {
	if _, ok := allDB[dbPath]; !ok {
		return http.StatusBadRequest, errors.New("There's no any db with such path (" + dbPath + ")")
	}

	err = allDB[dbPath].AddRecord(key, value)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusCreated, nil
}

func EditRecord(dbPath, oldKey, newKey, newValue string) (code int, err error) {
	if _, ok := allDB[dbPath]; !ok {
		return http.StatusBadRequest, errors.New("There's no any db with such path (" + dbPath + ")")
	}

	err = allDB[dbPath].EditRecord(oldKey, newKey, newValue)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func DeleteRecord(dbPath, key string) (code int, err error) {
	if _, ok := allDB[dbPath]; !ok {
		return http.StatusBadRequest, errors.New("There's no any db with such path (" + dbPath + ")")
	}

	err = allDB[dbPath].DeleteRecord(key)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// CloseDBs closes all databases
func CloseDBs() {
	for k := range allDB {
		allDB[k].Close()
		delete(allDB, k)
	}
	fmt.Println("[INFO] All databases were closed")
}
